package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "log"
)

func (c *client) do_hls_id() error {
   return c.Hls.Download(&c.Job, nil)
}

func main() {
   maya.SetProxy("", "*.ts")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_username_password() error {
   var err error
   c.Login, err = oldflix.FetchLogin(c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Hls   *maya.Hls
   Login *oldflix.Login
   //--------------
   Job maya.Job
   //--------------
   username string
   password string
   //--------------
   oldflix_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/oldflix.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   username := maya.StringFlag(&c.username, "u", "username")
   password := maya.StringFlag(&c.password, "p", "password")
   //-------------------------------------------------------------
   oldflix_id := maya.StringFlag(&c.oldflix_id, "o", "Oldflix ID")
   //-------------------------------------------------------------
   hls := maya.IntFlag(&c.Job.Hls, "h", "HLS ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if err != nil {
      return err
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   if oldflix_id.IsSet {
      return with_cache(c.do_oldflix_id)
   }
   if hls.IsSet {
      return with_cache(c.do_hls_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {username, password},
      {oldflix_id},
      {hls},
   })
}

func (c *client) do_oldflix_id() error {
   browse, err := c.Login.FetchBrowse(c.oldflix_id)
   if err != nil {
      return err
   }
   original, err := browse.GetOriginal()
   if err != nil {
      return err
   }
   watch, err := browse.FetchWatch(original.Id, c.Login.Token)
   if err != nil {
      return err
   }
   c.Hls, err = maya.ListHls(watch.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
