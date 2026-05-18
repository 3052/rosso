package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "log"
   "os"
   "path"
)

type client struct {
   cache          maya.Cache
   WidevineFolder maya.Flag[string]
   Email          maya.Flag[string] `depends:"Password"`
   Password       maya.Flag[string] `depends:"Email"`
   Address        maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/criterion"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.WidevineFolder.Set {
      return c.cache.Encode(WidevineFolder(c.WidevineFolder.Value))
   }
   if c.Email.Set {
      if c.Password.Set {
         return c.do_email_password()
      }
   }
   if c.Address.Set {
      return c.do_address()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "criterion", c)
}

func (c *client) do_dash_id() error {
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
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: file.FetchWidevine,
   })
}

type WidevineFolder string

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

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
