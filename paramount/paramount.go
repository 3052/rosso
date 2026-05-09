package paramount

import (
   "41.neocities.org/maya"
   "archive/zip"
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "regexp"
   "slices"
   "strings"
)

func CbsAppIds() string {
   var data strings.Builder
   for i, app := range CbsApps {
      if i >= 1 {
         data.WriteByte(' ')
      }
      data.WriteString(app.Id)
   }
   return data.String()
}

func GetCbsApp(id string) (*CbsApp, error) {
   for _, app := range CbsApps {
      if app.Id == id {
         return &app, nil
      }
   }
   return nil, fmt.Errorf("CBS app not found %q", id)
}

var CbsApps = []CbsApp{
   {
      Id:      "com.cbs.app",
      Host:    "www.paramountplus.com",
      Secret:  "7081400bd4143bf3",
      Version: "Paramount+ 16.8.0",
   },
   {
      Id:      "com.cbs.ca",
      Host:    "www.paramountplus.com",
      Secret:  "1c5d27627d71b420",
      Version: "Paramount+ 16.8.0",
   },
   {
      Id:      "com.cbs.tve",
      Host:    "www.cbs.com",
      Secret:  "cef32931dc01412e",
      Version: "CBS 15.6.0",
   },
}

type CbsApp struct {
   Id      string
   Host    string
   Secret  string
   Version string
}

type Session struct {
   LsSession    string `json:"ls_session"`
   Message      string
   StreamingUrl string // MPD
   Url          string // License Server
}

var hexPattern = regexp.MustCompile(`\x00\x10([0-9a-f]{16})\x00`)

// ExtractDexHexBytes returns a set (map) of unique 16-character hex strings
// found in .dex files
func ExtractDexHexBytes(name string) (map[string]struct{}, error) {
   results := make(map[string]struct{})
   reader, err := zip.OpenReader(name)
   if err != nil {
      return nil, err
   }
   for _, f := range reader.File {
      if strings.HasSuffix(f.Name, ".dex") {
         content, err := readZipFile(f)
         if err != nil {
            return nil, err
         }
         matches := hexPattern.FindAllSubmatch(content, -1)
         for _, match := range matches {
            results[string(match[1])] = struct{}{}
         }
      }
   }
   return results, nil
}

func readZipFile(f *zip.File) ([]byte, error) {
   rc, err := f.Open()
   if err != nil {
      return nil, err
   }
   defer rc.Close()
   return io.ReadAll(rc)
}

func (s *Session) GetManifest() (*url.URL, error) {
   return url.Parse(s.StreamingUrl)
}

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

func get_at(app_secret string) (string, error) {
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
   // 3 & 4. Create payload: "|" + app_secret
   data := []byte{'|'}
   data = append(data, app_secret...)
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

func pkcs7_pad(data []byte, block_size int) []byte {
   // Calculate the number of padding bytes needed.
   paddingLen := block_size - (len(data) % block_size)
   // Create a padding byte (the value is the length of the padding)
   padByte := byte(paddingLen)
   // Append the padding byte 'paddingLen' times
   for i := 0; i < paddingLen; i++ {
      data = append(data, padByte)
   }
   return data
}

func (s *Session) Fetch(body []byte) ([]byte, error) {
   target, err := url.Parse(s.Url)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"authorization": "Bearer " + s.LsSession}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type Cookie struct {
   Name  string
   Value string
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func (c *CbsApp) FetchCbsCom(username, password string) (*Cookie, error) {
   at, err := get_at(c.Secret)
   if err != nil {
      return nil, err
   }
   body := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme:   "https",
         Host:     c.Host,
         Path:     "/apps-api/v2.0/androidphone/auth/login.json",
         RawQuery: url.Values{"at": {at}}.Encode(),
      },
      map[string]string{
         "content-type": "application/x-www-form-urlencoded",
         "user-agent":   "!", // randomly fails if this is missing
      },
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return nil, err
   }
   for _, c := range resp.Cookies() {
      if c.Name == "CBS_COM" {
         return &Cookie{Name: c.Name, Value: c.Value}, nil
      }
   }
   return nil, errors.New("CBS_COM cookie not present")
}

func (c *CbsApp) fetch_session(platform, contentId string, cbs_com *Cookie) (*Session, error) {
   at, err := get_at(c.Secret)
   if err != nil {
      return nil, err
   }
   endpoint := "anonymous-session-token.json"
   header := map[string]string{}
   if cbs_com != nil {
      endpoint = "session-token.json"
      header["cookie"] = cbs_com.String()
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   c.Host,
         Path:   fmt.Sprintf("/apps-api/v3.1/%s/irdeto-control/%s", platform, endpoint),
         RawQuery: url.Values{
            "at":        {at},
            "contentId": {contentId},
         }.Encode(),
      },
      header,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Session
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

func (c *CbsApp) FetchWidevine(contentId string, cbsCom *Cookie) (*Session, error) {
   return c.fetch_session("androidphone", contentId, cbsCom)
}

func (c *CbsApp) FetchPlayReady(contentId string, cbsCom *Cookie) (*Session, error) {
   return c.fetch_session("xboxone", contentId, cbsCom)
}

func (c *CbsApp) FetchStreamingUrl(contentId string, cbsCom *Cookie) (*Session, error) {
   result, err := c.fetch_session("androidphone", contentId, cbsCom)
   if err != nil {
      return nil, err
   }
   if result.StreamingUrl == "" {
      return nil, errors.New("streamingUrl (MPD) is missing")
   }
   return result, nil
}
