package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   fetch := func(data []byte) ([]byte, error) {
      return amc.License(
         c.Source.KeySystems.ComWidevineAlpha.LicenseURL, c.BcovAuth, data,
      )
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, fetch)
}

func (c *client) do() error {
   err := cache.Setup("rosso/amc.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
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
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
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
      return with_cache(c.do_refresh)
   }
   if series.IsSet {
      return with_cache(c.do_series)
   }
   if season.IsSet {
      return with_cache(c.do_season)
   }
   if episode.IsSet {
      return with_cache(c.do_episode)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {refresh},
      {series},
      {season},
      {episode},
      {dash_id},
   })
}

var cache maya.Cache

func main() {
   maya.SetProxy("", "*.m4f")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
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
   c.Dash, err = c.Source.FetchDash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   AuthData *amc.AuthData
   BcovAuth string
   Dash     *amc.Dash
   Source   *amc.Source
   //------------------------
   Job maya.Job
   //------------------------
   email    string
   password string
   //------------------------
   series int
   //------------------------
   season int
   //------------------------
   episode int
   //------------------------
   dash_id string
}
