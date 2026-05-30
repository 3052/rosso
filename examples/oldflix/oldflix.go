package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   flags := maya.FlagSet{
      {Name: "username", Value: &c.username, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "username"},
      {Name: "browse-id", Value: &c.browse},
      {Name: "hls-id", Value: &c.hls},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if c.browse != "" {
      return c.do_browse()
   }
   if c.hls != "" {
      return c.do_hls()
   }
   return flags.Usage(os.Stderr, "oldflix")
}

func (c *client) do_browse() error {
   var login oldflix.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   browse, err := login.FetchBrowse(string(c.browse))
   if err != nil {
      return err
   }
   original, err := browse.GetOriginal()
   if err != nil {
      return err
   }
   watch, err := browse.FetchWatch(original.Id, login.Token)
   if err != nil {
      return err
   }
   manifest, err := maya.ListHls(&watch.Playlist[0].File.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_hls() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadHls(string(c.hls), &manifest, nil)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   username maya.FlagString
   password maya.FlagString
   browse   maya.FlagString
   hls      maya.FlagString

   cache maya.Cache
}

func (c *client) do_username_password() error {
   login, err := oldflix.FetchLogin(string(c.username), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}
