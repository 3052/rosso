package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (s *session) run_address() error {
   var client client_data
   err := s.cache.Get(&client)
   if err != nil {
      return err
   }
   err = client.Account.RefreshToken()
   if err != nil {
      return err
   }
   err = s.cache.Set(client)
   if err != nil {
      return err
   }
   entity, err := disney.GetEntity(s.address)
   if err != nil {
      return err
   }
   page, err := client.Account.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (s *session) run_season_id() error {
   var client client_data
   err := s.cache.Get(&client)
   if err != nil {
      return err
   }
   season, err := client.Account.Season(s.season_id)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (s *session) run_hls() error {
   var client client_data
   err := s.cache.Get(&client)
   if err != nil {
      return err
   }
   s.job.Send = client.Account.PlayReady
   return s.job.DownloadHls(client.Hls.Body, client.Hls.Url, s.hls)
}

func (s *session) run_media_id() error {
   var client client_data
   err := s.cache.Get(&client)
   if err != nil {
      return err
   }
   stream, err := client.Account.Stream(s.media_id)
   if err != nil {
      return err
   }
   client.Hls, err = stream.Hls()
   if err != nil {
      return err
   }
   err = s.cache.Set(client)
   if err != nil {
      return err
   }
   return maya.ListHls(client.Hls.Body, client.Hls.Url)
}

func (s *session) run() error {
   s.job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   s.job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := s.cache.Init("rosso/disney.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&s.email, "e", "", "email")
   flag.StringVar(&s.password, "p", "", "password")
   // 2
   flag.StringVar(&s.profile_id, "P", "", "profile ID")
   // 3
   flag.StringVar(&s.address, "a", "", "address")
   // 4
   flag.StringVar(&s.season_id, "s", "", "season ID")
   // 5
   flag.StringVar(&s.media_id, "m", "", "media ID")
   // 6
   flag.IntVar(&s.hls, "h", -1, "HLS ID")
   flag.StringVar(&s.job.CertificateChain, "C", s.job.CertificateChain, "certificate chain")
   flag.StringVar(&s.job.EncryptSignKey, "E", s.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if s.email != "" {
      if s.password != "" {
         return s.run_email_password()
      }
   }
   if s.profile_id != "" {
      return s.run_profile_id()
   }
   if s.address != "" {
      return s.run_address()
   }
   if s.season_id != "" {
      return s.run_season_id()
   }
   if s.media_id != "" {
      return s.run_media_id()
   }
   if s.hls >= 0 {
      return s.run_hls()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"P"},
      {"a"},
      {"s"},
      {"m"},
      {"h", "s", "E"},
   })
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return "", false
      }
      return "", true
   })
   err := new(session).run()
   if err != nil {
      log.Fatal(err)
   }
}

type session struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   profile_id string
   // 3
   address string
   // 4
   season_id string
   // 5
   media_id string
   // 6
   hls int
   job maya.PlayReadyJob
}

type client_data struct {
   Account         *disney.Account
   InactiveAccount *disney.InactiveAccount
   Hls             *disney.Hls
}

func (s *session) run_email_password() error {
   device, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   var client client_data
   client.InactiveAccount, err = device.Login(s.email, s.password)
   if err != nil {
      return err
   }
   for i, profile := range client.InactiveAccount.Data.Login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return s.cache.Set(client)
}

func (s *session) run_profile_id() error {
   var client client_data
   err := s.cache.Get(&client)
   if err != nil {
      return err
   }
   client.Account, err = client.InactiveAccount.SwitchProfile(s.profile_id)
   if err != nil {
      return err
   }
   return s.cache.Set(client)
}
