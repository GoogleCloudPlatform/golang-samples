package tictactoe

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zabawaba99/firego"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/opened", openedHandler)
	http.HandleFunc("/move", moveHandler)
	http.HandleFunc("/delete", deleteHandler)
}

var tmpl = template.Must(template.ParseFiles("template.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	config, err := ioutil.ReadFile("firebase_config.html")
	if err != nil {
		handleError(w, r, "Could not read config file", err)
		return
	}

	u := user.Current(ctx)

	var game *Game
	if gameID := r.FormValue("g"); gameID == "" {
		game = NewGame()
		game.K = datastore.NewKey(ctx, "Game", u.ID, 0, nil)
		game.UserX = u.ID

		if _, err := datastore.Put(ctx, game.K, game); err != nil {
			handleError(w, r, "Could not start game", err)
			return
		}
	} else {
		// Existing game, join it.
		if err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
			g, err := gameFromRequest(r)
			if err != nil {
				return err
			}
			game = g
			if game.UserO == "" {
				game.UserO = u.ID
				_, err := datastore.Put(ctx, game.K, game)
				return err
			}
			return nil
		}, nil); err != nil {
			handleError(w, r, "Could not join game", err)
			return
		}
	}

	gameKey := game.K.Encode()
	channelID := u.ID + gameKey

	tok, err := createToken(ctx, channelID)
	if err != nil {
		handleError(w, r, "Could not create auth token", err)
		return
	}

	d := struct {
		GameKey        string
		GameLink       string
		Me             string
		Token          string
		ChannelID      string
		FirebaseConfig template.HTML
		InitialMessage *Game
	}{
		FirebaseConfig: template.HTML(config),
		Token:          tok,
		GameKey:        gameKey,
		GameLink:       r.URL.Host + "/?g=" + gameKey,
		Me:             u.ID,
		ChannelID:      channelID,
		InitialMessage: game,
	}
	if err := tmpl.Execute(w, d); err != nil {
		handleError(w, r, "Could not execute template", err)
		return
	}
}

func createToken(ctx context.Context, channelID string) (string, error) {
	iss, err := appengine.ServiceAccount(ctx)
	if err != nil {
		return "", err
	}
	iat := time.Now().Unix()
	jwt := map[string]interface{}{
		"alg": "RS256",
		"iss": iss,
		"sub": iss,
		"aud": "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
		"iat": iat,
		"exp": iat + 3600, // 1 hour
		"uid": channelID,
	}
	body, err := json.Marshal(jwt)
	if err != nil {
		return "", err
	}
	header := base64.StdEncoding.EncodeToString([]byte(`{"typ":"JWT","alg":"RS256"}`))
	payload := append([]byte(header), byte('.'))
	payload = append(payload, []byte(base64.StdEncoding.EncodeToString(body))...)
	_, sig, err := appengine.SignBytes(ctx, payload)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", payload, base64.StdEncoding.EncodeToString(sig)), nil
}

func gameFromRequest(r *http.Request) (*Game, error) {
	ctx := appengine.NewContext(r)

	k, err := datastore.DecodeKey(r.FormValue("g"))
	if err != nil {
		return nil, fmt.Errorf("Invalid game ID: %v", err)
	}
	var g Game
	if err := datastore.Get(ctx, k, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func firebase(ctx context.Context) (*firego.Firebase, error) {
	hc, err := google.DefaultClient(ctx,
		"https://www.googleapis.com/auth/firebase.database",
		"https://www.googleapis.com/auth/userinfo.email",
	)
	if err != nil {
		return nil, err
	}
	base := os.Getenv("FIREBASE_BASE")
	if base == "" {
		// Check the environment variable for the base firebase URL.
		//
		// The config should look like:
		//
		// env_variables:
		//    FIREBASE_BASE: https://app-id.firebase.io.com
		//
		return nil, errors.New("Missing FIREBASE_BASE environment variable.")
	}
	return firego.New(base, hc), nil
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	g, err := gameFromRequest(r)
	if err != nil {
		handleError(w, r, "Could not get game", err)
		return
	}

	position, err := strconv.Atoi(r.FormValue("i"))
	if err != nil {
		handleError(w, r, "Invalid position", err)
		return
	}
	if position < 0 || position > 8 {
		handleError(w, r, "Expected position between 0 and 8", nil)
		return
	}

	u := user.Current(ctx).ID
	expectedUser := g.UserX
	if !g.MoveX {
		expectedUser = g.UserO
	}
	if u != expectedUser {
		handleError(w, r, "Not your move", nil)
		return
	}

	// Update the game board.
	if err := g.MoveAt(position); err != nil {
		handleError(w, r, "Could not move", err)
	}
	g.MoveX = !g.MoveX

	if winner, isWon := g.CheckWin(); isWon {
		if winner == "O" {
			g.Winner = g.UserO
		} else if winner == "X" {
			g.Winner = g.UserX
		} else {
			g.Winner = "No one"
		}
		g.WinningBoard = g.Board // TODO: implement patterns
	}

	if _, err := datastore.Put(ctx, g.K, g); err != nil {
		handleError(w, r, "Could not save game", err)
		return
	}
	sendUpdate(ctx, g)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	game, err := gameFromRequest(r)
	if err != nil {
		handleError(w, r, "Could not get game", err)
		return
	}
	fb, err := firebase(ctx)

	channelID := user.Current(ctx).ID + game.K.Encode()

	if err := fb.Child("channels").Child(channelID).Remove(); err != nil {
		handleError(w, r, "Could not delete channel", err)
		return
	}
	log.Infof(ctx, "Deleted channel %v", channelID)
}

func openedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	k, err := datastore.DecodeKey(r.FormValue("g"))
	if err != nil {
		handleError(w, r, "Invalid game ID", err)
		return
	}
	var g Game
	if err := datastore.Get(ctx, k, &g); err != nil {
		handleError(w, r, "Could not get game", err)
		return
	}
	sendUpdate(ctx, &g)
}

func sendUpdate(ctx context.Context, g *Game) {
	fb, err := firebase(ctx)
	if err != nil {
		log.Errorf(ctx, "getFirebase: %v", err)
	}
	chans := fb.Child("channels")

	gameKey := g.K.Encode()

	if g.UserO != "" {
		channelID := g.UserO + gameKey
		if err := chans.Child(channelID).Set(g); err != nil {
			log.Errorf(ctx, "Updating UserO (%s): %v", channelID, err)
		} else {
			log.Infof(ctx, "Update O sent.")
		}
	}

	if g.UserX != "" {
		channelID := g.UserX + gameKey
		if err := chans.Child(channelID).Set(g); err != nil {
			log.Errorf(ctx, "Updating UserX (%s): %v", channelID, err)
		} else {
			log.Infof(ctx, "Update X sent.")
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, message string, err error) {
	msg := message
	if err != nil {
		msg = fmt.Sprintf("%s: %v", message, err)
	}
	ctx := appengine.NewContext(r)
	http.Error(w, msg, 500)
	log.Errorf(ctx, "%s", msg, err)
}
