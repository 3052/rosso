package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/mubi.xml")
   if err != nil {
      return err
   }
   // 1
   flag.BoolVar(&c.code, "c", false, "link code")
   // 2
   flag.BoolVar(&c.session, "s", false, "session")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.code {
      return c.do_code()
   }
   if c.session {
      return c.do_session()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"a"},
      {"d", "C", "P"},
   })
}

type client struct {
   cache maya.Cache
   // 1
   code bool
   // 2
   session bool
   // 3
   address string
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
   c.job.Send = state.Session.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

type saved_state struct {
   Dash *mubi.Dash
   LinkCode *mubi.LinkCode
   Session *mubi.Session
}

func (c *client) do_code() error {
   var link_code mubi.LinkCode
   err := link_code.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(&link_code)
   return c.cache.Set(saved_state{LinkCode: &link_code})
}

func (c *client) do_session() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   state.Session, err = state.LinkCode.Session()
   if err != nil {
      return err
   }
   return c.cache.Set(state)
}
func (c *client) do_address() error {
   slug, err := mubi.FilmSlug(c.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FetchId(slug)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Get(&state)
   if err != nil {
      return err
   }
   err = state.Session.Viewing(film_id)
   if err != nil {
      return err
   }
   secure, err := state.Session.SecureUrl(film_id)
   if err != nil {
      return err
   }
   state.Dash, err = secure.Dash()
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
      return "", path.Ext(req.URL.Path) != ".dash"
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
