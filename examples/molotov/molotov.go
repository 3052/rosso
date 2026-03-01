package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "flag"
   "log"
   "net/http"
   "path"
)

type saved_state struct {
   Asset *molotov.Asset
}

func (c *command) do_dash() error {
   var asset molotov.Asset
   err := c.cache.Get("Asset", &asset)
   if err != nil {
      return err
   }
   var dash molotov.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   c.job.Send = asset.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("molotov")
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
}

func (c *command) do_email_password() error {
   var login molotov.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Login", login)
}

func (c *command) do_address() error {
   var login molotov.Login
   err := c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   err = login.Refresh()
   if err != nil {
      return err
   }
   err = c.cache.Set("Login", login)
   if err != nil {
      return err
   }
   var media molotov.MediaId
   err = media.Parse(c.address)
   if err != nil {
      return err
   }
   view, err := login.ProgramView(&media)
   if err != nil {
      return err
   }
   asset, err := login.Asset(view)
   if err != nil {
      return err
   }
   err = c.cache.Set("Asset", asset)
   if err != nil {
      return err
   }
   dash, err := asset.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

type command struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.WidevineJob
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
