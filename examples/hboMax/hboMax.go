package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
   "net/http"
)

func main() {
   maya.SetProxy("", "*.mp4")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.FetchPlayReady,
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
   search string
   //-------------------
   show   string
   season int
   //-------------------
   edit string
   //-------------------
   dash_id string
}

func (c *client) do_search() error {
   results, err := c.Login.Search(c.search)
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
   var (
      results []*hboMax.Entity
      err     error
   )
   if c.season >= 1 {
      results, err = c.Login.FetchSeason(c.show, c.season)
      if err != nil {
         return err
      }
      results = hboMax.EpisodeResults(results)
   } else {
      results, err = c.Login.FetchMovie(c.show)
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
   search := maya.StringFlag(&c.search, "s", "search")
   //-------------------------------------------------------------
   show := maya.StringFlag(&c.show, "SM", "show/movie ID")
   season := maya.IntFlag(&c.season, "S", "season")
   //-------------------------------------------------------------
   edit := maya.StringFlag(&c.edit, "e", "edit ID")
   //-------------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
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
   if search.IsSet {
      return with_cache(c.do_search)
   }
   if show.IsSet {
      return with_cache(c.do_show)
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
      {search},
      {show, season},
      {edit},
      {dash_id},
   })
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
   c.Playback, err = c.Login.FetchPlayReady(c.edit)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.FetchDash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
