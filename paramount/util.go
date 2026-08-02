// - what is `androidphone` MPD? 2160p
// - what is `xboxone` MPD? 1080p
// - what is the L3 cookie max? 576p
// - what is the L3 no cookie max? 576p
// - what is the SL2000 cookie max? 2160p
// - what is the SL2000 no cookie max? 1080p
package paramount

import (
   "archive/zip"
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "io"
   "log"
   "net/http"
   "regexp"
   "strings"
)

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

var hexPattern = regexp.MustCompile(`\x00\x10([0-9a-f]{16})\x00`)

func AppIds() string {
   var data strings.Builder
   for i, each := range Apps {
      if i >= 1 {
         data.WriteByte(' ')
      }
      data.WriteString(each.Id)
   }
   return data.String()
}

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

// doRequest logs the request method and URL, then sends the request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

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
   data = append(append(size, iv[:]...), data...)
   // 9. Return result base64 encoded
   return base64.StdEncoding.EncodeToString(data), nil
}

// newGetRequest builds a GET request with the given headers.
func newGetRequest(target string, header map[string]string) (*http.Request, error) {
   req, err := http.NewRequest("GET", target, nil)
   if err != nil {
      return nil, err
   }
   for key, val := range header {
      req.Header.Set(key, val)
   }
   return req, nil
}

// newPostRequest builds a POST request with the given headers and body.
func newPostRequest(target string, header map[string]string, body []byte) (*http.Request, error) {
   req, err := http.NewRequest("POST", target, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   for key, val := range header {
      req.Header.Set(key, val)
   }
   return req, nil
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

func readZipFile(f *zip.File) ([]byte, error) {
   rc, err := f.Open()
   if err != nil {
      return nil, err
   }
   defer rc.Close()
   return io.ReadAll(rc)
}

// util.go
