package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "address", Value: &c.address},
      {Name: "season-id", Value: &c.season},
      {Name: "audio-language", Value: &c.audio},
      {Name: "episode-id", Value: &c.episode, Needs: "audio-language"},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.PlayReady):
      return c.cache.Encode(c)
   case c.address != "":
      return c.do_address()
   case c.season != "":
      return c.do_season()
   case c.audio != "":
      return c.do_audio()
   case c.dash != "":
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "rakuten")
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
         address.ContentId, start.Profile.Classification, string(c.audio),
      )
   case address.IsTvShow():
      stream_info, err = rakuten.FetchEpisodeStreaming(
         string(c.episode), start.Profile.Classification, string(c.audio),
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

func (c *client) do_dash() error {
   var (
      manifest    maya.Manifest
      stream_info rakuten.StreamInfo
   )
   err := c.cache.Decode(&manifest, &stream_info)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: stream_info.FetchLicense,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   PlayReady maya.FlagString

   address maya.FlagString
   audio   maya.FlagString
   dash    maya.FlagString
   episode maya.FlagString
   season  maya.FlagString

   cache maya.Cache
}

func (c *client) do_address() error {
   address, err := rakuten.ParseAddress(string(c.address))
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

func (c *client) do_season() error {
   var start rakuten.Start
   err := c.cache.Decode(&start)
   if err != nil {
      return err
   }
   season, err := rakuten.FetchSeason(
      string(c.season), start.Profile.Classification, start.Market,
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
