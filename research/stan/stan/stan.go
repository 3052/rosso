package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/stan"
   "fmt"
   "log"
   "os"
   "strconv"
)

func indexByte(s string, b byte) int {
   for i := 0; i < len(s); i++ {
      if s[i] == b {
         return i
      }
   }
   return -1
}

// helpers

func lastSlash(s string) int {
   for i := len(s) - 1; i >= 0; i-- {
      if s[i] == '/' {
         return i
      }
   }
   return -1
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine  maya.FlagString
   address   maya.FlagString
   dash      maya.FlagString
   link_code maya.FlagBool
   stan_id   maya.FlagInt
   session   maya.FlagBool

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/stan/client"
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "link-code", Value: &c.link_code},
      {Name: "session", Value: &c.session},
      {Name: "address", Value: &c.address, Usage: "program URL"},
      {Name: "stan-id", Value: &c.stan_id},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.link_code {
      return c.do_link_code()
   }
   if c.session {
      return c.do_session()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.stan_id >= 1 {
      return c.do_stan()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "stan")
}

func (c *client) do_address() error {
   // Extract the program ID from the URL path
   // Stan URLs look like https://www.stan.com.au/program/123456
   base := string(c.address)
   // get the last path segment
   idStr := base
   if i := lastSlash(base); i >= 0 {
      idStr = base[i+1:]
   }
   // strip any query string
   if i := indexByte(idStr, '?'); i >= 0 {
      idStr = idStr[:i]
   }
   id, err := strconv.ParseInt(idStr, 10, 64)
   if err != nil {
      return err
   }
   var program stan.LegacyProgram
   err = program.New(id)
   if err != nil {
      return err
   }
   fmt.Println(program)
   return nil
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      stream   stan.ProgramStream
   )
   err := c.cache.Decode(&manifest, &stream)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: stream.License,
   })
}

func (c *client) do_link_code() error {
   activation_code, err := stan.FetchActivationCode()
   if err != nil {
      return err
   }
   fmt.Println(activation_code)
   return c.cache.Encode(activation_code)
}

func (c *client) do_session() error {
   var activation_code stan.ActivationCode
   err := c.cache.Decode(&activation_code)
   if err != nil {
      return err
   }
   web_token, err := activation_code.Token()
   if err != nil {
      return err
   }
   session, err := web_token.Session()
   if err != nil {
      return err
   }
   return c.cache.Encode(session)
}

func (c *client) do_stan() error {
   var session stan.AppSession
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   stream, err := session.Stream(int64(c.stan_id))
   if err != nil {
      return err
   }
   // Try each base URL host until one works
   var address *url.URL
   for _, host := range stan.BaseUrl {
      address, err = stream.BaseUrl(host)
      if err != nil {
         continue
      }
      // attempt to fetch the manifest to see if the host works
      resp, err := http.Get(address.String())
      if err != nil {
         continue
      }
      resp.Body.Close()
      if resp.StatusCode == 200 {
         break
      }
   }
   if address == nil {
      return fmt.Errorf("no working base URL found")
   }
   manifest, err := maya.ListDash(address.String())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, stream)
}

// examples/stan/stan.go
