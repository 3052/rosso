package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "fmt"
   "log"
   "path"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/criterion"); err != nil {
      return err
   }
   c.widevine = c.flag.AddValue("w", "Widevine")
   c.flag = append(c.flag, nil)
   c.email = c.flag.AddValue("e", "email")
   c.password = c.flag.AddValue("p", "password")
   c.flag = append(c.flag, nil)
   c.address = c.flag.AddValue("a", "address")
   c.dash = c.flag.AddValue("d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine(c.widevine.Value))
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

type widevine string

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
   return c.cache.Encode(manifest, file, token)
}

func (c *client) do_dash() error {
   var (
      file     criterion.File
      manifest maya.Manifest
      device   widevine
   )
   err := c.cache.Decode(&file, &manifest, &device)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: file.FetchWidevine,
   })
}

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

type client struct {
   cache    maya.Cache
   address  *maya.Flag
   dash     *maya.Flag
   email    *maya.Flag
   password *maya.Flag
   widevine *maya.Flag
   flag     maya.FlagSet
}
