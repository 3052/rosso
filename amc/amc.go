package amc

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Metadata struct {
   EpisodeNumber int
   Nid           int
   Title         string
}

func (m *Metadata) String() string {
   var data strings.Builder
   if m.EpisodeNumber >= 0 {
      data.WriteString("episode = ")
      data.WriteString(strconv.Itoa(m.EpisodeNumber))
   }
   if data.Len() >= 1 {
      data.WriteByte('\n')
   }
   data.WriteString("title = ")
   data.WriteString(m.Title)
   data.WriteString("\nnid = ")
   data.WriteString(strconv.Itoa(m.Nid))
   return data.String()
}

func BcJwt(header http.Header) string {
   return header.Get("x-amcn-bc-jwt")
}

func (c *Client) Refresh() error {
   var req http.Request
   req.Method = "POST"
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+c.Data.RefreshToken)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "gw.cds.amcn.com",
      Path:   "/auth-orchestration-id/api/v1/refresh",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(c)
}

func Unauth() (*Client, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "gw.cds.amcn.com",
      Path:   "/auth-orchestration-id/api/v1/unauth",
   }
   req.Header = http.Header{}
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-tenant", "amcn")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Client{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Client struct {
   Data struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
}

func (c *Client) Login(email, password string) error {
   data, err := json.Marshal(map[string]string{
      "email":    email,
      "password": password,
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://gw.cds.amcn.com", bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.URL.Path = "/auth-orchestration-id/api/v1/login"
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-group-id", "10")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(c)
}

func (c *Client) Playback(id int) ([]Source, http.Header, error) {
   data, err := json.Marshal(map[string]any{
      "adtags": map[string]any{
         "lat":          0,
         "mode":         "on-demand",
         "playerHeight": 0,
         "playerWidth":  0,
         "ppid":         0,
         "url":          "-",
      },
   })
   if err != nil {
      return nil, nil, err
   }
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "gw.cds.amcn.com",
      Path:   "/playback-id/api/v1/playback/" + strconv.Itoa(id),
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+c.Data.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")
   req.Body = io.NopCloser(bytes.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []Source
         }
      }
      Error string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, nil, err
   }
   if result.Error != "" {
      return nil, nil, errors.New(result.Error)
   }
   return result.Data.PlaybackJsonData.Sources, resp.Header, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

func (s *Source) Widevine(bcJwt string, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", s.KeySystems.ComWidevineAlpha.LicenseUrl,
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("bcov-auth", bcJwt)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Source struct {
   KeySystems struct {
      ComWidevineAlpha *struct {
         LicenseUrl string `json:"license_url"`
      } `json:"com.widevine.alpha"`
   } `json:"key_systems"`
   Src  string // URL to the MPD manifest
   Type string // e.g., "application/dash+xml"
}

func (s *Source) Dash() (*Dash, error) {
   resp, err := http.Get(s.Src)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

func GetDash(sources []Source) (*Source, error) {
   for _, source_data := range sources {
      if source_data.Type == "application/dash+xml" {
         return &source_data, nil
      }
   }
   return nil, errors.New("DASH source not found")
}

func join(data ...string) string {
   return strings.Join(data, "")
}
