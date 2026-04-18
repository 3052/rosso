package paramount

import (
   "archive/zip"
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "errors"
   "fmt"
   "io"
   "maps"
   "net/url"
   "regexp"
   "slices"
   "strings"
)

func (a *App) FetchWidevine(contentId, cbsCom string) (*Session, error) {
   return a.fetch_session("androidphone", contentId, cbsCom)
}

func (a *App) FetchPlayReady(contentId, cbsCom string) (*Session, error) {
   return a.fetch_session("xboxone", contentId, cbsCom)
}

func (a *App) FetchStreamingUrl(contentId, cbsCom string) (*Session, error) {
   result, err := a.fetch_session("androidphone", contentId, cbsCom)
   if err != nil {
      return nil, err
   }
   if result.StreamingUrl == "" {
      return nil, errors.New("streamingUrl (MPD) is missing")
   }
   return result, nil
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

var apps = map[string]App{
   "com.cbs.app": {
      Host:    "www.paramountplus.com",
      Version: "Paramount+ 16.8.0",
      Secret:  "7081400bd4143bf3",
   },
   "com.cbs.ca": {
      Host:    "www.paramountplus.com",
      Version: "Paramount+ 16.8.0",
      Secret:  "1c5d27627d71b420",
   },
   "com.cbs.tve": {
      Host:    "www.cbs.com",
      Version: "CBS 15.6.0",
      Secret:  "cef32931dc01412e",
   },
}

func GetAppKeys() []string {
   return slices.Sorted(maps.Keys(apps))
}

func get_at(appSecret string) (string, error) {
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

type App struct {
   Host    string
   Version string
   Secret  string
}

func GetApp(key string) (*App, error) {
   app, exists := apps[key]
   if !exists {
      return nil, fmt.Errorf("app not found: %s", key)
   }
   return &app, nil
}
