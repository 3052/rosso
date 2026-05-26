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

type client struct {
   App       maya.FlagString
   ContentId maya.FlagString
   PlayReady maya.FlagString

   cookie   maya.FlagBool
   dash     maya.FlagString
   password maya.FlagString
   username maya.FlagString

   cache maya.Cache
}

func (c *client) do_username_password() error {
   app, err := paramount.GetApp(string(c.App))
   if err != nil {
      return err
   }
   cbs_com, err := app.FetchCbsCom(string(c.username), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(cbs_com)
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
      {Name: "cbs-app", Value: &c.CbsApp, Usage: paramount.AppIds()},
      {Name: "username", Value: &c.username, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "username"},
      {Name: "content-id", Value: &c.ContentId},
      {Name: "use-cookie", Value: &c.cookie, Needs: "content-id"},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.App) {
      return c.cache.Encode(c)
   }
   if flags.IsSet(&c.PlayReady) {
      return c.cache.Encode(c)
   }
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if flags.IsSet(&c.ContentId) {
      return c.do_content_id()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "paramount")
}

///

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

func (c *client) do_dash() error {
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
