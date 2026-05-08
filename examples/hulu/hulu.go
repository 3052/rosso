package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hulu"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   playReady := maya.StringFlag(&c.job.PlayReady, "P", "PlayReady")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
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
      {playReady},
      {email, password},
      {address},
      {dash},
   })
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   job      maya.Job
   password string
}

///

func (c *client) do_email_password() error {
   var err error
   c.Device, err = hulu.FetchDevice(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_address() error {
   err := c.Device.TokenRefresh()
   if err != nil {
      return err
   }
   deep_link, err := c.Device.DeepLink(hulu.ParseId(c.address))
   if err != nil {
      return err
   }
   c.Playlist, err = c.Device.Playlist(deep_link.EabId)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.Playlist.GetManifest)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Playlist.FetchPlayReady)
}
