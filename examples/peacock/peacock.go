package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "os"
   "path"
)

func (*client) CachePath() string {
   return "rosso/examples/peacock/client"
}

type client struct {
   Widevine maya.FlagString
   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString

   cache maya.Cache
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
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
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
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "peacock")
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

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   id_session, err := peacock.FetchIdSession(string(c.email), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(id_session)
}

func (c *client) do_address() error {
   id_session := &peacock.Cookie{}
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
