package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   video, err := kanopy.ParseVideo(c.address)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = c.Login.Video(video.Alias)
      if err != nil {
         return err
      }
   }
   membership, err := c.Login.Membership()
   if err != nil {
      return err
   }
   plays, err := c.Login.Plays(membership.DomainId, video.VideoId)
   if err != nil {
      return err
   }
   for _, caption := range plays.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   c.Manifest, err = plays.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = c.Manifest.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Login.Widevine(c.Manifest.DrmLicenseId, data)
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

type client struct {
   Dash     *kanopy.Dash
   Login    *kanopy.Login
   Manifest *kanopy.Manifest
   //-------------------------------
   Job maya.Job
   //-------------------------------
   email    string
   password string
   //-------------------------------
   address string
   //-------------------------------
   dash_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/kanopy.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //---------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[widevine] {
      return cache.Write(c)
   }
   if set[email] {
      if set[password] {
         return c.do_email_password()
      }
   }
   if set[address] {
      return with_cache(c.do_address)
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Login, err = kanopy.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
