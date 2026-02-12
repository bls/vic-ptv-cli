package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// SignURL signs a PTV API URL path using HMAC-SHA1.
// The path should include /v3/... and any query parameters plus devid.
// Returns the full URL with the signature appended.
func SignURL(baseURL, path, devID, apiKey string) (string, error) {
	// Ensure the path includes devid
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	pathWithDevID := fmt.Sprintf("%s%sdevid=%s", path, separator, devID)

	// HMAC-SHA1 the full path using the API key
	key := []byte(apiKey)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(pathWithDevID))
	signature := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s%s&signature=%s", baseURL, pathWithDevID, strings.ToUpper(signature)), nil
}
