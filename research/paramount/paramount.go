package paramount

import (
   "archive/zip"
   "bytes"
   "io"
   "regexp"
   "strings"
)

var hexPattern = regexp.MustCompile(`\x00\x10([0-9a-f]{16})\x00`)

// ExtractDexHexBytes returns a set (map) of unique 16-character hex strings
// found in .dex files
func ExtractDexHexBytes(zipData []byte) (map[string]struct{}, error) {
   results := make(map[string]struct{})
   reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
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
type version struct {
   version    string
   app_secret string
}

var Apps = []struct {
   url      string
   id       string
   versions []version
}{
   {
      url: "https://apkmirror.com/apk/viacomcbs-streaming/paramount",
      id:  "com.cbs.app",
      versions: []version{
         {
            version:    "Paramount+ 16.8.0",
            app_secret: "7081400bd4143bf3",
         },
      },
   },
   {
      url: "https://apkmirror.com/apk/viacomcbs-streaming/paramount-4",
      id:  "com.cbs.ca",
      versions: []version{
         {
            version:    "Paramount+ 16.8.0",
            app_secret: "1c5d27627d71b420",
         },
      },
   },
   {
      url: "https://apkmirror.com/apk/cbs-interactive-inc/cbs",
      id:  "com.cbs.tve",
      versions: []version{
         {
            version:    "CBS 15.6.0",
            app_secret: "cef32931dc01412e",
         },
      },
   },
}
