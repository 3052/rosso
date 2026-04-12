package amc

import "encoding/json"

type AuthResponse struct {
   Success bool `json:"success"`
   Status  int  `json:"status"`
   Data    struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
      TokenType    string `json:"token_type"`
      ExpiresIn    int    `json:"expires_in"`
   } `json:"data"`
}

type ContentResponse struct {
   Success bool            `json:"success"`
   Status  int             `json:"status"`
   Data    json.RawMessage `json:"data"`
}

type PlaybackResponse struct {
   Success bool `json:"success"`
   Status  int  `json:"status"`
   Data    struct {
      PlaybackJsonData struct {
         VideoID string `json:"id"`
         Sources []struct {
            Codecs     string `json:"codecs"`
            Src        string `json:"src"`
            Type       string `json:"type"`
            KeySystems struct {
               ComWidevineAlpha struct {
                  LicenseURL string `json:"license_url"`
               } `json:"com.widevine.alpha"`
               ComMicrosoftPlayready struct {
                  LicenseURL string `json:"license_url"`
               } `json:"com.microsoft.playready"`
            } `json:"key_systems"`
         } `json:"sources"`
      } `json:"playbackJsonData"`
   } `json:"data"`
}

type PlaybackResult struct {
   Response PlaybackResponse
   BcovAuth string
}
