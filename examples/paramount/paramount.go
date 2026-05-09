package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
)

type content_id string

func (c *client) do_dash() error {
   var (
      cbs_app      paramount.CbsApp
      dash         maya.Dash
      paramount_id content_id
   )
   err := c.cache.Decode(&c.job, &cbs_app, &dash, &paramount_id)
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
   return dash.Download(c.dash, &c.job, session.Fetch)
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
   dash, err := maya.ListDash(&session.StreamingUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(content_id(c.paramount_id), dash)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache        maya.Cache
   app          string
   cookie       *maya.Flag
   dash         string
   job          maya.Job
   paramount_id string
   password     string
   username     string
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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/paramount"); err != nil {
      return err
   }
   c.cookie = maya.BoolFlag("c", "cookie")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   paramount_id := maya.StringFlag(&c.paramount_id, "p", "paramount ID")
   password := maya.StringFlag(&c.password, "P", "password")
   playReady := maya.StringFlag(&c.job.PlayReady, "PR", "PlayReady")
   username := maya.StringFlag(&c.username, "U", "username")
   app := maya.StringFlag(&c.app, "a", paramount.CbsAppIds())
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(c.job)
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
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {app},
      {username, password},
      {paramount_id, c.cookie},
      {dash, c.cookie},
   })
}
