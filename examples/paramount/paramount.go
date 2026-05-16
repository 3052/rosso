package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/paramount"); err != nil {
      return err
   }
   c.flag.AddValue(&c.app, "a", paramount.CbsAppIds())
   c.flag.AddValue(&c.playReady, "PR", "PlayReady")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.username, "U", "username")
   c.flag.AddValue(&c.password, "P", "password")
   c.flag = append(c.flag, nil)
   c.flag.Add(&c.cookie, "c", "cookie")
   c.flag.AddValue(&c.paramount_id, "p", "paramount ID")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.playReady.Set {
      return c.cache.Encode(playReady_device(c.playReady.Value))
   }
   if c.app.Set {
      return c.do_app()
   }
   if c.username.Set {
      if c.password.Set {
         return c.do_username_password()
      }
   }
   if c.paramount_id.Set {
      return c.do_paramount_id()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
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
