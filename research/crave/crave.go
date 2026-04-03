package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   media_id, err := crave.ParseMediaId(c.address)
   if err != nil {
      return err
   }
   c.Media, err = crave.FetchMedia(media_id)
   if err != nil {
      return err
   }
   c.ContentPackage, err = c.Media.FetchContentPackage()
   if err != nil {
      return err
   }
   manifest, err := c.ContentPackage.ManifestWidevine(
      c.Media.FirstContent.Id, c.Account.AccessToken,
   )
   if err != nil {
      return err
   }
   c.Dash, err = manifest.FetchDash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

var cache maya.Cache

type client struct {
   Account        *crave.Account
   ContentPackage *crave.ContentPackage
   Dash           *crave.Dash
   Media          *crave.Media
   //--------------------
   Job maya.Job
   //--------------------
   username string
   password string
   //--------------------
   profile string
   //--------------------
   address string
   //--------------------
   dash_id string
}

func main() {
   // MP4 need proxy so just use VPN
   maya.SetProxy("", "*.m4v")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_profile() error {
   err := c.Account.Login(c.profile)
   if err != nil {
      return err
   }
   subs, err := c.Account.FetchSubscriptions()
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

func (c *client) do() error {
   err := cache.Setup("rosso/crave.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //-----------------------------------------------------------
   username := maya.StringFlag(&c.username, "u", "username")
   password := maya.StringFlag(&c.password, "p", "password")
   //-----------------------------------------------------------
   profile := maya.StringFlag(&c.profile, "P", "profile")
   //-----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //-----------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
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
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {username, password},
      {profile},
      {address},
      {dash_id},
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

func (c *client) do_dash_id() error {
   fetch := func(data []byte) ([]byte, error) {
      return c.ContentPackage.LicenseWidevine(
         c.Media.FirstContent.Id, c.Account.AccessToken, data,
      )
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, fetch)
}
