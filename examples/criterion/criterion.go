package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "log"
   "os"
   "path"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine maya.FlagString

   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/criterion"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.ErrNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
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
   return flags.Usage(os.Stderr, "criterion")
}

///

func (c *client) do_email_password() error {
   token, err := criterion.FetchToken(c.Email.Value, c.Password.Value)
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
      token.AccessToken, path.Base(c.Address.Value),
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

func (c *client) do_dash() error {
   var (
      file     criterion.File
      manifest maya.Manifest
      widevine WidevineFolder
   )
   err := c.cache.Decode(&file, &manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: file.FetchWidevine,
   })
}
