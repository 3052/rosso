package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      options  maya.Options
      playback amc.Playback
      source   amc.Source
   )
   err := c.cache.Decode(&manifest, &options.Device, &playback, &source)
   if err != nil {
      return err
   }
   options.Drm = maya.DrmWidevine
   options.License = func(body []byte) ([]byte, error) {
      return amc.License(
         source.KeySystems.ComWidevineAlpha.LicenseURL,
         playback.BcovAuth,
         body,
      )
   }
   return maya.DownloadDash(c.dash, &manifest, &options)
}

func (c *client) do_season() error {
   var auth_data amc.AuthData
   err := c.cache.Decode(&auth_data)
   if err != nil {
      return err
   }
   season, err := amc.SeasonEpisodes(auth_data.AccessToken, c.season)
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

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache    maya.Cache
   dash     string
   email    string
   episode  int
   password string
   season   int
   series   int
   widevine string
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

func (c *client) do_refresh() error {
   var auth_data amc.AuthData
   err := c.cache.Decode(&auth_data)
   if err != nil {
      return err
   }
   err = auth_data.Refresh()
   if err != nil {
      return err
   }
   return c.cache.Encode(auth_data)
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
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   widevine := maya.StringFlag(&c.widevine, "w", "Widevine")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(maya.Device(c.widevine))
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

func (c *client) do_episode() error {
   var auth_data amc.AuthData
   err := c.cache.Decode(&auth_data)
   if err != nil {
      return err
   }
   playback, err := amc.GetPlayback(auth_data.AccessToken, c.episode)
   if err != nil {
      return err
   }
   source, err := playback.GetDash()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&source.Src.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playback, source)
}
