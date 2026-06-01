// utils.go
package amazon

import (
   "crypto/rand"
   "fmt"
   "html"
   "net/http"
   "regexp"
)

type Session struct {
   Client           *http.Client
   Email            string
   Password         string
   VideoID          string
   PlaybackEnvelope string
   TargetTitleID    string
   DeviceID         string
}

func GenerateUUID() string {
   b := make([]byte, 16)
   _, _ = rand.Read(b)
   b[6] = (b[6] & 0x0f) | 0x40
   b[8] = (b[8] & 0x3f) | 0x80
   return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func ExtractForm(htmlStr string, formName string) (string, map[string]string) {
   inputs := make(map[string]string)
   formPattern := regexp.MustCompile(`(?is)<form[^>]*name=["']?` + formName + `["']?[^>]*>(.*?)</form>`)
   match := formPattern.FindStringSubmatch(htmlStr)
   if len(match) == 0 {
      formPattern = regexp.MustCompile(`(?is)<form[^>]*id=["']?` + formName + `["']?[^>]*>(.*?)</form>`)
      match = formPattern.FindStringSubmatch(htmlStr)
      if len(match) == 0 {
         return "", inputs
      }
   }

   formTagPattern := regexp.MustCompile(`(?is)<form[^>]*>`)
   formTagMatch := formTagPattern.FindString(match[0])

   actionPattern := regexp.MustCompile(`(?is)action=["']?([^"'\s>]+)`)
   actionMatch := actionPattern.FindStringSubmatch(formTagMatch)
   action := ""
   if len(actionMatch) > 1 {
      action = html.UnescapeString(actionMatch[1])
   }

   inputPattern := regexp.MustCompile(`(?is)<input\s+([^>]+)>`)
   inputMatches := inputPattern.FindAllStringSubmatch(match[1], -1)

   namePattern := regexp.MustCompile(`(?is)name=["']?([^"'\s>]+)`)
   valuePattern := regexp.MustCompile(`(?is)value=["']([^"']*)["']`)
   valuePatternAlt := regexp.MustCompile(`(?is)value=([^"'\s>]+)`)

   for _, m := range inputMatches {
      attrs := m[1]
      nameMatch := namePattern.FindStringSubmatch(attrs)
      if len(nameMatch) > 1 {
         name := html.UnescapeString(nameMatch[1])
         value := ""
         valMatch := valuePattern.FindStringSubmatch(attrs)
         if len(valMatch) > 1 {
            value = html.UnescapeString(valMatch[1])
         } else {
            valMatchAlt := valuePatternAlt.FindStringSubmatch(attrs)
            if len(valMatchAlt) > 1 {
               value = html.UnescapeString(valMatchAlt[1])
            }
         }
         inputs[name] = value
      }
   }
   return action, inputs
}
