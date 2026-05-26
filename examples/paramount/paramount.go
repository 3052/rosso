package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

type client struct {
   PlayReady maya.FlagString

   CbsApp    maya.Flag[string] `usage:"com.cbs.app com.cbs.tve com.cbs.ca"`
   Username  maya.Flag[string] `depends:"Password"`
   Password  maya.Flag[string] `depends:"Username"`
   ContentId maya.Flag[string]
   UseCookie maya.Flag[bool] `depends:"ContentId"`
   DashId    maya.Flag[string]

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/paramount"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "cbs-app", Value: &c.CbsApp},
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.PlayReadyFolder.Set {
      return c.cache.Encode(c.PlayReadyFolder)
   }
   if c.CbsApp.Set {
      return c.do_cbs_app()
   }
   if c.Username.Set {
      if c.Password.Set {
         return c.do_username_password()
      }
   }
   if c.ContentId.Set {
      return c.do_content_id()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "paramount", c)
}

func (c *client) do_dash_id() error {
   var (
      cbs_app    paramount.CbsApp
      content_id ContentId
      manifest   maya.Manifest
      playReady  PlayReadyFolder
   )
   err := c.cache.Decode(&cbs_app, &content_id, &manifest, &playReady)
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
   session, err := cbs_app.FetchPlayReady(content_id.Value, cbs_com)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  playReady.Value,
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

func (c *client) do_content_id() error {
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
   session, err := cbs_app.FetchStreamingUrl(c.ContentId.Value, cbs_com)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&session.StreamingUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(c.ContentId, manifest)
}
