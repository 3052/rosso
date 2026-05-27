package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
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

func (c *client) do_audio_language() error {
   var (
      address     rakuten.Address
      start       rakuten.Start
      stream_info *rakuten.StreamInfo
   )
   err := c.cache.Decode(&address, &start)
   if err != nil {
      return err
   }
   switch {
   case address.IsMovie():
      stream_info, err = rakuten.FetchMovieStreaming(
         address.ContentId, start.Profile.Classification, c.AudioLanguage.Value,
      )
   case address.IsTvShow():
      stream_info, err = rakuten.FetchEpisodeStreaming(
         c.EpisodeId.Value, start.Profile.Classification, c.AudioLanguage.Value,
      )
   }
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&stream_info.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, stream_info)
}

type PlayReadyFolder maya.Flag[string]

type client struct {
   cache           maya.Cache
   PlayReadyFolder PlayReadyFolder
   Address         maya.Flag[string]
   SeasonId        maya.Flag[string]
   AudioLanguage   maya.Flag[string]
   EpisodeId       maya.Flag[string] `depends:"AudioLanguage"`
   DashId          maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rakuten"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.PlayReadyFolder.Set:
      return c.cache.Encode(c.PlayReadyFolder)
   case c.Address.Set:
      return c.do_address()
   case c.SeasonId.Set:
      return c.do_season_id()
   case c.AudioLanguage.Set:
      return c.do_audio_language()
   case c.DashId.Set:
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "rakuten", c)
}

func (c *client) do_dash_id() error {
   var (
      manifest    maya.Manifest
      playReady   PlayReadyFolder
      stream_info rakuten.StreamInfo
   )
   err := c.cache.Decode(&manifest, &playReady, &stream_info)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  playReady.Value,
      Drm:     maya.DrmPlayReady,
      License: stream_info.FetchLicense,
   })
}

func (c *client) do_address() error {
   address, err := rakuten.ParseAddress(c.Address.Value)
   if err != nil {
      return err
   }
   start, err := rakuten.FetchStart(address.MarketCode)
   if err != nil {
      return err
   }
   switch {
   case address.IsMovie():
      movie, err := rakuten.FetchMovie(
         address.ContentId, start.Profile.Classification, start.Market,
      )
      if err != nil {
         return err
      }
      fmt.Println(movie)
   case address.IsTvShow():
      show, err := rakuten.FetchTvShow(
         address.ContentId, start.Profile.Classification, start.Market,
      )
      if err != nil {
         return err
      }
      fmt.Println(show)
   }
   return c.cache.Encode(address, start)
}

func (c *client) do_season_id() error {
   var start rakuten.Start
   err := c.cache.Decode(&start)
   if err != nil {
      return err
   }
   season, err := rakuten.FetchSeason(
      c.SeasonId.Value, start.Profile.Classification, start.Market,
   )
   if err != nil {
      return err
   }
   for i, episode := range season.Episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}
