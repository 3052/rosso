// get_csrf_token.go
package amazon

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "regexp"
   "strings"
)

func GetCSRFToken(client *http.Client, endpoint string) (string, error) {
   req, err := http.NewRequest(http.MethodGet, endpoint, nil)
   if err != nil {
      return "", err
   }

   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   bodyStr := string(body)

   if strings.Contains(bodyStr, `input type="hidden" name="appAction" value="SIGNIN"`) {
      return "", errors.New("cookies are signed out, cannot get ontv CSRF token")
   }

   re := regexp.MustCompile(`<script type="text/template">(.+?)</script>`)
   matches := re.FindAllStringSubmatch(bodyStr, -1)

   for _, match := range matches {
      if len(match) > 1 {
         var data struct {
            Props struct {
               CodeEntry struct {
                  Token string `json:"token"`
               } `json:"codeEntry"`
            } `json:"props"`
         }
         if err := json.Unmarshal([]byte(match[1]), &data); err == nil && data.Props.CodeEntry.Token != "" {
            return data.Props.CodeEntry.Token, nil
         }
      }
   }

   return "", errors.New("unable to get ontv CSRF token")
}
