package paramount

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "slices"
)

// 1080p SL2000
// 1440p SL2000 + cookie
func PlayReady(at, contentId string, cbsCom *http.Cookie) (*SessionToken, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.cbs.com"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   if cbsCom != nil {
      req.AddCookie(cbsCom)
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/session-token.json"
   } else {
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/anonymous-session-token.json"
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Errors != "" {
      return nil, errors.New(result.Errors)
   }
   return &result, nil
}

type SessionToken struct {
   Errors    string
   LsSession string `json:"ls_session"`
   StreamingUrl string
   Url       string
}

func (s *SessionToken) Dash() (*Dash, error) {
   resp, err := http.Get(s.StreamingUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}
func (s *SessionToken) Send(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.LsSession)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(string(data))
   }
   return data, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

func FetchAppSecret() (string, error) {
   return "cef32931dc01412e", nil
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

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"
