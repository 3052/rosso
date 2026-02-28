package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("criterion")
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
   var token criterion.Token
   err := token.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Token", token)
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
      return "", path.Ext(req.URL.Path) != ".mp4"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_dash() error {
   var media criterion.MediaFile
   err := c.cache.Get("MediaFile", &media)
   if err != nil {
      return err
   }
   c.job.Send = media.Widevine
   var dash criterion.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}
func (c *command) do_address() error {
   var token criterion.Token
   err := c.cache.Get("Token", &token)
   if err != nil {
      return err
   }
   err = token.Refresh()
   if err != nil {
      return err
   }
   err = c.cache.Set("Token", token)
   if err != nil {
      return err
   }
   item, err := token.Item(path.Base(c.address))
   if err != nil {
      return err
   }
   files, err := token.Files(item)
   if err != nil {
      return err
   }
   media_file, err := files.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("MediaFile", media_file)
   if err != nil {
      return err
   }
   dash, err := media_file.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
