package paramount

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "slices"
   "strings"
)

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func FetchCbsCom(at, username, password string) (*http.Cookie, error) {
   body := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   req, err := http.NewRequest(
      "POST",
      "https://www.paramountplus.com/apps-api/v2.0/androidphone/auth/login.json",
      strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{"at": {at}}.Encode()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   // randomly fails if this is missing
   req.Header.Set("user-agent", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "CBS_COM" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}

///

var Apps = []struct {
   url        string
   id         string
   version    string
   app_secret string
}{
   {
      url:        "https://apkmirror.com/apk/viacomcbs-streaming/paramount",
      id:         "com.cbs.app",
      version:    "Paramount+ 16.8.0",
      app_secret: "7081400bd4143bf3",
   },
   {
      url:        "https://apkmirror.com/apk/cbs-interactive-inc/cbs",
      id:         "com.cbs.tve",
      version:    "CBS 15.6.0",
      app_secret: "cef32931dc01412e",
   },
   {
      url:        "https://apkmirror.com/apk/viacomcbs-streaming/paramount-4",
      id:         "com.cbs.ca",
      version:    "Paramount+ 16.8.0",
      app_secret: "1c5d27627d71b420",
   },
}

func pkcs7_pad(data []byte, blockSize int) []byte {
   // Calculate the number of padding bytes needed.
   // If data is already a multiple of blockSize, this results in a full block
   // of padding.
   paddingLen := blockSize - (len(data) % blockSize)
   // Create a padding byte (the value is the length of the padding)
   padByte := byte(paddingLen)
   // Append the padding byte 'paddingLen' times
   for i := 0; i < paddingLen; i++ {
      data = append(data, padByte)
   }
   return data
}

func GetAt(appSecret string) (string, error) {
   // 1. Decode hex secret key
   key, err := hex.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   // 2. Create aes cipher with key
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   // 3 & 4. Create payload: "|" + appSecret
   data := []byte{'|'}
   data = append(data, appSecret...)
   // 5. Apply PKCS7 Padding (Separate Function)
   data = pkcs7_pad(data, aes.BlockSize)
   // Prepare Empty IV (16 bytes of zeros)
   var iv [aes.BlockSize]byte
   // 6. CBC encrypt with empty IV
   // We encrypt 'data' in place
   cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(data, data)
   // 8. Create Header for block size (uint16)
   size := binary.BigEndian.AppendUint16(nil, aes.BlockSize)
   // 7 & 8. Combine [Size] + [IV] + [Encrypted Data]
   data = slices.Concat(size, iv[:], data)
   // 9. Return result base64 encoded
   return base64.StdEncoding.EncodeToString(data), nil
}

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"
