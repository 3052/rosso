package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/amc"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "refresh", Value: &c.refresh},
      {Name: "series", Value: &c.series},
      {Name: "season", Value: &c.season},
      {Name: "episode-or-movie", Value: &c.episode_or_movie},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.series >= 1 {
      return c.do_series()
   }
   if c.season >= 1 {
      return c.do_season()
   }
   if c.episode_or_movie >= 1 {
      return c.do_episode_or_movie()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "amc")
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playback amc.Playback
      source   amc.Source
   )
   err := c.cache.Decode(&manifest, &playback, &source)
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
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}

func (c *client) do_email_password() error {
   auth_data, err := amc.Unauth()
   if err != nil {
      return err
   }
   auth_data, err = amc.Login(
      auth_data.AccessToken, string(c.email), string(c.password),
   )
   if err != nil {
      return err
   }
   return c.cache.Encode(auth_data)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine maya.FlagString

   email            maya.FlagString
   password         maya.FlagString
   refresh          maya.FlagBool
   series           maya.FlagInt
   season           maya.FlagInt
   episode_or_movie maya.FlagInt
   dash             maya.FlagString

   cache maya.Cache
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
   series, err := amc.SeriesDetail(auth_data.AccessToken, int(c.series))
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
   season, err := amc.SeasonEpisodes(auth_data.AccessToken, int(c.season))
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

func (c *client) do_episode_or_movie() error {
   var auth_data amc.AuthData
   err := c.cache.Decode(&auth_data)
   if err != nil {
      return err
   }
   playback, err := amc.GetPlayback(
      auth_data.AccessToken, int(c.episode_or_movie),
   )
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
