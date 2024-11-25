package encoders

import (
	"encoding/base64"
	"fmt"
)

func EncodeUserUUID(username string) []byte {
	byteUsername := []byte(username)
	uuid := make([]byte, base64.StdEncoding.EncodedLen(len(byteUsername)))
	base64.StdEncoding.Encode(uuid, byteUsername)
	return uuid
}

func DecodeBearerToken(bearer string) (string, error) {
	username, err := base64.StdEncoding.DecodeString(bearer)
	if err != nil {
		return "", fmt.Errorf("error decoding bearer token: %v", err)
	}
	return string(username), nil
}
