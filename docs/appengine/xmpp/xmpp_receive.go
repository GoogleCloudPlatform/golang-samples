// [START example_sending]
package demo

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/xmpp"
)

func init() {
	http.HandleFunc("/send", sendHandler)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	m := &xmpp.Message{
		To:   []string{"example@gmail.com"},
		Body: "Someone has sent you a gift: http://example.com/gifts/",
	}
	err := m.Send(ctx)
	if err != nil {
		// Send an email message instead...
	}
}

// [END example_sending]
