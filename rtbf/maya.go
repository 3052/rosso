package rtbf

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func (a *Account) Identity() (*Identity, error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.getJWT", url.Values{
         "APIKey":      {api_key},
         "login_token": {a.SessionInfo.CookieValue},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Identity
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.ErrorMessage != "" {
      return nil, errors.New(result.ErrorMessage)
   }
   return &result, nil
}

func FetchAssetId(path string) (string, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "bff-service.rtbf.be",
         Path:   "/auvio/v1.23/pages" + path,
      },
      nil,
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return "", errors.New(resp.Status)
   }
   var page struct {
      Data struct {
         Content struct {
            AssetId string
            Media   *struct {
               AssetId string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&page)
   if err != nil {
      return "", err
   }
   content := page.Data.Content
   if content.AssetId != "" {
      return content.AssetId, nil
   }
   if content.Media != nil {
      return content.Media.AssetId, nil
   }
   return "", errors.New("assetId not found")
}

func (i *Identity) Session() (*Session, error) {
   body, err := json.Marshal(map[string]any{
      "device": map[string]string{
         "deviceId": "",
         "type":     "WEB",
      },
      "jwt": i.IdToken,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "exposure.api.redbee.live",
         Path:   "/v2/customer/RTBF/businessunit/Auvio/auth/gigyaLogin",
      },
      map[string]string{"content-type": "application/json"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Session{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (e *Entitlement) FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "exposure.api.redbee.live",
         Path:   "/v2/license/customer/RTBF/businessunit/Auvio/widevine",
         RawQuery: url.Values{
            "contentId":  {e.AssetId},
            "ls_session": {e.PlayToken},
         }.Encode(),
      },
      map[string]string{"content-type": "application/x-protobuf"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      var value struct {
         Message string
      }
      err = json.NewDecoder(resp.Body).Decode(&value)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(value.Message)
   }
   return io.ReadAll(resp.Body)
}
