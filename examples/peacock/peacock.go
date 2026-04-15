package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "net/http"
   "path"
)

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, c.Playout.FetchWidevine)
}

func (c *client) do_address() error {
   token, err := peacock.FetchToken(c.IdSession)
   if err != nil {
      return err
   }
   c.Playout, err = token.FetchPlayout(path.Base(c.address))
   if err != nil {
      return err
   }
   endpoint, err := c.Playout.GetFastly()
   if err != nil {
      return err
   }
   dash, err := endpoint.ParseDash()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_email_password() error {
   var err error
   c.IdSession, err = peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

var cache maya.Cache

func main() {
   maya.SetProxy("", "*.m4s")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash      *maya.Dash
   IdSession *http.Cookie
   Playout   *peacock.Playout
   //----------------------
   Job maya.Job
   //----------------------
   email    string
   password string
   //----------------------
   address string
}

func (c *client) do() error {
   err := cache.Setup("rosso/peacock.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
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
   if dash.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}
