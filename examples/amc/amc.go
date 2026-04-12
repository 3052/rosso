package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "fmt"
   "log"
)

func main() {
   maya.SetProxy("", "*.m4f")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   BcJwt  string
   AuthData *amc.AuthData
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

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Source.Widevine(c.BcJwt, data)
      },
   )
}
func (c *client) do() error {
   err := cache.Setup("rosso/amc.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "E", "email")
   password := maya.StringFlag(&c.password, "P", "password")
   //----------------------------------------------------------
   refresh := maya.BoolFlag("r", "refresh")
   //----------------------------------------------------------
   series := maya.IntFlag(&c.series, "s", "series ID")
   //----------------------------------------------------------
   season := maya.IntFlag(&c.season, "S", "season ID")
   //----------------------------------------------------------
   episode := maya.IntFlag(&c.episode, "e", "episode or movie ID")
   //----------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if refresh.IsSet {
      return with_cache(c.do_refresh)
   }
   if series.IsSet {
      return with_cache(c.do_series)
   }
   if season.IsSet {
      return with_cache(c.do_season)
   }
   if episode.IsSet {
      return with_cache(c.do_episode)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {refresh},
      {series},
      {season},
      {episode},
      {dash_id},
   })
}

var cache maya.Cache
