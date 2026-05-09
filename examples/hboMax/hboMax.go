package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
)

func (c *client) do_edit() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   playback, err := hboMax.PlayReadyRequest(login.Token, c.edit)
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(playback.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, playback)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_initiate() error {
   st, err := hboMax.StRequest()
   if err != nil {
      return err
   }
   initiate, err := hboMax.InitiateRequest(st, c.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return c.cache.Encode(st)
}

func (c *client) do_login() error {
   st := &hboMax.Cookie{}
   err := c.cache.Decode(st)
   if err != nil {
      return err
   }
   login, err := hboMax.LoginRequest(st)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

func (c *client) do_search() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   results, err := hboMax.SearchRequest(login.Token, c.search)
   if err != nil {
      return err
   }
   results, err = hboMax.SearchResults(results)
   if err != nil {
      return err
   }
   for i, result := range results {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(result)
   }
   return nil
}

func (c *client) do_show() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   var results []*hboMax.Entity
   if c.season >= 1 {
      results, err = hboMax.SeasonRequest(login.Token, c.show, c.season)
      if err != nil {
         return err
      }
      results = hboMax.SeasonResults(results)
   } else {
      results, err = hboMax.MovieRequest(login.Token, c.show)
      if err != nil {
         return err
      }
      results = hboMax.MovieResults(results)
   }
   for i, result := range results {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(result)
   }
   return nil
}

type client struct {
   cache  maya.Cache
   dash   string
   edit   string
   job    maya.Job
   market string
   search string
   season int
   show   string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hboMax"); err != nil {
      return err
   }
   edit := maya.StringFlag(&c.edit, "e", "edit ID")
   initiate := maya.BoolFlag("i", "initiate")
   login := maya.BoolFlag("l", "login")
   market := maya.StringFlag(&c.market, "m", fmt.Sprint(hboMax.Markets))
   search := maya.StringFlag(&c.search, "s", "search")
   season := maya.IntFlag(&c.season, "S", "season")
   show := maya.StringFlag(&c.show, "SM", "show/movie ID")
   playReady := maya.StringFlag(&c.job.PlayReady, "p", "PlayReady")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(c.job)
   }
   if initiate.IsSet {
      if market.IsSet {
         return c.do_initiate()
      }
   }
   if login.IsSet {
      return c.do_login()
   }
   if search.IsSet {
      return c.do_search()
   }
   if show.IsSet {
      return c.do_show()
   }
   if edit.IsSet {
      return c.do_edit()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {initiate, market},
      {login},
      {search},
      {show, season},
      {edit},
      {dash},
   })
}

func (c *client) do_dash() error {
   var (
      dash     maya.Dash
      playback hboMax.Playback
   )
   err := c.cache.Decode(&c.job, &dash, &playback)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, playback.PlayReadyRequest)
}
