package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Login.FetchWidevine(c.Manifest.DrmLicenseId, data)
      },
   )
}

func main() {
   maya.SetProxy("", "*.m4s")
   log.SetFlags(log.Ltime)
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
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
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

func (c *client) do_email_password() error {
   var err error
   c.Login, err = kanopy.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
func (c *client) do_address() error {
   video, err := kanopy.ParseVideo(c.address)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = c.Login.FetchVideo(video.Alias)
      if err != nil {
         return err
      }
   }
   membership, err := c.Login.FetchMembership()
   if err != nil {
      return err
   }
   plays, err := c.Login.FetchPlays(membership.DomainId, video.VideoId)
   if err != nil {
      return err
   }
   for _, caption := range plays.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   c.Manifest, err = plays.GetDash()
   if err != nil {
      return err
   }
   c.Dash, err = c.Manifest.FetchDash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
