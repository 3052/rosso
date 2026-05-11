package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "log"
)

func (c *client) do_hls() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadHls(c.hls, &manifest, nil)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_username_password() error {
   login, err := oldflix.FetchLogin(c.username, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

type client struct {
   cache      maya.Cache
   flag       maya.FlagSet
   hls        string
   oldflix_id string
   password   string
   username   string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/oldflix"); err != nil {
      return err
   }
   oldflix_id := c.flag.String(&c.oldflix_id, "o", "Oldflix ID")
   password := c.flag.String(&c.password, "p", "password")
   username := c.flag.String(&c.username, "u", "username")
   hls := c.flag.String(&c.hls, "h", "HLS ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   if oldflix_id.IsSet {
      return c.do_oldflix_id()
   }
   if hls.IsSet {
      return c.do_hls()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {username, password},
      {oldflix_id},
      {hls},
   })
}

func (c *client) do_oldflix_id() error {
   var login oldflix.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   browse, err := login.FetchBrowse(c.oldflix_id)
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
