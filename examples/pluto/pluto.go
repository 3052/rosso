package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/pluto.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.movie, "m", "", "movie URL")
   // 2
   flag.StringVar(&c.show, "s", "", "show URL")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.movie != "" {
      return c.do_movie()
   }
   if c.show != "" {
      return c.do_show()
   }
   if c.episode != "" {
      return c.do_episode()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"m"},
      {"s"},
      {"e"},
      {"d", "c", "p"},
   })
}

func (c *client) do_movie() error {
   var series pluto.Series
   err := series.Fetch(path.Base(c.movie))
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(series.GetMovieURL())
   if err != nil {
      return err
   }
   err = c.cache.Set(saved_state{Dash: &dash})
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *client) do_show() error {
   var series pluto.Series
   err := series.Fetch(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return c.cache.Set(saved_state{Series: &series})
}
type saved_state struct {
   Dash *pluto.Dash
   Series *pluto.Series
}

func (c *client) do_episode() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   link, err := state.Series.GetEpisodeURL(c.episode)
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(link)
   if err != nil {
      return err
   }
   state.Dash = &dash
   err = c.cache.Set(state)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = pluto.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache maya.Cache
   // 1
   movie string
   // 2
   show string
   // 3
   episode string
   // 4
   dash string
   job  maya.WidevineJob
}
