package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

const maxWebhookBody = 1 << 20

func validWebhook(r *http.Request, secret string) bool {
	if secret == "" {
		return true
	}
	switch r.PathValue("host") {
	case "github":
		body, err := io.ReadAll(io.LimitReader(r.Body, maxWebhookBody))
		if err != nil {
			return false
		}
		signature, ok := strings.CutPrefix(r.Header.Get("X-Hub-Signature-256"), "sha256=")
		if !ok {
			return false
		}
		got, err := hex.DecodeString(signature)
		if err != nil {
			return false
		}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		return hmac.Equal(got, mac.Sum(nil))
	case "azure":
		return equal(r.Header.Get("X-Fuda-Secret"), secret)
	}
	return false
}
