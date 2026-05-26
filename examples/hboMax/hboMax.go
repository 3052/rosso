package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
   "os"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playback hboMax.Playback
   )
   err := c.cache.Decode(&manifest, &playback)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
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

type client struct {
   PlayReady maya.FlagString

   dash     maya.FlagString
   edit     maya.FlagString
   initiate maya.FlagString
   login    maya.FlagBool
   movie    maya.FlagString
   search   maya.FlagString
   season   maya.FlagInt
   show     maya.FlagString

   cache maya.Cache
}

func (c *client) do_initiate() error {
   st, err := hboMax.StRequest()
   if err != nil {
      return err
   }
   initiate, err := hboMax.InitiateRequest(st, string(c.initiate))
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
   results, err := hboMax.SearchRequest(login.Token, string(c.search))
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

func (c *client) do_movie() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   results, err := hboMax.MovieRequest(login.Token, string(c.movie))
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

func (c *client) do_show_season() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   results, err := hboMax.SeasonRequest(
      login.Token, string(c.show), int(c.season),
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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hboMax"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "initiate", Value: &c.initiate, Usage: hboMax.Markets},
      {Name: "login", Value: &c.login},
      {Name: "search", Value: &c.search},
      {Name: "movie-id", Value: &c.movie},
      {Name: "show-id", Value: &c.show, Needs: "season"},
      {Name: "season", Value: &c.season, Needs: "show-id"},
      {Name: "edit-id", Value: &c.edit},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.PlayReady) {
      return c.cache.Encode(c)
   }
   if c.initiate != "" {
      return c.do_initiate()
   }
   if c.login {
      return c.do_login()
   }
   if c.search != "" {
      return c.do_search()
   }
   if c.movie != "" {
      return c.do_movie()
   }
   if c.show != "" {
      if c.season >= 1 {
         return c.do_show_season()
      }
   }
   if c.edit != "" {
      return c.do_edit()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "hboMax")
}

func (c *client) do_edit() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   playback, err := hboMax.PlayReadyRequest(login.Token, string(c.edit))
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(playback.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playback)
}
