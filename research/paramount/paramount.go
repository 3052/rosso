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

// 576p L3
func Widevine(at, contentId string, cbsCom *http.Cookie) (*SessionToken, error) {
   var url_data url.URL
   url_data.Scheme = "https"
   url_data.Host = "www.cbs.com"
   if cbsCom != nil {
      url_data.Path = "/apps-api/v3.1/androidphone/irdeto-control/session-token.json"
   } else {
      url_data.Path = "/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json"
   }
   url_data.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   var req http.Request
   req.URL = &url_data
   req.Header = http.Header{}
   if cbsCom != nil {
      req.AddCookie(cbsCom)
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   defer resp.Body.Close()
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.StreamingUrl == "" {
      return nil, errors.New("StreamingUrl")
   }
   return &result, nil
}

type SessionToken struct {
   Message      string
   LsSession    string `json:"ls_session"`
   Url          string
   StreamingUrl string
}

func (s *SessionToken) Send(data []byte) ([]byte, error) {
   url_data, err := url.Parse(s.Url)
   if err != nil {
      return nil, err
   }
   //url_data.Path = "/playready/rightsmanager.asmx"
   var req http.Request
   req.Method = "POST"
   req.URL = url_data
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+s.LsSession)
   req.Body = io.NopCloser(bytes.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
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
