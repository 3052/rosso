package main

import (
   "crypto/rand"
   "crypto/sha256"
   "encoding/base64"
)

// generateRandomString generates a URL-safe random string of the given length.
func generateRandomString(length int) (string, error) {
   b := make([]byte, length)
   _, err := rand.Read(b)
   if err != nil {
      return "", err
   }
   return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

// pkcePair generates a code_verifier and its corresponding code_challenge (S256).
func pkcePair() (verifier string, challenge string, err error) {
   verifier, err = generateRandomString(43)
   if err != nil {
      return "", "", err
   }

   h := sha256.Sum256([]byte(verifier))
   challenge = base64.RawURLEncoding.EncodeToString(h[:])
   return verifier, challenge, nil
}
