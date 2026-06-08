package amazon

import (
   "net/url"
   "regexp"
)

var (
   // Finds the first form named "signIn" and grabs its body
   signInFormRegex  = regexp.MustCompile(`(?is)<form[^>]+name=["']signIn["'][^>]*>(.*?)</form>`)
   formTagRegex     = regexp.MustCompile(`(?i)<form[^>]+name=["']signIn["'][^>]*>`)
   actionRegex      = regexp.MustCompile(`(?i)action=["']([^"']+)["']`)
   hiddenInputRegex = regexp.MustCompile(`(?i)<input[^>]+type=["']hidden["'][^>]*>`)
   nameRegex        = regexp.MustCompile(`(?i)name=["']([^"']+)["']`)
   valueRegex       = regexp.MustCompile(`(?i)value=["']([^"']*)["']`)
)

func extractFormActionAndHiddenInputs(html, baseUrl string) (string, map[string]string) {
   hiddenParams := make(map[string]string)
   actionUrl := ""

   if formTag := formTagRegex.FindString(html); formTag != "" {
      if m := actionRegex.FindStringSubmatch(formTag); len(m) > 1 {
         actionUrl = m[1]
      }
   }

   if actionUrl != "" {
      if base, err := url.Parse(baseUrl); err == nil {
         if rel, err := url.Parse(actionUrl); err == nil {
            actionUrl = base.ResolveReference(rel).String()
         }
      }
   }

   if m := signInFormRegex.FindStringSubmatch(html); len(m) > 1 {
      formBody := m[1]
      inputs := hiddenInputRegex.FindAllString(formBody, -1)
      for _, input := range inputs {
         var name, value string
         if nm := nameRegex.FindStringSubmatch(input); len(nm) > 1 {
            name = nm[1]
         }
         if vm := valueRegex.FindStringSubmatch(input); len(vm) > 1 {
            value = vm[1]
         }
         if name != "" {
            hiddenParams[name] = value
         }
      }
   }
   return actionUrl, hiddenParams
}
