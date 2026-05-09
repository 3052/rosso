package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
)

func (c *client) do_dash() error {
   var (
      asset molotov.Asset
      dash  maya.Dash
   )
   err := c.cache.Decode(&c.job, &asset, &dash)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, asset.FetchWidevine)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   job      maya.Job
   password string
}

func (c *client) do_email_password() error {
   auth, err := molotov.FetchAuth(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(auth)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/molotov"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return c.do_address()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},

      {address},
      {dash},
   })
}

///

func (c *client) do_address() error {
   var auth molotov.Auth
   err := c.cache.Decode(&auth)
   if err != nil {
      return err
   }
   program, err := molotov.ParseProgram(c.address)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   play, err := auth.FetchPlay(program)
   if err != nil {
      return err
   }
   asset, err := auth.FetchAsset(play)
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(asset.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(asset, auth, dash)
}
