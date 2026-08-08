package amc

import (
   _ "embed"
   "encoding/json"
   "fmt"
   "net/http"
)

//go:embed playback.json
var playback_json []byte

type ContentNode struct {
   Type             string        `json:"type"`
   Properties       *Properties   `json:"properties,omitempty"`
   TabletProperties *Properties   `json:"tablet_properties,omitempty"`
   Children         []ContentNode `json:"children,omitempty"`
}

func SeasonEpisodes(authToken string, id int) (*ContentNode, error) {
   resp, err := doGet(
      fmt.Sprint("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/season-episodes/id/", id),
      map[string]string{
         "authorization":   "Bearer " + authToken,
         "x-amcn-network":  "amcplus",
         "x-amcn-platform": "android",
         "x-amcn-tenant":   "amcn",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("season episodes failed with status: %d", resp.StatusCode)
   }
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

func SeriesDetail(authToken string, id int) (*ContentNode, error) {
   resp, err := doGet(
      fmt.Sprint("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/series-detail/id/", id),
      map[string]string{
         "authorization":   "Bearer " + authToken,
         "x-amcn-network":  "amcplus",
         "x-amcn-platform": "android",
         "x-amcn-tenant":   "amcn",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("series detail failed with status: %d", resp.StatusCode)
   }
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

func (c *ContentNode) EpisodesMetadata() []*Metadata {
   var metadata []*Metadata
   var walk func(node ContentNode)
   walk = func(node ContentNode) {
      props := node.Properties
      if props != nil && props.Metadata != nil {
         if node.Type == "card" && props.ContentType == "episode" {
            metadata = append(metadata, props.Metadata)
         }
      }
      for _, child := range node.Children {
         walk(child)
      }
   }
   walk(*c)
   return metadata
}

func (c *ContentNode) SeasonsMetadata() []*Metadata {
   var metadata []*Metadata
   var walk func(node ContentNode)
   walk = func(node ContentNode) {
      props := node.Properties
      if props != nil && props.Metadata != nil {
         if node.Type == "tab_bar_item" && props.Metadata.SeasonNumber > 0 {
            metadata = append(metadata, props.Metadata)
         }
      }
      for _, child := range node.Children {
         walk(child)
      }
   }
   walk(*c)
   return metadata
}

type Playback struct {
   BcovAuth string
   Sources  []Source
}

func GetPlayback(authToken string, videoId int) (*Playback, error) {
   resp, err := doPost(
      fmt.Sprint("https://gw.cds.amcn.com/playback-id/api/v1/playback/", videoId),
      map[string]string{
         "authorization":       "Bearer " + authToken,
         "content-type":        "application/json",
         "x-amcn-language":     "en",
         "x-amcn-network":      "amcplus",
         "x-amcn-platform":     "web",
         "x-amcn-service-id":   "amcplus",
         "x-amcn-tenant":       "amcn",
         "x-amcn-device-ad-id": "-",
         "x-ccpa-do-not-sell":  "doNotPassData",
      },
      playback_json,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("playback failed with status: %d", resp.StatusCode)
   }
   var result struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []Source
         }
      }
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &Playback{
      BcovAuth: resp.Header.Get("x-amcn-bc-jwt"),
      Sources:  result.Data.PlaybackJsonData.Sources,
   }, nil
}

func (*Playback) CachePath() string {
   return "rosso/amc/Playback"
}

func (p *Playback) GetDash() (*Source, error) {
   for _, source_data := range p.Sources {
      if source_data.Type == "application/dash+xml" {
         return &source_data, nil
      }
   }
   return nil, fmt.Errorf("application/dash+xml source not found")
}

type Source struct {
   Codecs     string
   KeySystems KeySystems `json:"key_systems"`
   Src        string     // MPD
   Type       string
}

func (*Source) CachePath() string {
   return "rosso/amc/Source"
}

// content.go
