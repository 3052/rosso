package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
   "os"
)

type client struct {
   cache          maya.Cache
   WidevineFolder maya.Flag[string]
   Email          maya.Flag[string] `depends:"Password"`
   Password       maya.Flag[string] `depends:"Email"`
   Address        maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/kanopy"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.WidevineFolder.Set {
      return c.cache.Encode(WidevineFolder(c.WidevineFolder.Value))
   }
   if c.Email.Set {
      if c.Password.Set {
         return c.do_email_password()
      }
   }
   if c.Address.Set {
      return c.do_address()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "kanopy", c)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   login, err := kanopy.LoginUser(c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

func (c *client) do_dash_id() error {
   var (
      login         kanopy.Login
      manifest      kanopy.Manifest
      maya_manifest maya.Manifest
      widevine      WidevineFolder
   )
   err := c.cache.Decode(&login, &manifest, &maya_manifest, &widevine)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return kanopy.CreateLicense(&login, &manifest, body)
   }
   return maya.DownloadDash(c.DashId.Value, &maya_manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}

type WidevineFolder string

func (c *client) do_address() error {
   login := &kanopy.Login{}
   err := c.cache.Decode(login)
   if err != nil {
      return err
   }
   video, err := kanopy.ParseVideo(c.Address.Value)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = kanopy.GetVideo(login, video.Alias)
      if err != nil {
         return err
      }
   }
   memberships, err := kanopy.GetMemberships(login)
   if err != nil {
      return err
   }
   play, err := kanopy.CreatePlay(login, &memberships[0], video)
   if err != nil {
      return err
   }
   for _, caption := range play.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   manifest, err := play.GetDash()
   if err != nil {
      return err
   }
   maya_manifest, err := maya.ListDash(&manifest.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, maya_manifest)
}
