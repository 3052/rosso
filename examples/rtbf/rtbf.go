package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "flag"
   "log"
   "net/http"
)

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Entitlement.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/rtbf.xml")
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

func (c *client) do_address() error {
   path, err := rtbf.GetPath(c.address)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Get(&state)
   if err != nil {
      return err
   }
   identity, err := state.Account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   state.Entitlement, err = session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   format, err := state.Entitlement.Dash()
   if err != nil {
      return err
   }
   state.Dash, err = format.Dash()
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
   Account *rtbf.Account
   Dash *rtbf.Dash
   Entitlement *rtbf.Entitlement
}

func (c *client) do_email_password() error {
   var account rtbf.Account
   err := account.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Account: &account})
}

func main() {
   maya.SetProxy(func(*http.Request) (string, bool) {
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
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.WidevineJob
}
