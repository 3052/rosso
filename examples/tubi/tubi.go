package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/tubi"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do_tubi() error {
   var content tubi.Content
   err := content.Fetch(c.tubi)
   if err != nil {
      return err
   }
   var state saved_state
   state.VideoResource = &content.VideoResources[0]
   state.Dash, err = state.VideoResource.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

type saved_state struct {
   Dash *tubi.Dash
   VideoResource *tubi.VideoResource
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.VideoResource.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".mp4"
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache maya.Cache
   // 1
   tubi  int
   // 2
   dash string
   job  maya.WidevineJob
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/tubi.xml")
   if err != nil {
      return err
   }
   // 1
   flag.IntVar(&c.tubi, "t", 0, "Tubi ID")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.tubi >= 1 {
      return c.do_tubi()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"t"},
      {"d", "c", "p"},
   })
}
