package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
   "net/http"
)

func (c *client) do() error {
   err := cache.Setup("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "p", "PlayReady")
   //-------------------------------------------------------------
   initiate := maya.BoolFlag("i", "initiate")
   market := maya.StringFlag(&c.market, "m", fmt.Sprint(hboMax.Markets))
   //-------------------------------------------------------------
   login := maya.BoolFlag("l", "login")
   //-------------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   season := maya.IntFlag(&c.season, "s", "season")
   //-------------------------------------------------------------
   edit := maya.StringFlag(&c.edit, "e", "edit ID")
   //-------------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   maya.SetProxy("", "*.mp4")
   if playReady.IsSet {
      return cache.Write(c)
   }
   if initiate.IsSet {
      if market.IsSet {
         return c.do_initiate()
      }
   }
   if login.IsSet {
      return with_cache(c.do_login)
   }
   if address.IsSet {
      return with_cache(c.do_address)
   }
   if edit.IsSet {
      return with_cache(c.do_edit)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {initiate, market},
      {login},
      {address, season},
      {edit},
      {dash_id},
   })
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.PlayReady,
   )
}

func (c *client) do_login() error {
   var err error
   c.Login, err = hboMax.FetchLogin(c.St)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   show_id, err := hboMax.ParseShowId(c.address)
   if err != nil {
      return err
   }
   var page *hboMax.Page
   if c.season >= 1 {
      page, err = c.Login.Season(show_id, c.season)
   } else {
      page, err = c.Login.Movie(show_id)
   }
   if err != nil {
      return err
   }
   page.FilterAndSort()
   for i, video := range page.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

func (c *client) do_initiate() error {
   var err error
   c.St, err = hboMax.FetchSt()
   if err != nil {
      return err
   }
   initiate, err := hboMax.FetchInitiate(c.St, c.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return cache.Write(c)
}

func (c *client) do_edit() error {
   var err error
   c.Playback, err = c.Login.PlayReady(c.edit)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *http.Cookie
   //-------------------
   Job maya.Job
   //-------------------
   market string
   //-------------------
   address string
   season  int
   //-------------------
   edit string
   //-------------------
   dash_id string
}
