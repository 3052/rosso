// common.go
package amazon

import (
   "html"
   "regexp"
   "strings"
)

// PageData holds the extracted form action URL and hidden input fields
// required to successfully pass Amazon's CSRF and session validations.
type PageData struct {
   ActionURL    string
   HiddenParams map[string]string
}

// ExtractPageData parses the HTML to find the login form action URL and hidden tokens.
func ExtractPageData(htmlStr string) *PageData {
   data := &PageData{
      HiddenParams: make(map[string]string),
   }

   // Extract form action URL
   reAction := regexp.MustCompile(`(?i)<form[^>]+name=["']signIn["'][^>]*action=["']([^"']+)["']`)
   actionMatch := reAction.FindStringSubmatch(htmlStr)
   if len(actionMatch) > 1 {
      action := html.UnescapeString(actionMatch[1])
      if strings.HasPrefix(action, "/") {
         action = "https://www.amazon.com" + action
      }
      data.ActionURL = action
   }

   // Extract hidden inputs
   reInput := regexp.MustCompile(`(?i)<input[^>]+type=["']?hidden["']?[^>]*>`)
   reName := regexp.MustCompile(`(?i)name=["']([^"']+)["']`)
   reValue := regexp.MustCompile(`(?i)value=["']([^"']*)["']`)

   matches := reInput.FindAllString(htmlStr, -1)
   for _, match := range matches {
      nameMatch := reName.FindStringSubmatch(match)
      valueMatch := reValue.FindStringSubmatch(match)

      if len(nameMatch) > 1 {
         name := nameMatch[1]
         val := ""
         if len(valueMatch) > 1 {
            val = html.UnescapeString(valueMatch[1])
         }
         data.HiddenParams[name] = val
      }
   }
   return data
}
