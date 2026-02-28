package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *client) do_language() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch state.Media.Type {
   case rakuten.MovieType:
      stream, err = state.Media.MovieStream(
         c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case rakuten.TvShowType:
      stream, err = state.Media.EpisodeStream(
         c.episode, c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   state.Dash, err = stream.Dash()
   if err != nil {
      return err
   }
   state.Episode = c.episode
   state.Language = c.language
   err = c.cache.Set(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

type saved_state struct {
   Dash     *rakuten.Dash
   Episode  string
   Language string
   Media    *rakuten.Media
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch state.Media.Type {
   case rakuten.MovieType:
      stream, err = state.Media.MovieStream(
         state.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   case rakuten.TvShowType:
      stream, err = state.Media.EpisodeStream(
         state.Episode,
         state.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   c.job.Send = stream.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do_season() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   season, err := state.Media.RequestSeason(c.season)
   if err != nil {
      return err
   }
   for i, item := range season.Episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&item)
   }
   return nil
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      // everything needs proxy
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return "", false
      }
      return "", true
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache maya.Cache
   // 1
   address string
   // 2
   season string
   // 3
   episode  string
   language string
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.season, "s", "", "season ID")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.language, "A", "", "audio language")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.season != "" {
      return c.do_season()
   }
   if c.language != "" {
      return c.do_language()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"a"},
      {"s"},
      {"e", "A"},
      {"d", "c", "p"},
   })
}

func (c *client) do_address() error {
   var media rakuten.Media
   err := media.ParseURL(c.address)
   if err != nil {
      return err
   }
   switch media.Type {
   case rakuten.MovieType:
      item, err := media.RequestMovie()
      if err != nil {
         return err
      }
      fmt.Println(item)
   case rakuten.TvShowType:
      item, err := media.RequestTvShow()
      if err != nil {
         return err
      }
      fmt.Println(item)
   }
   return c.cache.Set(saved_state{Media: &media})
}
