package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

type PlayReadyFolder string

func (c *client) do_dash() error {
   var (
      manifest    maya.Manifest
      playReady   PlayReadyFolder
      stream_info rakuten.StreamInfo
   )
   err := c.cache.Decode(&manifest, &playReady, &stream_info)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: stream_info.FetchLicense,
   })
}

///

func (c *client) do_address() error {
   parsed, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   address := rakuten.ParseUrl(parsed)

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

func (c *client) do_season() error {
   var start rakuten.Start
   err := c.cache.Decode(&start)
   if err != nil {
      return err
   }
   season, err := rakuten.FetchSeason(
      c.season.Value, start.Profile.Classification, start.Market,
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

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_audio() error {
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
         address.ContentId, start.Profile.Classification, c.audio.Value,
      )
   case address.IsTvShow():
      stream_info, err = rakuten.FetchEpisodeStreaming(
         c.episode.Value, start.Profile.Classification, c.audio.Value,
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

type client struct {
   cache     maya.Cache
   flag      maya.FlagSet
   address   maya.Flag
   audio     maya.Flag
   dash      maya.Flag
   episode   maya.Flag
   season    maya.Flag
   playReady maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rakuten"); err != nil {
      return err
   }
   c.flag.AddValue(&c.playReady, "p", "PlayReady")
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.season, "s", "season ID")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.audio, "A", "audio language")
   c.flag.AddValue(&c.episode, "e", "episode ID")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.playReady.Set:
      return c.cache.Encode(PlayReadyFolder(c.playReady.Value))
   case c.address.Set:
      return c.do_address()
   case c.season.Set:
      return c.do_season()
   case c.audio.Set:
      return c.do_audio()
   case c.dash.Set:
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}
