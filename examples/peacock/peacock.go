package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "path"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playout  peacock.Playout
      widevine widevine_folder
   )
   err := c.cache.Decode(&manifest, &playout, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(widevine),
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
   id_session, err := peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(id_session)
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   flag     maya.FlagSet
   password string
   widevine string
}

type widevine_folder string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/peacock"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   email := c.flag.String(&c.email, "e", "email")
   password := c.flag.String(&c.password, "p", "password")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(widevine_folder(c.widevine))
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return c.do_address()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
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
   playout, err := token.FetchPlayout(path.Base(c.address))
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
