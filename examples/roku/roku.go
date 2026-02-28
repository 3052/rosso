package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *client) do_connection() error {
   connection, err := roku.NewConnection(nil)
   if err != nil {
      return err
   }
   state := saved_state{Connection: connection}
   state.LinkCode, err = state.Connection.LinkCode()
   if err != nil {
      return err
   }
   fmt.Println(state.LinkCode)
   return c.cache.Set(state)
}

type saved_state struct {
   Dash       *roku.Dash
   Connection *roku.Connection
   LinkCode   *roku.LinkCode
   Playback   *roku.Playback
   User       *roku.User
}

func (c *client) do_set_user() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   state.User, err = state.Connection.User(state.LinkCode)
   if err != nil {
      return err
   }
   return c.cache.Set(state)
}

func (c *client) do_roku() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   var connection *roku.Connection
   if c.get_user {
      connection, err = roku.NewConnection(state.User)
   } else {
      connection, err = roku.NewConnection(nil)
   }
   if err != nil {
      return err
   }
   state.Playback, err = connection.Playback(c.roku)
   if err != nil {
      return err
   }
   state.Dash, err = state.Playback.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
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
   connection bool
   // 2
   set_user bool
   // 3
   roku     string
   get_user bool
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Playback.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/roku.xml")
   if err != nil {
      return err
   }
   // 1
   flag.BoolVar(&c.connection, "c", false, "connection")
   // 2
   flag.BoolVar(&c.set_user, "s", false, "set user")
   // 3
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_user, "g", false, "get user")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.connection {
      return c.do_connection()
   }
   if c.set_user {
      return c.do_set_user()
   }
   if c.roku != "" {
      return c.do_roku()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"r", "g"},
      {"d", "C", "P"},
   })
}
