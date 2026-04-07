package main

import (
   "crypto/hmac"
   "crypto/md5"
   "crypto/sha256"
   "encoding/base64"
   "fmt"
)

const (
   clientKey    = "web.NhFyz4KsZ54"
   clientSecret = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

// GenerateAuthorizationHeader dynamically creates the "Client key=...,time=...,sig=..." string
func GenerateAuthorizationHeader(url string, body []byte, timestamp int64) string {
   // Decode the static Base64URL secret
   secretBytes, _ := base64.RawURLEncoding.DecodeString(clientSecret)

   // 1. MD5 Hash the JSON body
   bodyHashRaw := md5.Sum(body)

   // 2. Base64URL encode the MD5 hash
   bodyHashStr := base64.RawURLEncoding.EncodeToString(bodyHashRaw[:])

   // 3. Concatenate requested URL + body hash + timestamp
   message := fmt.Sprintf("%s%s%d", url, bodyHashStr, timestamp)

   // 4. HMAC-SHA256 the concatenated message
   mac := hmac.New(sha256.New, secretBytes)
   mac.Write([]byte(message))

   // 5. Base64URL encode the HMAC result
   sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

   return fmt.Sprintf("Client key=%s,time=%d,sig=%s", clientKey, timestamp, sig)
}
