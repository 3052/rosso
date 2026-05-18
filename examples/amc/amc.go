package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      playback amc.Playback
      source   amc.Source
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &playback, &source, &widevine)
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
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}

type WidevineFolder string

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
   auth_data, err = amc.Login(
      auth_data.AccessToken, c.email.Value, c.password.Value,
   )
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
   if err = c.cache.Decode(&auth_data); err != nil {
      return err
   }
   series, err := amc.SeriesDetail(auth_data.AccessToken, c.SeriesId.Value)
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

///

func (c *client) do_season() error {
   id, err := c.season.ParseInt()
   if err != nil {
      return err
   }
   var auth_data amc.AuthData
   if err = c.cache.Decode(&auth_data); err != nil {
      return err
   }
   season, err := amc.SeasonEpisodes(auth_data.AccessToken, id)
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
   video_id, err := c.episode.ParseInt()
   if err != nil {
      return err
   }
   var auth_data amc.AuthData
   if err = c.cache.Decode(&auth_data); err != nil {
      return err
   }
   playback, err := amc.GetPlayback(auth_data.AccessToken, video_id)
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

type client struct {
   cache maya.Cache

   DashId   maya.Flag
   email    maya.Flag
   episode  maya.Flag
   password maya.Flag
   season   maya.Flag
   SeriesId maya.Flag
   widevine maya.Flag
   refresh  maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/amc"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.email, "E", "email")
   c.flag.AddValue(&c.password, "P", "password")
   c.flag = append(c.flag, nil)
   c.flag.Add(&c.refresh, "r", "refresh")
   c.flag.AddValue(&c.SeriesId, "s", "series ID")
   c.flag.AddValue(&c.season, "S", "season ID")
   c.flag.AddValue(&c.episode, "e", "episode or movie ID")
   c.flag.AddValue(&c.DashId, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(WidevineFolder(c.widevine.Value))
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if c.refresh.Set {
      return c.do_refresh()
   }
   if c.SeriesId.Set {
      return c.do_series()
   }
   if c.season.Set {
      return c.do_season()
   }
   if c.episode.Set {
      return c.do_episode()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   fmt.Println(c.flag)
   return nil
}
