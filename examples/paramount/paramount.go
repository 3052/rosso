package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
   "os"
)

type client struct {
   cache           maya.Cache
   CbsApp          maya.Flag[string] `usage:"com.cbs.app com.cbs.tve com.cbs.ca"`
   PlayReadyFolder maya.Flag[string]
   Username        maya.Flag[string] `depends:"Password"`
   Password        maya.Flag[string] `depends:"Username"`
   ParamountId     maya.Flag[string]
   UseCookie       maya.Flag[bool] `depends:"ParamountId"`
   DashId          maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/paramount"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.PlayReadyFolder.Set {
      return c.cache.Encode(PlayReadyFolder(c.PlayReadyFolder.Value))
   }
   if c.CbsApp.Set {
      return c.do_cbs_app()
   }
   if c.Username.Set {
      if c.Password.Set {
         return c.do_username_password()
      }
   }
   if c.ParamountId.Set {
      return c.do_paramount_id()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "paramount", c)
}

type PlayReadyFolder string

type ParamountIdString string

func (c *client) do_dash_id() error {
   var (
      cbs_app      paramount.CbsApp
      manifest     maya.Manifest
      paramount_id ParamountIdString
      playReady    PlayReadyFolder
   )
   err := c.cache.Decode(&cbs_app, &manifest, &paramount_id, &playReady)
   if err != nil {
      return err
   }
   var cbs_com *paramount.Cookie
   if c.UseCookie.Set {
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
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: session.Fetch,
   })
}

func (c *client) do_cbs_app() error {
   cbs_app, err := paramount.GetCbsApp(c.CbsApp.Value)
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
   cbs_com, err := cbs_app.FetchCbsCom(c.Username.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(cbs_com)
}

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
   if c.UseCookie.Set {
      cbs_com = &paramount.Cookie{}
      err = c.cache.Decode(cbs_com)
      if err != nil {
         return err
      }
   }
   session, err := cbs_app.FetchStreamingUrl(c.ParamountId.Value, cbs_com)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&session.StreamingUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(ParamountIdString(c.ParamountId.Value), manifest)
}
