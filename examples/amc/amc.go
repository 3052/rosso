package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/amc.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "E", "email")
   password := maya.StringFlag(&c.password, "P", "password")
   //----------------------------------------------------------
   refresh := maya.BoolFlag("r", "refresh")
   //----------------------------------------------------------
   series := maya.IntFlag(&c.series, "s", "series ID")
   //----------------------------------------------------------
   season := maya.IntFlag(&c.season, "S", "season ID")
   //----------------------------------------------------------
   episode := maya.IntFlag(&c.episode, "e", "episode or movie ID")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if refresh.IsSet {
      return c.run(c.do_refresh)
   }
   if series.IsSet {
      return c.run(c.do_series)
   }
   if season.IsSet {
      return c.run(c.do_season)
   }
   if episode.IsSet {
      return c.run(c.do_episode)
   }
   if dash.IsSet {
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {refresh},
      {series},
      {season},
      {episode},
      {dash},
   })
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
}

type client struct {
   // cache
   AuthData *amc.AuthData
   BcovAuth string
   Dash     *maya.Dash
   Job      maya.Job
   Source   *amc.Source
   // flags
   email    string
   episode  int
   password string
   season   int
   series   int
   // state
   cache_err error
}

func (c *client) do_dash() error {
   fetch := func(data []byte) ([]byte, error) {
      return amc.License(
         c.Source.KeySystems.ComWidevineAlpha.LicenseURL, c.BcovAuth, data,
      )
   }
   return c.Dash.Download(&c.Job, fetch)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_episode() error {
   playback, err := amc.Playback(c.AuthData.AccessToken, c.episode)
   if err != nil {
      return err
   }
   c.Source, err = playback.Data.DashSource()
   if err != nil {
      return err
   }
   c.BcovAuth = playback.BcovAuth
   c.Dash, err = maya.ListDash(c.Source.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_email_password() error {
   var err error
   c.AuthData, err = amc.Unauth()
   if err != nil {
      return err
   }
   c.AuthData, err = amc.Login(c.AuthData.AccessToken, c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_refresh() error {
   var err error
   c.AuthData, err = amc.Refresh(c.AuthData.RefreshToken)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_series() error {
   series, err := amc.SeriesDetail(c.AuthData.AccessToken, c.series)
   if err != nil {
      return err
   }
   for i, season := range series.SeasonsMetadata() {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(season)
   }
   return nil
}

func (c *client) do_season() error {
   season, err := amc.SeasonEpisodes(c.AuthData.AccessToken, c.season)
   if err != nil {
      return err
   }
   for i, episode := range season.EpisodesMetadata() {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(episode)
   }
   return nil
}
