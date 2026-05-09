package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
)

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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/molotov"); err != nil {
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

func (c *client) do_email_password() error {
   var err error
   c.Auth, err = molotov.FetchAuth(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_address() error {
   program, err := molotov.ParseProgram(c.address)
   if err != nil {
      return err
   }
   err = c.Auth.Refresh()
   if err != nil {
      return err
   }
   play, err := c.Auth.FetchPlay(program)
   if err != nil {
      return err
   }
   c.Asset, err = c.Auth.FetchAsset(play)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.Asset.GetManifest)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Asset.FetchWidevine)
}
