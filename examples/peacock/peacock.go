package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "path"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/peacock.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
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
      return c.run(c.do_address)
   }
   if dash.IsSet {
      return c.run(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
}

func (c *client) do_email_password() error {
   var err error
   c.IdSession, err = peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, c.Playout.FetchWidevine)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
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
   c.Dash, err = maya.ListDash(endpoint.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   // cache
   Dash      *maya.Dash
   IdSession string
   Job       maya.Job
   Playout   *peacock.Playout
   // flags
   address  string
   email    string
   password string
   // state
   cache_err error
}
