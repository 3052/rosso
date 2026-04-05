package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/oldflix"
   "log"
)

func (c *client) do_hls_id() error {
   return c.Job.DownloadHls(c.Hls.Body, c.Hls.Url, c.hls_id, nil)
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
   hls_id := maya.IntFlag(&c.hls_id, "h", "HLS ID")
   err = maya.ParseFlags()
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
   if hls_id.IsSet {
      return with_cache(c.do_hls_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {username, password},
      {oldflix_id},
      {hls_id},
   })
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "")
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
   Hls   *oldflix.Hls
   Login *oldflix.Login
   //--------------
   Job maya.Job
   //--------------
   username string
   password string
   //--------------
   oldflix_id string
   //--------------
   hls_id int
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
   c.Hls, err = watch.FetchHls()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListHls(c.Hls.Body, c.Hls.Url)
}
