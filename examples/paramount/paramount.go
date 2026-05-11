package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
)

func (c *client) do_dash() error {
   var (
      cbs_app      paramount.CbsApp
      manifest     maya.Manifest
      paramount_id content_id
      playReady    playReady_folder
   )
   err := c.cache.Decode(&cbs_app, &manifest, &paramount_id, &playReady)
   if err != nil {
      return err
   }
   var cbs_com *paramount.Cookie
   if c.cookie.IsSet {
      cbs_com = &paramount.Cookie{}
      err = c.cache.Decode(cbs_com)
      if err != nil {
         return err
      }
   }
   session, err := cbs_app.FetchPlayReady(string(paramount_id), cbs_com)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: session.Fetch,
   })
}

func (c *client) do_app() error {
   cbs_app, err := paramount.GetCbsApp(c.app)
   if err != nil {
      return err
   }
   return c.cache.Encode(cbs_app)
}

func (c *client) do_username_password() error {
   var cbs_app paramount.CbsApp
   err := c.cache.Decode(&cbs_app)
   if err != nil {
      return err
   }
   cbs_com, err := cbs_app.FetchCbsCom(c.username, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(cbs_com)
}

type content_id string

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   app          string
   cache        maya.Cache
   cookie       *maya.Flag
   dash         string
   flag         maya.FlagSet
   paramount_id string
   password     string
   username     string
   playReady    string
}

type playReady_folder string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/paramount"); err != nil {
      return err
   }
   c.cookie = c.flag.Bool("c", "cookie")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   paramount_id := c.flag.String(&c.paramount_id, "p", "paramount ID")
   password := c.flag.String(&c.password, "P", "password")
   username := c.flag.String(&c.username, "U", "username")
   app := c.flag.String(&c.app, "a", paramount.CbsAppIds())
   playReady := c.flag.String(&c.playReady, "PR", "PlayReady")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(playReady_folder(c.playReady))
   }
   if app.IsSet {
      return c.do_app()
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   if paramount_id.IsSet {
      return c.do_paramount_id()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {playReady},
      {app},
      {username, password},
      {paramount_id, c.cookie},
      {dash, c.cookie},
   })
}

func (c *client) do_paramount_id() error {
   var cbs_app paramount.CbsApp
   err := c.cache.Decode(&cbs_app)
   if err != nil {
      return err
   }
   var cbs_com *paramount.Cookie
   if c.cookie.IsSet {
      cbs_com = &paramount.Cookie{}
      err = c.cache.Decode(cbs_com)
      if err != nil {
         return err
      }
   }
   session, err := cbs_app.FetchStreamingUrl(c.paramount_id, cbs_com)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&session.StreamingUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(content_id(c.paramount_id), manifest)
}
