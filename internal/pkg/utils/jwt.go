package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"hafiedh.com/downloader/internal/config"
)

type (
	UploadRequest struct {
		FileName          string `json:"fileName"`
		UseUniqueFileName string `json:"useUniqueFileName"`
		Folder            string `json:"folder"`
		IsPrivateFile     string `json:"isPrivateFile"`
		Iat               int64  `json:"iat"`
		Exp               int64  `json:"exp"`
	}
)

func ImageKitJwtSign(req UploadRequest) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
		"kid": config.GetString("imagekit.publicKey"),
	}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)

	req.Iat = time.Now().Unix()
	req.Exp = req.Iat + 3600
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	secret := config.GetString("imagekit.privateKey")
	signatureData := encodedHeader + "." + encodedPayload
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signatureData))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	token := encodedHeader + "." + encodedPayload + "." + signature
	return token, nil
}
