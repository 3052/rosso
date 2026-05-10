package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "path"
)

func (c *client) do_dash() error {
   var (
      dash    maya.Dash
      playout peacock.Playout
   )
   err := c.cache.Decode(&c.job, &dash, &playout)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, playout.FetchWidevine)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   job      maya.Job
   password string
}

func (c *client) do_email_password() error {
   id_session, err := peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(id_session)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/peacock"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
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
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},

      {address},
      {dash},
   })
}

///

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
   dash, err := maya.ListDash(endpoint)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, playout)
}
