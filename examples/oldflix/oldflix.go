package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "fmt"
   "log"
)

type client struct {
   cache     maya.Cache
   Username  maya.Flag
   Password  maya.Flag
   OldflixId maya.Flag
   HlsId     maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/oldflix"); err != nil {
      return err
   }
   c.flag.AddValue(&c.Username, "u", "username")
   c.flag.AddValue(&c.Password, "p", "password")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.OldflixId, "o", "Oldflix ID")
   c.flag.AddValue(&c.HlsId, "h", "HLS ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.Username.Set {
      if c.Password.Set {
         return c.do_username_password()
      }
   }
   if c.OldflixId.Set {
      return c.do_oldflix_id()
   }
   if c.HlsId.Set {
      return c.do_hls_id()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_oldflix_id() error {
   var login oldflix.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   browse, err := login.FetchBrowse(c.OldflixId.Value)
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

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_hls_id() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadHls(c.HlsId.Value, &manifest, nil)
}

func (c *client) do_username_password() error {
   login, err := oldflix.FetchLogin(c.Username.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}
