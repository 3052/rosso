package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache        maya.Cache
   cookie       *maya.Flag
   err          error
   paramount_id string
   password     string
   username     string
   job          maya.Job
   cbs_app      paramount.CbsApp
}

///

func (c *client) do() error {
   if err := c.cache.Setup("rosso/paramount"); err != nil {
      return err
   }
   password := maya.StringFlag(&c.password, "P", "password")
   username := maya.StringFlag(&c.username, "U", "username")
   c.cookie = maya.BoolFlag("c", "cookie")
   paramount_id := maya.StringFlag(&c.paramount_id, "p", "paramount ID")

   c.err = c.cache.Decode(c)

   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")
   playReady := maya.StringFlag(&c.job.PlayReady, "PR", "PlayReady")
   app := maya.StringFlag(&c.cbs_app.Id, "a", paramount.CbsAppIds())

   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
      return cache.Write(c)
   }
   if app.IsSet {
      return cache.Write(c)
   }
   if username.IsSet {
      if password.IsSet {
         return c.run(c.do_username_password)
      }
   }
   if paramount_id.IsSet {
      return c.run(c.do_paramount)
   }
   if dash.IsSet {
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {app},
      {username, password},
      {paramount_id, c.cookie},
      {dash, c.cookie},
   })
}

func (c *client) do_username_password() error {
   app, err := paramount.GetApp(c.App)
   if err != nil {
      return err
   }
   c.CbsCom, err = app.FetchCbsCom(c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_paramount() error {
   app, err := paramount.GetApp(c.App)
   if err != nil {
      return err
   }
   var cbs_com string
   if c.cookie.IsSet {
      cbs_com = c.CbsCom
   }
   session, err := app.FetchStreamingUrl(c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(session.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash() error {
   app, err := paramount.GetApp(c.App)
   if err != nil {
      return err
   }
   var cbs_com string
   if c.cookie.IsSet {
      cbs_com = c.CbsCom
   }
   session, err := app.FetchPlayReady(c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   return c.Dash.Download(&c.Job, session.Fetch)
}
