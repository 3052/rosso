package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/canal"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
)

func (c *command) do_tracking() error {
   var session canal.Session
   err := c.cache.Get("Session", &session)
   if err != nil {
      return err
   }
   player, err := session.Player(c.tracking)
   if err != nil {
      return err
   }
   err = c.cache.Set("Player", player)
   if err != nil {
      return err
   }
   dash, err := player.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("canal")
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.tracking, "t", "", "tracking")
   flag.IntVar(&c.season, "s", 0, "season")
   // 5
   flag.BoolVar(&c.subtitles, "S", false, "subtitles")
   // 6
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.tracking != "" {
      if c.season >= 1 {
         return c.do_tracking_season()
      }
      return c.do_tracking()
   }
   if c.subtitles {
      return c.do_subtitles()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"r"},
      {"a"},
      {"t", "s"},
      {"S"},
      {"d", "C", "P"},
   })
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".dash"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_dash() error {
   var player canal.Player
   err := c.cache.Get("Player", &player)
   if err != nil {
      return err
   }
   c.job.Send = player.Widevine
   var dash canal.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func get(address string) error {
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   file, err := os.Create(path.Base(address))
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   if err != nil {
      return err
   }
   return nil
}

func (c *command) do_email_password() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   login, err := ticket.Login(c.email, c.password)
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Fetch(login.SsoToken)
   if err != nil {
      return err
   }
   return c.cache.Set("Session", session)
}

func (c *command) do_address() error {
   tracking, err := canal.FetchTracking(c.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking =", tracking)
   return nil
}

func (c *command) do_refresh() error {
   var session canal.Session
   err := c.cache.Get("Session", &session)
   if err != nil {
      return err
   }
   err = session.Fetch(session.SsoToken)
   if err != nil {
      return err
   }
   return c.cache.Set("Session", session)
}

func (c *command) do_tracking_season() error {
   var session canal.Session
   err := c.cache.Get("Session", &session)
   if err != nil {
      return err
   }
   episodes, err := session.Episodes(c.tracking, c.season)
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}

type command struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   refresh bool
   // 3
   address string
   // 4
   tracking string
   season   int
   // 5
   subtitles bool
   // 6
   dash string
   job  maya.WidevineJob
}

func (c *command) do_subtitles() error {
   var player canal.Player
   err := c.cache.Get("Player", &player)
   if err != nil {
      return err
   }
   for _, subtitles := range player.Subtitles {
      err = get(subtitles.Url)
      if err != nil {
         return err
      }
   }
   return nil
}
