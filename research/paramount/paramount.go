package paramount

import "regexp"

var hexPattern = regexp.MustCompile(`([0-9a-f]{16})\x00\$`)

func ExtractHexBytes(data []byte) []string {
   matches := hexPattern.FindAllSubmatch(data, -1)
   var results []string
   for _, match := range matches {
      results = append(results, string(match[1]))
   }
   return results
}

var AppSecrets = []struct {
   version string
   url     string
   id      string
}{
   {
      url:     "https://apkmirror.com/apk/viacomcbs-streaming/paramount",
      version: "Paramount+ 16.8.0",
      id:      "com.cbs.app",
   },
   {
      url:     "https://apkmirror.com/apk/viacomcbs-streaming/paramount-4",
      version: "Paramount+ 16.8.0",
      id:      "com.cbs.ca",
   },
   {
      url:     "https://apkmirror.com/apk/cbs-interactive-inc/cbs",
      version: "CBS 15.6.0",
      id:      "com.cbs.tve",
   },
}
