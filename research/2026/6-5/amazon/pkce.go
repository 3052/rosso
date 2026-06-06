// pkce.go
package amazon

import (
   "crypto/rand"
   "crypto/sha256"
   "encoding/base64"
)

// GeneratePKCE creates a random code verifier and its corresponding S256 code challenge.
func GeneratePKCE() (codeVerifier, codeChallenge string, err error) {
   // Generate a 32-byte random string for the code verifier
   b := make([]byte, 32)
   if _, err := rand.Read(b); err != nil {
      return "", "", err
   }

   // Base64Url encode without padding to create the verifier
   codeVerifier = base64.RawURLEncoding.EncodeToString(b)

   // SHA256 hash the verifier
   hash := sha256.Sum256([]byte(codeVerifier))

   // Base64Url encode the hash without padding to create the challenge
   codeChallenge = base64.RawURLEncoding.EncodeToString(hash[:])

   return codeVerifier, codeChallenge, nil
}
