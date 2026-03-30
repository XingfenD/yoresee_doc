package utils

import "encoding/base64"

func DecodeBase64(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}
