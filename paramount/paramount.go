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
   "fmt"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strings"
)

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func FetchCbsCom(at, username, password string) (*http.Cookie, error) {
   data := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.paramountplus.com",
      Path:     "/apps-api/v2.0/androidphone/auth/login.json",
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   req.Header = http.Header{}
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   // randomly fails if this is missing
   req.Header.Set("user-agent", "!")
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
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

func join(items ...string) string {
   return strings.Join(items, "")
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

func FetchAppSecret() (string, error) {
   resp, err := http.Head("https://www.paramountplus.com")
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   switch resp.Header.Get("x-real-server") {
   case "us_www_web_prod_vip1":
      return AppSecrets[0].Us, nil
   case "international_www_web_prod_vip1":
      return AppSecrets[0].International, nil
   }
   return "", errors.New("unexpected or missing server header value")
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

func (i *Item) Dash() (*Dash, error) {
   resp, err := http.Get(i.StreamingUrl)
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

func FetchItem(at, cid string, cbsCom *http.Cookie) (*Item, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.paramountplus.com",
      Path:     join("/apps-api/v2.0/androidphone/video/cid/", cid, ".json"),
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   if cbsCom != nil {
      req.AddCookie(cbsCom)
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      ItemList []Item
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   item := result.ItemList[0]
   if item.StreamingUrl == "" {
      return nil, fmt.Errorf(
         "subscription level = %v, CBS_COM = %#v",
         item.SubscriptionLevel, cbsCom,
      )
   }
   return &item, nil
}

type Item struct {
   StreamingUrl      string
   SubscriptionLevel string
}

type SessionToken struct {
   Errors    string
   LsSession string `json:"ls_session"`
   Url       string
}

// 1080p SL2000
// 1440p SL2000 + cookie
func PlayReady(at, contentId string, cbsCom *http.Cookie) (*SessionToken, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
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

// 576p L3
func Widevine(at, contentId string) (*SessionToken, error) {
   var req http.Request
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.Path = "/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(data.String())
   }
   defer resp.Body.Close()
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
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

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

var AppSecrets = []struct {
   Version       string
   Us            string
   International string
}{
   {
      Version:       "16.4.1",
      Us:            "7cd07f93a6e44cf7",
      International: "68b4475a49bed95a",
   },
   {
      Version:       "16.0.0",
      Us:            "9fc14cb03691c342",
      International: "6c68178445de8138",
   },
}
