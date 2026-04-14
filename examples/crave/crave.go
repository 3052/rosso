package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
)

func (c *client) do_profile() error {
   err := c.Account.Login(c.profile)
   if err != nil {
      return err
   }
   subs, err := crave.FetchSubscriptions(c.Account.AccessToken)
   if err != nil {
      return err
   }
   for i, sub := range subs {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&sub)
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

type client struct {
   Account        *crave.Account
   ContentPackage *crave.ContentPackage
   Media          *crave.Media
   Dash           *maya.Dash
   //--------------------
   Job maya.Job
   //--------------------
   username string
   password string
   //--------------------
   profile string
   //--------------------
   address string
}

func (c *client) do() error {
   err := cache.Setup("rosso/crave.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "PR", "PlayReady")
   //-----------------------------------------------------------
   username := maya.StringFlag(&c.username, "u", "username")
   password := maya.StringFlag(&c.password, "p", "password")
   //-----------------------------------------------------------
   profile := maya.StringFlag(&c.profile, "P", "profile")
   //-----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //-----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   maya.SetProxy("", "*.m4v")
   if playReady.IsSet {
      return cache.Write(c)
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   if profile.IsSet {
      return with_cache(c.do_profile)
   }
   if address.IsSet {
      return with_cache(c.do_address)
   }
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {username, password},
      {profile},
      {address},
      {dash},
   })
}

func (c *client) do_username_password() error {
   var err error
   c.Account, err = crave.Login(c.username, c.password)
   if err != nil {
      return err
   }
   profiles, err := c.Account.FetchProfiles()
   if err != nil {
      return err
   }
   for i, profile := range profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(profile)
   }
   return cache.Write(c)
}

func (c *client) do_dash() error {
   fetch := func(data []byte) ([]byte, error) {
      return c.ContentPackage.LicensePlayReady(
         c.Media.FirstContent.Id, c.Account.AccessToken, data,
      )
   }
   return c.Dash.Download(&c.Job, fetch)
}

func (c *client) do_address() error {
   var err error
   c.Media, err = crave.ParseMedia(c.address)
   if err != nil {
      return err
   }
   if c.Media.FirstContent.Id == 0 {
      c.Media, err = crave.FetchMedia(c.Media.Id)
      if err != nil {
         return err
      }
   }
   c.ContentPackage, err = crave.FetchContentPackage(
      c.Account.AccessToken, c.Media.FirstContent.Id,
   )
   if err != nil {
      return err
   }
   manifest, err := c.ContentPackage.ManifestPlayReady(
      c.Media.FirstContent.Id, c.Account.AccessToken,
   )
   if err != nil {
      return err
   }
   dash, err := crave.ParseDash(manifest.Playback)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

var cache maya.Cache
