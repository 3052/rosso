package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   auth_data, err := amc.Unauth()
   if err != nil {
      return err
   }
   auth_data, err = amc.Login(auth_data.AccessToken, c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(auth_data)
}

type client struct {
   cache    maya.Cache
   job      maya.Job
   email    string
   episode  int
   password string
   season   int
   series   int
   //BcovAuth string
   //Dash     *maya.Dash
   //Source   *amc.Source
}

func (c *client) do_refresh() error {
   var auth_data amc.AuthData
   return c.cache.Update(&auth_data, func() error {
      return auth_data.Refresh()
   })
}

func (c *client) do_series() error {
   var auth_data amc.AuthData
   err := c.cache.Decode(&auth_data)
   if err != nil {
      return err
   }
   series, err := amc.SeriesDetail(auth_data.AccessToken, c.series)
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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/amc"); err != nil {
      return err
   }
   email := maya.StringFlag(&c.email, "E", "email")
   password := maya.StringFlag(&c.password, "P", "password")
   refresh := maya.BoolFlag("r", "refresh")
   series := maya.IntFlag(&c.series, "s", "series ID")
   season := maya.IntFlag(&c.season, "S", "season ID")
   episode := maya.IntFlag(&c.episode, "e", "episode or movie ID")
   c.cache.Decode(&c.job)
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if refresh.IsSet {
      return c.do_refresh()
   }
   if series.IsSet {
      return c.do_series()
   }
   if season.IsSet {
      return c.do_season()
   }
   if episode.IsSet {
      return c.do_episode()
   }
   if dash.IsSet {
      return c.do_dash()
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

///

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

func (c *client) do_dash() error {
   fetch := func(data []byte) ([]byte, error) {
      return amc.License(
         c.Source.KeySystems.ComWidevineAlpha.LicenseURL, c.BcovAuth, data,
      )
   }
   return c.Dash.Download(&c.job, fetch)
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
   return c.cache.Write(c)
}
