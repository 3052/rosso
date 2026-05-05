package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/hulu.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "P", "PlayReady")
   //-------------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
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
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
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

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Playlist.FetchPlayReady)
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
   return cache.Write(c)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   var err error
   c.Device, err = hulu.FetchDevice(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

var cache maya.Cache

type client struct {
   // cache
   Dash     *maya.Dash
   Device   *hulu.Device
   Job      maya.Job
   Playlist *hulu.Playlist
   // flags
   address  string
   email    string
   password string
   // state
   cache_err error
}
