package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const version = "v0"
const slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
const slackSignatureHeader = "X-Slack-Signature"

func getSignature(baseString string, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(baseString))
	computedReqSignature := h.Sum(nil)

	return computedReqSignature
}

func checkTimestamp(timeStamp int64) bool {
	t := time.Since(time.Unix(timeStamp, 0)).Minutes()
	if t > 5 {
		return false
	}
	return true
}

func verifyRequestSignature(r *http.Request, slackSigningSecret string) (bool, error) {
	timeStamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	t, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		log.Fatal("error: couldn't parse the timestamp header from Slack")
	}

	if !checkTimestamp(t) {
		return false, fmt.Errorf("timestamp too old")
	}

	if timeStamp == "" || slackSignature == "" {
		return false, fmt.Errorf("either timeStamp or signature headers were blank")
	}

	theBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("couldn't read the request's body")
	}

	// Reset the body so other calls won't fail.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(theBody))

	baseString := fmt.Sprintf("%s:%s:%s", version, timeStamp, theBody)

	computedReqSignature := getSignature(baseString, slackSigningSecret)

	byteSignature, err := hex.DecodeString(strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version)))

	log.Println("Slack Signature: ", slackSignature)
	log.Println("Computed Signature: ", hex.EncodeToString(computedReqSignature))
	log.Println("Base String: ", baseString)

	if err != nil {
		return false, err
	}

	// request signature and computed signature match.
	if hmac.Equal(computedReqSignature, byteSignature) {
		return true, nil
	}

	// request signature and computed signature do not match.
	return false, nil
}
