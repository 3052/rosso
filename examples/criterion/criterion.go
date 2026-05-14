package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "fmt"
   "log"
   "path"
)

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   address  maya.Flag
   dash     maya.Flag
   email    maya.Flag
   password maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/criterion"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.email, "e", "email")
   c.flag.AddValue(&c.password, "p", "password")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine_device(c.widevine.Value))
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if c.address.Set {
      return c.do_address()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_dash() error {
   var (
      device   widevine_device
      file     criterion.File
      manifest maya.Manifest
   )
   err := c.cache.Decode(&device, &file, &manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: file.FetchWidevine,
   })
}

type widevine_device string

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   token, err := criterion.FetchToken(c.email.Value, c.password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(token)
}

func (c *client) do_address() error {
   var token criterion.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   err = token.Refresh()
   if err != nil {
      return err
   }
   files_href, err := criterion.FetchFilesHref(
      token.AccessToken, path.Base(c.address.Value),
   )
   if err != nil {
      return err
   }
   files, err := criterion.FetchFiles(token.AccessToken, files_href)
   if err != nil {
      return err
   }
   file, err := criterion.GetDash(files)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&file.Links.Source.Href.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(file, manifest, token)
}
