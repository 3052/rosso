// examples/unext/unext.go
package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/unext"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine     maya.FlagString
   email        maya.FlagString
   password     maya.FlagString
   title_code   maya.FlagString
   episode_code maya.FlagString
   dash_id      maya.FlagString
   play_mode    maya.FlagString
   refresh      maya.FlagBool

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/unext/client"
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "refresh", Value: &c.refresh},
      {Name: "title-code", Value: &c.title_code},
      {Name: "episode-code", Value: &c.episode_code},
      {Name: "play-mode", Value: &c.play_mode, Needs: "episode-code", Usage: "caption dub"},
      {Name: "dash-id", Value: &c.dash_id},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.title_code != "" {
      return c.do_title_code()
   }
   if c.episode_code != "" {
      if c.play_mode != "" {
         return c.do_episode_code()
      }
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "unext")
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      playlist unext.PlaylistUrl
   )
   err := c.cache.Decode(&manifest, &playlist)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device: string(c.Widevine),
      Drm:    maya.DrmWidevine,
      License: func(challenge []byte) ([]byte, error) {
         licenseURL, err := playlist.WidevineLicenseURL()
         if err != nil {
            return nil, err
         }
         return unext.Step6GetLicense(licenseURL, playlist.PlayToken, challenge)
      },
   })
}

func (c *client) do_email_password() error {
   challengeID, err := unext.Step1GetChallenge()
   if err != nil {
      return err
   }
   postAuth, err := unext.Step2Login(string(c.email), string(c.password), challengeID)
   if err != nil {
      return err
   }
   authCode, err := unext.Step3GetAuthCode(postAuth)
   if err != nil {
      return err
   }
   tokens, err := unext.Step4GetToken(authCode)
   if err != nil {
      return err
   }
   return c.cache.Encode(tokens)
}

func (c *client) do_episode_code() error {
   tokens := &unext.TokenResponse{}
   err := c.cache.Decode(tokens)
   if err != nil {
      return err
   }
   playlist, err := unext.Step5GetPlaylist(
      tokens.AccessToken, string(c.episode_code), string(c.play_mode),
   )
   if err != nil {
      return err
   }
   mpdURL, err := playlist.MPDURL()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(mpdURL)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playlist)
}

func (c *client) do_refresh() error {
   var tokens unext.TokenResponse
   if err := c.cache.Decode(&tokens); err != nil {
      return err
   }
   if err := tokens.Refresh(); err != nil {
      return err
   }
   return c.cache.Encode(&tokens)
}

func (c *client) do_title_code() error {
   tokens := &unext.TokenResponse{}
   err := c.cache.Decode(tokens)
   if err != nil {
      return err
   }
   codes, err := unext.GetEpisodeCodes(tokens.AccessToken, string(c.title_code))
   if err != nil {
      return err
   }
   for _, code := range codes {
      log.Printf("episode: %s", code)
   }
   return nil
}
