package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_series() error {
   series, err := c.Client.Series(c.series)
   if err != nil {
      return err
   }
   seasons, err := series.Seasons()
   if err != nil {
      return err
   }
   for i, season := range seasons {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(season)
   }
   return nil
}

func (c *client) do_season() error {
   season, err := c.Client.Season(c.season)
   if err != nil {
      return err
   }
   episodes, err := season.Episodes()
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(episode)
   }
   return nil
}

func (c *client) do_episode() error {
   sources, header, err := c.Client.Playback(c.episode)
   if err != nil {
      return err
   }
   c.Source, err = amc.GetDash(sources)
   if err != nil {
      return err
   }
   c.BcJwt = amc.BcJwt(header)
   c.Dash, err = c.Source.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   BcJwt  string
   Client *amc.Client
   Dash   *amc.Dash
   Source *amc.Source
   //------------------------
   Job maya.Job
   //------------------------
   email    string
   password string
   //------------------------
   series int
   //------------------------
   season int
   //------------------------
   episode int
   //------------------------
   dash_id string
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Source.Widevine(c.BcJwt, data)
      },
   )
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4f")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
func (c *client) do() error {
   err := cache.Setup("rosso/amc.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringVar(&c.email, "E", "email")
   password := maya.StringVar(&c.password, "P", "password")
   //----------------------------------------------------------
   refresh := maya.BoolVar(new(bool), "r", "refresh")
   //----------------------------------------------------------
   series := maya.IntVar(&c.series, "s", "series ID")
   //----------------------------------------------------------
   season := maya.IntVar(&c.season, "S", "season ID")
   //----------------------------------------------------------
   episode := maya.IntVar(&c.episode, "e", "episode or movie ID")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[widevine] {
      return cache.Write(c)
   }
   if set[email] {
      if set[password] {
         return c.do_email_password()
      }
   }
   if set[refresh] {
      return with_cache(c.do_refresh)
   }
   if set[series] {
      return with_cache(c.do_series)
   }
   if set[season] {
      return with_cache(c.do_season)
   }
   if set[episode] {
      return with_cache(c.do_episode)
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {email, password},
      {refresh},
      {series},
      {season},
      {episode},
      {dash_id},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Client, err = amc.Unauth()
   if err != nil {
      return err
   }
   err = c.Client.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_refresh() error {
   err := c.Client.Refresh()
   if err != nil {
      return err
   }
   return cache.Write(c)
}
