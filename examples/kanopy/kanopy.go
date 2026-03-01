package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do_kanopy() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   member, err := state.Login.Membership()
   if err != nil {
      return err
   }
   plays, err := state.Login.Plays(member, c.kanopy)
   if err != nil {
      return err
   }
   state.PlayManifest, err = plays.Dash()
   if err != nil {
      return err
   }
   state.Dash, err = state.PlayManifest.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set(&state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func main() {
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache maya.Cache
   // 1
   proxy string
   // 2
   email    string
   password string
   // 3
   kanopy int
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/kanopy.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 3
   flag.IntVar(&c.kanopy, "k", 0, "Kanopy ID")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return c.proxy, path.Ext(req.URL.Path) != ".m4s"
   })
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.kanopy >= 1 {
      return c.do_kanopy()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"k"},
      {"d", "C", "P"},
   })
}

func (c *client) do_email_password() error {
   var login kanopy.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Login: &login})
}

type saved_state struct {
   Dash         *kanopy.Dash
   Login        *kanopy.Login
   PlayManifest *kanopy.PlayManifest
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return state.Login.Widevine(state.PlayManifest, data)
   }
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}
