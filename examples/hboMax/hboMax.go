package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

func (c *client) do_initiate() error {
   st, err := hboMax.StRequest()
   if err != nil {
      return err
   }
   initiate, err := hboMax.InitiateRequest(st, c.Initiate.Value)
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
   results, err := hboMax.SearchRequest(login.Token, c.Search.Value)
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

func (c *client) do_show_id_season() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   results, err := hboMax.SeasonRequest(
      login.Token, c.ShowId.Value, c.Season.Value,
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

func (c *client) do_movie_id() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   results, err := hboMax.MovieRequest(login.Token, c.MovieId.Value)
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

type PlayReadyFolder maya.Flag[string]

type client struct {
   cache           maya.Cache
   PlayReadyFolder PlayReadyFolder
   Initiate        maya.Flag[string] `usage:"amer apac emea latam"`
   Login           maya.Flag[bool]
   Search          maya.Flag[string]
   MovieId         maya.Flag[string]
   ShowId          maya.Flag[string] `depends:"Season"`
   Season          maya.Flag[int]    `depends:"ShowId"`
   EditId          maya.Flag[string]
   DashId          maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hboMax"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.PlayReadyFolder.Set {
      return c.cache.Encode(c.PlayReadyFolder)
   }
   if c.Initiate.Set {
      return c.do_initiate()
   }
   if c.Login.Set {
      return c.do_login()
   }
   if c.Search.Set {
      return c.do_search()
   }
   if c.MovieId.Set {
      return c.do_movie_id()
   }
   if c.ShowId.Set {
      if c.Season.Set {
         return c.do_show_id_season()
      }
   }
   if c.EditId.Set {
      return c.do_edit_id()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "hboMax", c)
}

func (c *client) do_dash_id() error {
   var (
      manifest  maya.Manifest
      playReady PlayReadyFolder
      playback  hboMax.Playback
   )
   err := c.cache.Decode(&manifest, &playReady, &playback)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  playReady.Value,
      Drm:     maya.DrmPlayReady,
      License: playback.PlayReadyRequest,
   })
}

func (c *client) do_edit_id() error {
   var login hboMax.Login
   err := c.cache.Decode(&login)
   if err != nil {
      return err
   }
   playback, err := hboMax.PlayReadyRequest(login.Token, c.EditId.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(playback.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playback)
}
