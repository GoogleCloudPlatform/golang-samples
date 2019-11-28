package slack

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func badRequest() *http.Request {
	wrongSHA2Signature := "v0=146abde6763faeba19adc4d9fe4961668f4be11f7405a1c05b636f29312eac2e"
	body := "somebody"

	req := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Slack-Request-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Add("X-Slack-Signature", wrongSHA2Signature)

	return req
}

func goodRequest() *http.Request {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	body := "somebody"
	base := fmt.Sprintf("v0:%s:%s", ts, body)
	correctSHA2Signature := fmt.Sprintf("v0=%s", hex.EncodeToString(getSignature(base, "talesfromthecrypt")))

	req := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Slack-Request-Timestamp", ts)
	req.Header.Add("X-Slack-Signature", correctSHA2Signature)

	return req
}

func badTimeStamp() *http.Request {
	ts := "1504928418"
	body := "somebody"
	base := fmt.Sprintf("v0:%s:%s", ts, body)
	correctSHA2Signature := fmt.Sprintf("v0=%s", hex.EncodeToString(getSignature(base, "talesfromthecrypt")))

	req := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Slack-Request-Timestamp", ts)
	req.Header.Add("X-Slack-Signature", correctSHA2Signature)

	return req
}

func TestGoodSignature(t *testing.T) {

	secret := "talesfromthecrypt"
	goodReq := goodRequest()

	goodResult, err := verifyRequestSignature(goodReq, secret)
	if err != nil {
		t.Errorf("TestGoodSignature got an error: %s", err)
	}
	if !goodResult {
		t.Errorf("TestGoodSignature failed. Got: %v | Expected: %v", goodResult, !goodResult)
	}
}

func TestBadSignature(t *testing.T) {

	secret := "talesfromthecrypt"

	badReq := badRequest()

	badResult, err := verifyRequestSignature(badReq, secret)
	if err != nil {
		t.Errorf("TestBadSignature got an error: %s", err)
	}
	if badResult {
		t.Errorf("TestBadSignature failed. Got: %v | Expected: %v", badResult, badResult)
	}
}

func TestBadTimestamp(t *testing.T) {

	secret := "talesfromthecrypt"

	badTSReq := badTimeStamp()

	badResult, err := verifyRequestSignature(badTSReq, secret)
	if err == nil {
		t.Errorf("TestBadTimestamp expected an error but got nil")
	}
	if badResult {
		t.Errorf("TestBadTimestamp failed. Got: %v | Expected: %v", badResult, badResult)
	}
}
