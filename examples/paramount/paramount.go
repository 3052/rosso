package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
)

type client struct {
   cache        maya.Cache
   cookie       maya.Flag
   flag         maya.FlagSet
   app          maya.Flag
   dash         maya.Flag
   paramount_id maya.Flag
   password     maya.Flag
   username     maya.Flag
   playReady    maya.Flag
}

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
   if playReady.Set {
      return c.cache.Encode(playReady_device(c.playReady))
   }
   if app.Set {
      return c.do_app()
   }
   if username.Set {
      if password.Set {
         return c.do_username_password()
      }
   }
   if paramount_id.Set {
      return c.do_paramount_id()
   }
   if dash.Set {
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

type playReady_device string

func (c *client) do_dash() error {
   var (
      cbs_app      paramount.CbsApp
      device       playReady_device
      manifest     maya.Manifest
      paramount_id content_id
   )
   err := c.cache.Decode(&cbs_app, &device, &manifest, &paramount_id)
   if err != nil {
      return err
   }
   var cbs_com *paramount.Cookie
   if c.cookie.Set {
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
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmPlayReady,
      License: session.Fetch,
   })
}

func (c *client) do_app() error {
   cbs_app, err := paramount.GetCbsApp(c.app.Value)
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
   cbs_com, err := cbs_app.FetchCbsCom(c.username.Value, c.password.Value)
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

func (c *client) do_paramount_id() error {
   var cbs_app paramount.CbsApp
   err := c.cache.Decode(&cbs_app)
   if err != nil {
      return err
   }
   var cbs_com *paramount.Cookie
   if c.cookie.Set {
      cbs_com = &paramount.Cookie{}
      err = c.cache.Decode(cbs_com)
      if err != nil {
         return err
      }
   }
   session, err := cbs_app.FetchStreamingUrl(c.paramount_id.Value, cbs_com)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&session.StreamingUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(content_id(c.paramount_id.Value), manifest)
}
