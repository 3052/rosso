package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/maya/unext"
   "log"
   "net/http"
   "net/http/cookiejar"
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
   Widevine maya.FlagString
   email    maya.FlagString
   password maya.FlagString
   title    maya.FlagString
   episode  maya.FlagString
   dash     maya.FlagString

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
      {Name: "title-code", Value: &c.title},
      {Name: "episode-code", Value: &c.episode},
      {Name: "dash", Value: &c.dash},
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
   if c.title != "" {
      return c.do_title_code()
   }
   if c.episode != "" {
      return c.do_episode_code()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "unext")
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playlist unext.PlaylistUrl
   )
   err := c.cache.Decode(&manifest, &playlist)
   if err != nil {
      return err
   }
   httpClient := &http.Client{}
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device: string(c.Widevine),
      Drm:    maya.DrmWidevine,
      License: func(challenge []byte) ([]byte, error) {
         licenseURL, err := playlist.WidevineLicenseURL()
         if err != nil {
            return nil, err
         }
         return unext.Step6GetLicense(httpClient, licenseURL, playlist.PlayToken, challenge)
      },
   })
}

func (c *client) do_email_password() error {
   verifier, challenge, err := unext.PkcePair()
   if err != nil {
      return err
   }
   state, err := unext.GenerateRandomString(43)
   if err != nil {
      return err
   }
   nonce, err := unext.GenerateRandomString(43)
   if err != nil {
      return err
   }

   jar, err := cookiejar.New(nil)
   if err != nil {
      return err
   }
   httpClient := &http.Client{
      Jar: jar,
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   challengeID, err := unext.Step1GetChallenge(httpClient, state, nonce)
   if err != nil {
      return err
   }
   postAuth, err := unext.Step2Login(httpClient, string(c.email), string(c.password), challengeID)
   if err != nil {
      return err
   }
   authCode, err := unext.Step3GetAuthCode(httpClient, postAuth, challenge)
   if err != nil {
      return err
   }
   tokens, err := unext.Step4GetToken(httpClient, authCode, verifier)
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
   httpClient := &http.Client{}

   // Assumes Step5GetPlaylist accepts the episode code as a parameter
   playlist, err := unext.Step5GetPlaylist(httpClient, tokens.AccessToken, string(c.episode))
   if err != nil {
      return err
   }
   mpdURL, err := playlist.MPDURL()
   if err != nil {
      return err
   }

   resp, err := httpClient.Get(mpdURL.String())
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   manifest, err := maya.ListDash(resp.Body)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playlist)
}

func (c *client) do_title_code() error {
   tokens := &unext.TokenResponse{}
   err := c.cache.Decode(tokens)
   if err != nil {
      return err
   }
   httpClient := &http.Client{}
   codes, err := unext.GetEpisodeCodes(httpClient, tokens.AccessToken, string(c.title))
   if err != nil {
      return err
   }
   for _, code := range codes {
      log.Printf("episode: %s", code)
   }
   return nil
}
