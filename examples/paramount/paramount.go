package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/paramount.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "PR", "PlayReady")
   //--------------------------------------------------------------
   app := maya.StringFlag(&c.App, "a", fmt.Sprint(paramount.GetAppKeys()))
   //--------------------------------------------------------------
   username := maya.StringFlag(&c.username, "U", "username")
   password := maya.StringFlag(&c.password, "P", "password")
   //--------------------------------------------------------------
   paramount_id := maya.StringFlag(&c.ParamountId, "p", "paramount ID")
   c.cookie = maya.BoolFlag("c", "cookie")
   //--------------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
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

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
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

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

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

type client struct {
   // cache
   App         string
   CbsCom      string
   Dash        *maya.Dash
   Job         maya.Job
   ParamountId string
   // flags
   cookie   *maya.Flag
   password string
   username string
   // state
   cache_err error
}
