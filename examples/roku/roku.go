package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
)

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Playback.GetWidevineLicense)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache       maya.Cache
   dash        string
   job         maya.Job
   roku_id     string
   use_account *maya.Flag
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
   return c.cache.Encode(activation_status)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/roku"); err != nil {
      return err
   }
   account_activation := maya.BoolFlag("a", "account activation")
   activation_status := maya.BoolFlag("A", "activation status")
   c.use_account = maya.BoolFlag("u", "use account")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   roku_id := maya.StringFlag(&c.roku_id, "r", "Roku ID")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
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
   return maya.PrintFlags([][]*maya.Flag{
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
   dash, err := maya.ListDash(&playback.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(account_token, dash, playback)
}
