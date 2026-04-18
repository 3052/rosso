package crave

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "net/url"
)

// crave-web:default
const crave_web = "Basic Y3JhdmUtd2ViOmRlZmF1bHQ="

// 699710369328da351ac33c63
func (a *Account) Login(profileId string) error {
   body := url.Values{
      "grant_type":    {"refresh_token"},
      "profile_id":    {profileId},
      "refresh_token": {a.RefreshToken},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/login/v2.2",
      },
      map[string]string{
         "authorization": crave_web,
         "content-type":  "application/x-www-form-urlencoded",
      },
      []byte(body),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return fmt.Errorf("profile login failed with: %v", resp.Status)
   }
   return json.NewDecoder(resp.Body).Decode(a)
}

func (a *Account) FetchProfiles() ([]*Profile, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/profile/v2/account/" + a.AccountId,
      },
      map[string]string{"authorization": "Bearer " + a.AccessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("failed to fetch profiles with: %v", resp.Status)
   }
   var profiles []*Profile
   if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
      return nil, err
   }
   return profiles, nil
}

func FetchSubscriptions(accessToken string) ([]Subscription, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/subscription/v5",
      },
      map[string]string{"authorization": "Bearer " + accessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Subscriptions []Subscription
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Subscriptions, nil
}

func FetchContentPackage(accessToken string, contentId int) (*ContentPackage, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "playback.rte-api.bellmedia.ca",
         Path:   fmt.Sprint("/content/", contentId),
      },
      map[string]string{
         "authorization":       "Bearer " + accessToken,
         "x-client-platform":   "platform_jasper_web", // platform_jasper_html
         "x-playback-language": Language,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      ContentPackage ContentPackage
      Error          string // 2026-04-14
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   return &result.ContentPackage, nil
}
