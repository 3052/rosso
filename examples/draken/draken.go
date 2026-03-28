package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/draken"
   "log"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/draken.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //----------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return with_cache(c.do_address)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Playback.Widevine(c.Login.Token, data)
      },
   )
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_email_password() error {
   var err error
   c.Login, err = draken.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   viewer, err := draken.FetchViewer(path.Base(c.address))
   if err != nil {
      return err
   }
   entitlement, err := viewer.Entitlement(c.Login.Token)
   if err != nil {
      return err
   }
   c.Playback, err = viewer.Playback(c.Login.Token, entitlement.Token)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Dash     *draken.Dash
   Login    *draken.Login
   Playback *draken.Playback
   //-----------------------
   Job maya.Job
   //-----------------------
   email    string
   password string
   //-----------------------
   address string
   //-----------------------
   dash_id string
}
