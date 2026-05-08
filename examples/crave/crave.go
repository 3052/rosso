package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   profile_token := &crave.ProfileToken{}
   err := c.cache.Decode(profile_token)
   if err != nil {
      return err
   }
   media, err := crave.ParseMedia(c.address)
   if err != nil {
      return err
   }
   if media.FirstContent.Id == 0 {
      media, err = crave.GetMedia(media.Id)
      if err != nil {
         return err
      }
   }
   playback, err := crave.GetPlayback(profile_token, media)
   if err != nil {
      return err
   }
   stream, err := crave.GetStream(profile_token, playback)
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(stream)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, media, playback)
}

func (c *client) do_dash() error {
   if c.err != nil {
      return c.err
   }
   var (
      dash          maya.Dash
      playback      crave.Playback
      profile_token crave.ProfileToken
   )
   err := c.cache.Decode(&dash, &playback, &profile_token)
   if err != nil {
      return err
   }
   return dash.Download(&c.job, func(data []byte) ([]byte, error) {
      return crave.AcquireLicense(data, &profile_token, &playback)
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache    maya.Cache
   job      maya.Job
   address  string
   password string
   profile  string
   username string
   err      error
}

func (c *client) do_username_password() error {
   account_token, err := crave.PerformLogin(c.username, c.password)
   if err != nil {
      return err
   }
   profiles, err := crave.GetProfiles(account_token)
   if err != nil {
      return err
   }
   for i, profile := range profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return c.cache.Encode(account_token)
}

func (c *client) do_profile() error {
   account_token := &crave.AccountToken{}
   err := c.cache.Decode(account_token)
   if err != nil {
      return err
   }
   profile_token, err := crave.SwitchProfile(account_token, c.profile)
   if err != nil {
      return err
   }
   subs, err := crave.GetSubscriptions(profile_token)
   if err != nil {
      return err
   }
   for i, sub := range subs {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&sub)
   }
   return c.cache.Encode(profile_token)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/crave"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   password := maya.StringFlag(&c.password, "p", "password")
   profile := maya.StringFlag(&c.profile, "P", "profile")
   username := maya.StringFlag(&c.username, "u", "username")
   c.err = c.cache.Decode(&c.job)
   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")
   playReady := maya.StringFlag(&c.job.PlayReady, "PR", "PlayReady")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(c.job)
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   if profile.IsSet {
      return c.do_profile()
   }
   if address.IsSet {
      return c.do_address()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {username, password},
      {profile},
      {address},
      {dash},
   })
}
