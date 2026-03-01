package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Playout.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

type saved_state struct {
   Cookie  *http.Cookie
   Dash    *peacock.Dash
   Playout *peacock.Playout
}

func (c *client) do_address() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   var token peacock.Token
   err = token.Fetch(state.Cookie)
   if err != nil {
      return err
   }
   state.Playout, err = token.Playout(path.Base(c.address))
   if err != nil {
      return err
   }
   endpoint, err := state.Playout.Fastly()
   if err != nil {
      return err
   }
   state.Dash, err = endpoint.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_email_password() error {
   cookie, err := peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Cookie: cookie})
}

type client struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/peacock.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
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
