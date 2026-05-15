package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

type client struct {
   address   string
   audio     string
   cache     maya.Cache
   dash      string
   episode   string
   flag      maya.FlagSet
   season    string
   playReady string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rakuten"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   audio := c.flag.String(&c.audio, "A", "audio language")
   episode := c.flag.String(&c.episode, "e", "episode ID")
   season := c.flag.String(&c.season, "s", "season ID")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   playReady := c.flag.String(&c.playReady, "p", "PlayReady")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case playReady.IsSet:
      return c.cache.Encode(playReady_device(c.playReady))
   case address.IsSet:
      return c.do_address()
   case season.IsSet:
      return c.do_season()
   case audio.IsSet:
      return c.do_audio()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {playReady},
      {address},
      {season},
      {audio, episode},
      {dash},
   })
}

type playReady_device string

func (c *client) do_dash() error {
   var (
      device      playReady_device
      manifest    maya.Manifest
      stream_info rakuten.StreamInfo
   )
   err := c.cache.Decode(&device, &manifest, &stream_info)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmPlayReady,
      License: stream_info.FetchLicense,
   })
}

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
