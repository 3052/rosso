package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
)

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Asset.FetchWidevine)
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
   Asset *molotov.Asset
   Dash  *maya.Dash
   Auth  *molotov.Auth
   //------------------
   Job maya.Job
   //-------------
   email    string
   password string
   //-------------
   address string
}

func (c *client) do() error {
   err := cache.Setup("rosso/molotov.xml")
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
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
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
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Auth, err = molotov.FetchAuth(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   program, err := molotov.ParseProgram(c.address)
   if err != nil {
      return err
   }
   c.Auth, err = c.Auth.Refresh()
   if err != nil {
      return err
   }
   play, err := c.Auth.FetchPlay(program)
   if err != nil {
      return err
   }
   c.Asset, err = c.Auth.FetchAsset(play)
   if err != nil {
      return err
   }
   dash, err := c.Asset.ParseDash()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
