package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
)

func (c *client) do_address() error {
   var err error
   c.Device, err = c.Device.TokenRefresh()
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
   maya.SetProxy("", "*.mp4", "*.mp4a")
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

type client struct {
   Dash     *maya.Dash
   Playlist *hulu.Playlist
   Device   *hulu.Device
   //--------------------
   Job maya.Job
   //--------------------
   email    string
   password string
   //--------------------
   address string
}

var cache maya.Cache

func (c *client) do() error {
   err := cache.Setup("rosso/hulu.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "P", "PlayReady")
   //-------------------------------------------------------------
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
   if playReady.IsSet {
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
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Playlist.FetchPlayReady)
}
