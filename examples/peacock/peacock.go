package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "os"
   "path"
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
   token    maya.FlagBool
   address  maya.FlagString
   dash_id  maya.FlagString
   vcodec   maya.FlagString

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/peacock/client"
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
      {Name: "token", Value: &c.token},
      {Name: "address", Value: &c.address, Needs: "vcodec"},
      {Name: "vcodec", Value: &c.vcodec, Needs: "address", Usage: "H264 H265"},
      {Name: "dash-id", Value: &c.dash_id},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.email != "" {
      return c.do_email()
   }
   if flags.IsSet(&c.token) {
      return c.do_token()
   }
   if c.address != "" {
      if c.vcodec != "" {
         return c.do_address()
      }
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "peacock")
}

func (c *client) do_address() error {
   var token peacock.TokenResponse
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   contentID := path.Base(string(c.address))
   playout, err := peacock.PlayoutVod(&token, contentID, string(c.vcodec), "WIDEVINE")
   if err != nil {
      return err
   }
   endpoint, err := playout.Fastly()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(endpoint)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playout)
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      playout  peacock.PlayoutVodResponse
   )
   err := c.cache.Decode(&manifest, &playout)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: playout.AcquireLicense,
   })
}

func (c *client) do_email() error {
   signIn, err := peacock.SignIn(&peacock.SignInParams{
      UserIdentifier: string(c.email),
      Password:       string(c.password),
      RememberMe:     true,
   })
   if err != nil {
      return err
   }
   authResp, err := peacock.OAuthAuthorize(signIn.Cookie)
   if err != nil {
      return err
   }
   return c.cache.Encode(authResp, signIn)
}

func (c *client) do_token() error {
   var (
      authResp peacock.OAuthAuthorizeResponse
      signIn   peacock.SignInResponse
   )
   err := c.cache.Decode(&authResp, &signIn)
   if err != nil {
      return err
   }
   token, err := peacock.ExchangeToken(&authResp, &signIn)
   if err != nil {
      return err
   }
   return c.cache.Encode(token)
}
