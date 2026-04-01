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
   media, err := crave.FetchMedia(media_id)
   if err != nil {
      return err
   }
   content, err := media.FetchContentPackage()
   if err != nil {
      return err
   }
   manifest, err := content.FetchManifest(
      media.FirstContent.Id, c.Account.AccessToken,
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

type client struct {
   Account *crave.Account
   Dash    *crave.Dash
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

func (c *client) do_profile() error {
   err := c.Account.Login(c.profile)
   if err != nil {
      return err
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
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {username, password},
      {profile},
      {address},
   })
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
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
   return nil
}
