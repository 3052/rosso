package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playback roku.Playback
      widevine device
   )
   err := c.cache.Decode(&manifest, &playback, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
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

type client struct {
   cache       maya.Cache
   dash        string
   flag        maya.FlagSet
   roku_id     string
   use_account *maya.Flag
   widevine    string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/roku"); err != nil {
      return err
   }
   account_activation := c.flag.Bool("a", "account activation")
   activation_status := c.flag.Bool("A", "activation status")
   c.use_account = c.flag.Bool("u", "use account")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   roku_id := c.flag.String(&c.roku_id, "r", "Roku ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(device(c.widevine))
   }
   if account_activation.IsSet {
      return c.do_account_activation()
   }
   if activation_status.IsSet {
      return c.do_activation_status()
   }
   if roku_id.IsSet {
      return c.do_roku_id()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {account_activation},
      {activation_status},
      {roku_id, c.use_account},
      {dash},
   })
}

func (c *client) do_roku_id() error {
   var status *roku.ActivationStatus
   if c.use_account.IsSet {
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
   playback, err := roku.GetPlayback(account_token, c.roku_id)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&playback.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(account_token, manifest, playback)
}
