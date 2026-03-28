package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
)

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
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
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
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Login, err = molotov.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   url, err := molotov.ParseUrl(c.address)
   if err != nil {
      return err
   }
   err = c.Login.Refresh()
   if err != nil {
      return err
   }
   program, err := url.FetchProgram(c.Login.Auth.AccessToken)
   if err != nil {
      return err
   }
   c.Asset, err = program.Asset(c.Login.Auth.AccessToken)
   if err != nil {
      return err
   }
   c.Dash, err = c.Asset.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Asset.Widevine,
   )
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   Asset *molotov.Asset
   Dash  *molotov.Dash
   Login *molotov.Login
   //------------------
   Job maya.Job
   //-------------
   email    string
   password string
   //-------------
   address string
   //------------
   dash_id string
}
