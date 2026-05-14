package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "fmt"
   "log"
)

func (c *client) do_oldflix_id() error {
   var login oldflix.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   browse, err := login.FetchBrowse(c.oldflix_id.Value)
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

type client struct {
   cache      maya.Cache
   flag       maya.FlagSet
   hls        maya.Flag
   oldflix_id maya.Flag
   password   maya.Flag
   username   maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/oldflix"); err != nil {
      return err
   }
   c.flag.AddValue(&c.username, "u", "username")
   c.flag.AddValue(&c.password, "p", "password")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.oldflix_id, "o", "Oldflix ID")
   c.flag.AddValue(&c.hls, "h", "HLS ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.username.Set {
      if c.password.Set {
         return c.do_username_password()
      }
   }
   if c.oldflix_id.Set {
      return c.do_oldflix_id()
   }
   if c.hls.Set {
      return c.do_hls()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_hls() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadHls(c.hls.Value, &manifest, nil)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_username_password() error {
   login, err := oldflix.FetchLogin(c.username.Value, c.password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}
