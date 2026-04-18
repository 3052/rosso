package rtbf

import (
   "errors"
   "net/url"
)

type Identity struct {
   ErrorMessage string
   IdToken      string `json:"id_token"`
}

type Session struct {
   SessionToken string
}

// hard coded in JavaScript
const api_key = "4_Ml_fJ47GnBAW6FrPzMxh0w"

type Entitlement struct {
   AssetId   string
   Formats   []Format
   Message   string
   PlayToken string
}

type Format struct {
   Format       string
   MediaLocator string // MPD
}

func (e *Entitlement) GetDash() (*Format, error) {
   for _, format_data := range e.Formats {
      if format_data.Format == "DASH" {
         return &format_data, nil
      }
   }
   return nil, errors.New("DASH format not found")
}

func (f *Format) GetManifest() (*url.URL, error) {
   return url.Parse(f.MediaLocator)
}

func GetPath(urlData string) (string, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   if url_parse.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return url_parse.Path, nil
}

type Account struct {
   ErrorMessage string
   SessionInfo  struct {
      CookieValue string
   }
}
