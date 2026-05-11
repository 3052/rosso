package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hboMax"); err != nil {
      return err
   }
   edit := c.flag.String(&c.edit, "e", "edit ID")
   initiate := c.flag.Bool("i", "initiate")
   login := c.flag.Bool("l", "login")
   market := c.flag.String(&c.market, "m", fmt.Sprint(hboMax.Markets))
   search := c.flag.String(&c.search, "s", "search")
   season := c.flag.Int(&c.season, "S", "season")
   show := c.flag.String(&c.show, "SM", "show/movie ID")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   playReady := c.flag.String(&c.playReady, "p", "PlayReady")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(device(c.playReady))
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
   return maya.PrintFlags([]maya.FlagSet{
      {playReady},
      {initiate, market},
      {login},
      {search},
      {show, season},
      {edit},
      {dash},
   })
}

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
   manifest, err := maya.ListDash(playback.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playback)
}

func (c *client) do_dash() error {
   var (
      manifest  maya.Manifest
      playReady device
      playback  hboMax.Playback
   )
   err := c.cache.Decode(&manifest, &playReady, &playback)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: playback.PlayReadyRequest,
   })
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
   cache     maya.Cache
   dash      string
   edit      string
   flag      maya.FlagSet
   market    string
   search    string
   season    int
   show      string
   playReady string
}

type device string
