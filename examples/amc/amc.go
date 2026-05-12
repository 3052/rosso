package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

type client struct {
   cache    maya.Cache
   dash *maya.Flag
   email *maya.Flag
   episode *maya.Flag
   flag     maya.FlagSet
   password *maya.Flag
   season *maya.Flag
   series *maya.Flag
   widevine *maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/amc"); err != nil {
      return err
   }
   c.dash = c.flag.String("d", "DASH ID")
   c.email = c.flag.String("E", "email")
   c.episode = c.flag.String("e", "episode or movie ID")
   c.password = c.flag.String("P", "password")
   c.season = c.flag.String("S", "season ID")
   c.series = c.flag.String("s", "series ID")
   c.widevine = c.flag.String("w", "Widevine")
   refresh := c.flag.Bool("r", "refresh")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine(c.widevine))
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if refresh.Set {
      return c.do_refresh()
   }
   if c.series.Set {
      return c.do_series()
   }
   if c.season.Set {
      return c.do_season()
   }
   if c.episode.Set {
      return c.do_episode()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {email, password},
      {refresh},
      {series},
      {season},
      {episode},
      {dash},
   })
}

type widevine string

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

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playback amc.Playback
      source   amc.Source
      device widevine
   )
   err := c.cache.Decode(&manifest, &playback, &source, &device)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return amc.License(
         source.KeySystems.ComWidevineAlpha.LicenseURL,
         playback.BcovAuth,
         body,
      )
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: license,
   })
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
