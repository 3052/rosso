package plex

import (
   "encoding/json"
   urlpkg "net/url"

   "41.neocities.org/maya"
)

type Subscription struct {
   Active   bool     `json:"active"`
   Status   string   `json:"status"`
   Features []string `json:"features"`
}

type Profile struct {
   AutoSelectAudio              bool `json:"autoSelectAudio"`
   DefaultAudioAccessibility    int  `json:"defaultAudioAccessibility"`
   AutoSelectSubtitle           int  `json:"autoSelectSubtitle"`
   DefaultSubtitleAccessibility int  `json:"defaultSubtitleAccessibility"`
   DefaultSubtitleForced        int  `json:"defaultSubtitleForced"`
   WatchedIndicator             int  `json:"watchedIndicator"`
   MediaReviewsVisibility       int  `json:"mediaReviewsVisibility"`
   MediaPostsVisibility         bool `json:"mediaPostsVisibility"`
}

type Service struct {
   Identifier string  `json:"identifier"`
   Endpoint   string  `json:"endpoint"`
   Token      *string `json:"token"`
   Status     string  `json:"status"`
   Secret     *string `json:"secret"`
}

type User struct {
   Id                   int          `json:"id"`
   Uuid                 string       `json:"uuid"`
   Username             string       `json:"username"`
   Title                string       `json:"title"`
   Email                string       `json:"email"`
   FriendlyName         string       `json:"friendlyName"`
   Confirmed            bool         `json:"confirmed"`
   JoinedAt             int          `json:"joinedAt"`
   EmailOnlyAuth        bool         `json:"emailOnlyAuth"`
   HasPassword          bool         `json:"hasPassword"`
   Protected            bool         `json:"protected"`
   Thumb                string       `json:"thumb"`
   AuthToken            string       `json:"authToken"`
   MailingListActive    bool         `json:"mailingListActive"`
   ScrobbleTypes        string       `json:"scrobbleTypes"`
   Subscription         Subscription `json:"subscription"`
   Restricted           bool         `json:"restricted"`
   Anonymous            bool         `json:"anonymous"`
   Home                 bool         `json:"home"`
   Guest                bool         `json:"guest"`
   HomeSize             int          `json:"homeSize"`
   HomeAdmin            bool         `json:"homeAdmin"`
   MaxHomeSize          int          `json:"maxHomeSize"`
   Profile              Profile      `json:"profile"`
   Services             []Service    `json:"services"`
   ExperimentalFeatures bool         `json:"experimentalFeatures"`
   TwoFactorEnabled     bool         `json:"twoFactorEnabled"`
   BackupCodesCreated   bool         `json:"backupCodesCreated"`
}

func CreateAnonymousUser() (*User, error) {
   endpoint := &urlpkg.URL{
      Scheme: "https",
      Host:   "plex.tv",
      Path:   "/api/v2/users/anonymous",
   }

   headers := map[string]string{
      "X-Plex-Client-Identifier": "!",
      "X-Plex-Product":           "Plex Mediaverse",
      "Accept":                   "application/json",
   }

   resp, err := maya.Post(endpoint, headers, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var user User
   if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
      return nil, err
   }

   return &user, nil
}
