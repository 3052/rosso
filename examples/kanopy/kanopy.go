package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
)

func main() {
   maya.SetProxy("", "*.m4s")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, func(data []byte) ([]byte, error) {
      return c.Session.GetWidevine(c.Manifest, data)
   })
}

func (c *client) do() error {
   err := cache.Setup("rosso/kanopy.xml")
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
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Session, err = kanopy.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Session  *kanopy.Session
   Manifest *kanopy.Manifest
   //-------------------------------
   Dash *maya.Dash
   Job  maya.Job
   //-------------------------------
   email    string
   password string
   //-------------------------------
   address string
}

func (c *client) do_address() error {
   video, err := kanopy.ParseVideo(c.address)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = c.Session.GetVideo(video.Alias)
      if err != nil {
         return err
      }
   }
   memberships, err := c.Session.GetMemberships()
   if err != nil {
      return err
   }
   play, err := c.Session.CreatePlay(&memberships[0], video)
   if err != nil {
      return err
   }
   for _, caption := range play.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   c.Manifest, err = play.DashManifest()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.Manifest.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
