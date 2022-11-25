package algorithms

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

const fingerPrintKey = "sdfsdfsfsdf"

func GetFingerPrint(data []string) string {
	strings.Join(data, ".")
	alg := hmac.New(sha256.New, []byte(fingerPrintKey))
	alg.Write([]byte(strings.Join(data, ".")))
	return base64.StdEncoding.EncodeToString(alg.Sum(nil))
}
