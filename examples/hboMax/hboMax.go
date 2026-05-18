package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
)

func (c *client) do_show_season() error {
   var login hboMax.Login
   if err = c.cache.Decode(&login); err != nil {
      return err
   }
   results, err := hboMax.SeasonRequest(
      login.Token, c.show.Value, c.season.Value,
   )
   if err != nil {
      return err
   }
   for i, result := range hboMax.SeasonResults(results) {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(result)
   }
   return nil
}

func (c *client) do_show() error {
   var login hboMax.Login
   if err = c.cache.Decode(&login); err != nil {
      return err
   }
   results, err := hboMax.MovieRequest(login.Token, c.show.Value)
   if err != nil {
      return err
   }
   for i, result := range hboMax.MovieResults(results) {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(result)
   }
   return nil
}

type PlayReadyFolder string

func (c *client) do_dash() error {
   var (
      manifest  maya.Manifest
      playReady PlayReadyFolder
      playback  hboMax.Playback
   )
   err := c.cache.Decode(&manifest, &playReady, &playback)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: playback.PlayReadyRequest,
   })
}

///

func (c *client) do_edit() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   playback, err := hboMax.PlayReadyRequest(login.Token, c.edit.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(playback.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playback)
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
   initiate, err := hboMax.InitiateRequest(st, c.market.Value)
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
   results, err := hboMax.SearchRequest(login.Token, c.search.Value)
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

type client struct {
   cache     maya.Cache
   flag      maya.FlagSet
   dash      maya.Flag
   edit      maya.Flag
   market    maya.Flag
   search    maya.Flag
   show      maya.Flag
   playReady maya.Flag
   season    maya.Flag
   initiate  maya.Flag
   login     maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hboMax"); err != nil {
      return err
   }
   c.flag.AddValue(&c.playReady, "p", "PlayReady")
   c.flag = append(c.flag, nil)
   c.flag.Add(&c.initiate, "i", "initiate")
   c.flag.AddValue(&c.market, "m", fmt.Sprint(hboMax.Markets))
   c.flag = append(c.flag, nil)
   c.flag.Add(&c.login, "l", "login")
   c.flag.AddValue(&c.search, "s", "search")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.show, "SM", "show/movie ID")
   c.flag.AddValue(&c.season, "S", "season")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.edit, "e", "edit ID")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.playReady.Set {
      return c.cache.Encode(PlayReadyFolder(c.playReady.Value))
   }
   if c.initiate.Set {
      if c.market.Set {
         return c.do_initiate()
      }
   }
   if c.login.Set {
      return c.do_login()
   }
   if c.search.Set {
      return c.do_search()
   }
   if c.show.Set {
      if c.season.Set {
         return c.do_show_season()
      }
      return c.do_show()
   }
   if c.edit.Set {
      return c.do_edit()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}
