package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "log"
)

func (c *client) do_hls() error {
   var hls maya.Hls
   err := c.cache.Decode(&c.job, &hls)
   if err != nil {
      return err
   }
   return hls.Download(c.hls, &c.job, nil)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache      maya.Cache
   hls        int
   job        maya.Job
   oldflix_id string
   password   string
   username   string
}

func (c *client) do_username_password() error {
   login, err := oldflix.FetchLogin(c.username, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/oldflix"); err != nil {
      return err
   }
   oldflix_id := maya.StringFlag(&c.oldflix_id, "o", "Oldflix ID")
   password := maya.StringFlag(&c.password, "p", "password")
   username := maya.StringFlag(&c.username, "u", "username")
   hls := maya.IntFlag(&c.hls, "h", "HLS ID")
   if err := maya.ParseFlags(); err != nil {
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
   return maya.PrintFlags([][]*maya.Flag{
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
   hls, err := maya.ListHls(&watch.Playlist[0].File.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(hls)
}
