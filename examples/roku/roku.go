package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
)

type WidevineFolder string

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playback roku.Playback
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &playback, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: playback.LicenseWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

func (c *client) do_account_activation() error {
   account_token, err := roku.GetAccountToken(nil)
   if err != nil {
      return err
   }
   account_activation, err := roku.CreateAccountActivation(account_token)
   if err != nil {
      return err
   }
   fmt.Println(account_activation)
   return c.cache.Encode(account_activation, account_token)
}

func (c *client) do_activation_status() error {
   account_activation := &roku.AccountActivation{}
   account_token := &roku.AccountToken{}
   err := c.cache.Decode(account_activation, account_token)
   if err != nil {
      return err
   }
   activation_status, err := roku.GetActivationStatus(
      account_token, account_activation,
   )
   if err != nil {
      return err
   }
   return c.cache.Encode(activation_status)
}

func (c *client) do_roku_id() error {
   var status *roku.ActivationStatus
   if c.use_account.Set {
      status = &roku.ActivationStatus{}
      err := c.cache.Decode(&status)
      if err != nil {
         return err
      }
   }
   account_token, err := roku.GetAccountToken(status)
   if err != nil {
      return err
   }
   playback, err := roku.GetPlayback(account_token, c.roku_id.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&playback.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(account_token, manifest, playback)
}

type client struct {
   cache              maya.Cache
   flag               maya.FlagSet
   use_account        maya.Flag
   dash               maya.Flag
   roku_id            maya.Flag
   widevine           maya.Flag
   account_activation maya.Flag
   activation_status  maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/roku"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag.Add(&c.account_activation, "a", "account activation")
   c.flag.Add(&c.activation_status, "A", "activation status")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.roku_id, "r", "Roku ID")
   c.flag.Add(&c.use_account, "u", "use account")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(WidevineFolder(c.widevine.Value))
   }
   if c.account_activation.Set {
      return c.do_account_activation()
   }
   if c.activation_status.Set {
      return c.do_activation_status()
   }
   if c.roku_id.Set {
      return c.do_roku_id()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}
