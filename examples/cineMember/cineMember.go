package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) do_dash() error {
   var dash cineMember.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) run() error {
   c.cache.Init("cineMember")
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
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
      {"d"},
   })
}

func (c *command) do_address() error {
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   var session cineMember.Session
   err = c.cache.Get("Session", &session)
   if err != nil {
      return err
   }
   stream, err := session.Stream(id)
   if err != nil {
      return err
   }
   link, err := stream.Dash()
   if err != nil {
      return err
   }
   dash, err := link.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) do_email_password() error {
   var session cineMember.Session
   err := session.Fetch()
   if err != nil {
      return err
   }
   err = session.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Session", session)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.Job
}
