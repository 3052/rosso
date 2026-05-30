package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
   "os"
)

func (*client) CachePath() string {
   return "rosso/examples/paramount/client"
}

type client struct {
   App       maya.FlagString
   ContentId maya.FlagString
   PlayReady maya.FlagString
   cookie    maya.FlagBool
   dash      maya.FlagString
   password  maya.FlagString
   username  maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "app", Value: &c.App, Usage: paramount.AppIds()},
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

func (c *client) do_content_id() error {
   app, err := paramount.GetApp(string(c.App))
   if err != nil {
      return err
   }
   var cbs_com *paramount.Cookie
   if c.cookie {
      cbs_com = &paramount.Cookie{}
      err = c.cache.Decode(cbs_com)
      if err != nil {
         return err
      }
   }
   session, err := app.FetchStreamingUrl(string(c.ContentId), cbs_com)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&session.StreamingUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(c, manifest)
}

func (c *client) do_dash() error {
   // 1. manifest
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   // 2. app
   app, err := paramount.GetApp(string(c.App))
   if err != nil {
      return err
   }
   // 3. cookie
   var cbs_com *paramount.Cookie
   if c.cookie {
      cbs_com = &paramount.Cookie{}
      err = c.cache.Decode(cbs_com)
      if err != nil {
         return err
      }
   }
   // 4. session
   session, err := app.FetchPlayReady(string(c.ContentId), cbs_com)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: session.Fetch,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
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
