package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

type saved_state struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *hboMax.St
}

func (c *client) do_login() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   state.Login, err = state.St.Login()
   if err != nil {
      return err
   }
   return c.cache.Set(state)
}

func (c *client) do() error {
   c.job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   c.job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := c.cache.Init("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   // 1
   flag.BoolVar(&c.initiate, "i", false, "device initiate")
   flag.StringVar(
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
   )
   // 2
   flag.BoolVar(&c.login, "l", false, "device login")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   flag.IntVar(&c.season, "s", 0, "season")
   // 4
   flag.StringVar(&c.edit, "e", "", "edit ID")
   // 5
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.initiate {
      return c.do_initiate()
   }
   if c.login {
      return c.do_login()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.edit != "" {
      return c.do_edit()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
}

func (c *client) do_address() error {
   var show hboMax.ShowKey
   err := show.Parse(c.address)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Get(&state)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if c.season >= 1 {
      videos, err = state.Login.Season(&show, c.season)
   } else {
      videos, err = state.Login.Movie(&show)
   }
   if err != nil {
      return err
   }
   videos.FilterAndSort()
   for i, video := range videos.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

func (c *client) do_edit() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   state.Playback, err = state.Login.PlayReady(c.edit)
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

func (c *client) do_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   initiate, err := st.Initiate(c.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return c.cache.Set(saved_state{St: &st})
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Playback.PlayReady
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

type client struct {
   cache maya.Cache
   // 1
   initiate bool
   market   string
   // 2
   login bool
   // 3
   address string
   season  int
   // 4
   edit string
   // 5
   dash string
   job  maya.PlayReadyJob
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
