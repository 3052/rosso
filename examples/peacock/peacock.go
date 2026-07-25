// examples/peacock/peacock.go
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
   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   otp      maya.FlagString

   otpToken maya.FlagString

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
      {Name: "email", Value: &c.email},
      {Name: "otp", Value: &c.otp, Needs: "email"},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.email != "" && c.otp == "" {
      return c.do_email()
   }
   if c.otp != "" {
      return c.do_otp()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "peacock")
}

func (c *client) do_address() error {
   id_session := &peacock.IdSession{}
   err := c.cache.Decode(id_session)
   if err != nil {
      return err
   }
   token, err := peacock.FetchToken(id_session)
   if err != nil {
      return err
   }
   playout, err := token.FetchPlayout(
      path.Base(string(c.address)),
   )
   if err != nil {
      return err
   }
   endpoint, err := playout.GetFastly()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(endpoint)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playout)
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playout  peacock.Playout
   )
   err := c.cache.Decode(&manifest, &playout)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: playout.FetchWidevine,
   })
}

func (c *client) do_email() error {
   var session peacock.IdSession
   token, err := session.InitiateOTP(string(c.email))
   if err != nil {
      return err
   }
   c.otpToken = maya.FlagString(token)
   return c.cache.Encode(c)
}

func (c *client) do_otp() error {
   var session peacock.IdSession
   session.Cookie = nil // will be set by VerifyOTP
   err := session.VerifyOTP(string(c.otpToken), string(c.otp))
   if err != nil {
      return err
   }
   c.otpToken = "" // token consumed
   if err := c.cache.Encode(&session); err != nil {
      return err
   }
   return c.cache.Encode(c)
}
