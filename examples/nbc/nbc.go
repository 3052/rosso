package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/nbc"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do_address() error {
   name, err := nbc.GetName(c.address)
   if err != nil {
      return err
   }
   metadata, err := nbc.FetchMetadata(name)
   if err != nil {
      return err
   }
   stream, err := metadata.Stream()
   if err != nil {
      return err
   }
   dash, err := stream.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set(dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
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

func (c *client) do_dash() error {
   var dash nbc.Dash
   err := c.cache.Get(&dash)
   if err != nil {
      return err
   }
   c.job.Send = nbc.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

type client struct {
   cache maya.Cache
   // 1
   address string
   // 2
   dash string
   job  maya.WidevineJob
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/nbc.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
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
      {"a"},
      {"d", "c", "p"},
   })
}
