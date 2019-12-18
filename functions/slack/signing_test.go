package slack

import (
	"encoding/hex"
	"fmt"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGoodSignature(t *testing.T) {

	secret := "talesfromthecrypt"
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	body := "somebody"
	base := fmt.Sprintf("v0:%s:%s", ts, body)
	correctSHA2Signature := fmt.Sprintf("v0=%s", hex.EncodeToString(getSignature([]byte(base), []byte("talesfromthecrypt"))))

	goodReq := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader(body))
	goodReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	goodReq.Header.Add("X-Slack-Request-Timestamp", ts)
	goodReq.Header.Add("X-Slack-Signature", correctSHA2Signature)

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
	wrongSHA2Signature := "v0=146abde6763faeba19adc4d9fe4961668f4be11f7405a1c05b636f29312eac2e"
	body := "somebody"

	badReq := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader(body))
	badReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	badReq.Header.Add("X-Slack-Request-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	badReq.Header.Add("X-Slack-Signature", wrongSHA2Signature)

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
	ts := "1504928418"
	body := "somebody"
	base := fmt.Sprintf("v0:%s:%s", ts, body)
	correctSHA2Signature := fmt.Sprintf("v0=%s", hex.EncodeToString(getSignature([]byte(base), []byte("talesfromthecrypt"))))

	badTSReq := httptest.NewRequest("POST", "https://someurl.com", strings.NewReader(body))
	badTSReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	badTSReq.Header.Add("X-Slack-Request-Timestamp", ts)
	badTSReq.Header.Add("X-Slack-Signature", correctSHA2Signature)

	badResult, err := verifyRequestSignature(badTSReq, secret)
	if err == nil {
		t.Errorf("TestBadTimestamp expected an error but got nil")
	}
	if badResult {
		t.Errorf("TestBadTimestamp failed. Got: %v | Expected: %v", badResult, badResult)
	}
}
