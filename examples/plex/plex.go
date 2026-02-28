package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "flag"
   "log"
   "net/http"
   "path"
)

type saved_state struct {
   Dash      *plex.Dash
   MediaPart *plex.MediaPart
   User      *plex.User
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

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return state.User.Widevine(state.MediaPart, data)
   }
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/plex.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.x_forwarded_for, "x", "", "x-forwarded-for")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"a", "x"},
      {"d", "c", "p"},
   })
}

type client struct {
   cache maya.Cache
   // 1
   address         string
   x_forwarded_for string
   // 2
   dash string
   job  maya.WidevineJob
}

func (c *client) do_address() error {
   var user plex.User
   err := user.Fetch()
   if err != nil {
      return err
   }
   address, err := plex.GetPath(c.address)
   if err != nil {
      return err
   }
   metadata, err := user.RatingKey(address)
   if err != nil {
      return err
   }
   metadata, err = user.Media(metadata, c.x_forwarded_for)
   if err != nil {
      return err
   }
   var state saved_state
   state.MediaPart, err = metadata.Dash()
   if err != nil {
      return err
   }
   state.Dash, err = user.Dash(state.MediaPart, c.x_forwarded_for)
   if err != nil {
      return err
   }
   state.User = &user
   err = c.cache.Set(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}
