package hboMax

import (
   "bytes"
   "cmp"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strings"
)

func (p *Playback) GetManifest() (*url.URL, error) {
   return url.Parse(strings.Replace(p.Fallback.Manifest.Url, "_fallback", "", 1))
}

type Login struct {
   Token string
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL = ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code = ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

var Markets = []string{
   "amer",
   "apac",
   "emea",
   "latam",
}

func MovieResults(entities []*Entity) []*Entity {
   var movies []*Entity
   for _, item := range entities {
      // Identify the primary video entity for the movie
      if item.Type == "video" && item.Attributes.VideoType == "MOVIE" {
         movies = append(movies, item)
      }
   }
   return movies
}

func SeasonResults(entities []*Entity) []*Entity {
   var results []*Entity
   for _, item := range entities {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         results = append(results, item)
      }
   }
   // Sort episodes by EpisodeNumber using the modern slices.SortFunc
   slices.SortFunc(results, func(entityA, entityB *Entity) int {
      return cmp.Compare(entityA.Attributes.EpisodeNumber, entityB.Attributes.EpisodeNumber)
   })
   return results
}

// Entity represents a single unified node in the Max API response
type Entity struct {
   Attributes struct {
      Name          string
      Alias         string
      ShowType      string
      VideoType     string
      MaterialType  string
      Description   string
      SeasonNumber  int
      EpisodeNumber int
      AirDate       string
   }
   Id            string
   Relationships struct {
      Edit struct {
         Data Resource
      }
      Items struct {
         Data []Resource
      }
      Show struct {
         Data Resource
      }
      Video struct {
         Data Resource
      }
   }
   Type string
}

func SearchResults(entities []*Entity) ([]*Entity, error) {
   // Pre-allocate map capacity for better performance
   entitiesMap := make(map[string]*Entity, len(entities))
   var searchResultsCollection *Entity

   // Combine map building and target collection searching into a single loop
   for _, each := range entities {
      entitiesMap[each.Id] = each

      if searchResultsCollection == nil && each.Type == "collection" && each.Attributes.Alias == "search-page-rail-results" {
         searchResultsCollection = each
      }
   }

   if searchResultsCollection == nil {
      return nil, fmt.Errorf("could not find the search results collection in the response payload")
   }

   var results []*Entity
   for _, itemRes := range searchResultsCollection.Relationships.Items.Data {
      colItem, exists := entitiesMap[itemRes.Id]
      if !exists {
         continue
      }

      targetId := colItem.Relationships.Show.Data.Id
      if targetId == "" {
         targetId = colItem.Relationships.Video.Data.Id
      }

      if targetId == "" {
         continue
      }

      mediaEntity, exists := entitiesMap[targetId]
      if !exists {
         continue
      }

      // Append the actual show/movie entity
      results = append(results, mediaEntity)
   }
   return results, nil
}

// Resource represents a relationship pointer in the JSON:API graph
type Resource struct {
   Id   string
   Type string
}

type Scheme struct {
   LicenseUrl string
}

const (
   disco_client = "!:!:beam:!"
   device_info  = "!/!(!/!;!/!;!/!)"
)

// String implements the fmt.Stringer interface to provide a clean visual
// output for the Entity
func (e *Entity) String() string {
   data := &strings.Builder{}
   if e.Attributes.MaterialType == "EPISODE" {
      fmt.Fprintf(data, "Episode: %d\n", e.Attributes.EpisodeNumber)
   }
   if e.Attributes.ShowType != "" {
      fmt.Fprintf(data, "Show Type: %s\n", e.Attributes.ShowType)
   } else if e.Attributes.VideoType != "" {
      fmt.Fprintf(data, "Video Type: %s\n", e.Attributes.VideoType)
   }
   fmt.Fprintf(data, "Name: %s\n", e.Attributes.Name)
   if e.Type == "video" {
      fmt.Fprintf(data, "Edit ID: %s\n", e.Relationships.Edit.Data.Id)
   } else {
      fmt.Fprintf(data, "ID: %s\n", e.Id)
   }
   return strings.TrimSpace(data.String())
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Errors   []Error
   Fallback struct {
      Manifest struct {
         Url string // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

func (e *Error) Error() string {
   var data strings.Builder
   // 1. print code
   data.WriteString("code = ")
   data.WriteString(e.Code)
   // 2, 3, 4. if detail print detail, if message print message, if both print
   // one
   if e.Detail != "" {
      data.WriteString("\ndetail = ")
      data.WriteString(e.Detail)
   } else if e.Message != "" {
      data.WriteString("\nmessage = ")
      data.WriteString(e.Message)
   }
   return data.String()
}

type Error struct {
   Code    string // 2026-04-10
   Detail  string // 2026-04-10
   Message string // 2026-04-10
}

func playback_request(token, edit_id, drm string) (*Playback, error) {
   body, err := json.Marshal(map[string]any{
      "editId":               edit_id,
      "consumptionType":      "streaming",
      "appBundle":            "",         // required
      "applicationSessionId": "",         // required
      "firstPlay":            false,      // required
      "gdpr":                 false,      // required
      "playbackSessionId":    "",         // required
      "userPreferences":      struct{}{}, // required
      "capabilities": map[string]any{
         "contentProtection": map[string]any{
            "contentDecryptionModules": []any{
               map[string]string{
                  "drmKeySystem": drm,
               },
            },
         },
         "manifests": map[string]any{
            "formats": map[string]any{
               "dash": struct{}{}, // required
            }, // required
         }, // required
      }, // required
      "deviceInfo": map[string]any{
         "player": map[string]any{
            "mediaEngine": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
            "playerView": map[string]int{
               "height": 0, // required
               "width":  0, // required
            }, // required
            "sdk": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
         }, // required
      }, // required
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://default.prd.api.hbomax.com", bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo"
   req.Header.Set("authorization", "Bearer "+token)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

func PlayReadyRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "playready")
}

func WidevineRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "widevine")
}

func entity_request(token string, endpoint *url.URL) ([]*Entity, error) {
   // Scheme
   endpoint.Scheme = "https"
   // Host
   endpoint.Host = "default.prd.api.hbomax.com"
   // RawQuery
   query := endpoint.Query()
   query.Set("include", "default")
   endpoint.RawQuery = query.Encode()
   req := http.Request{
      URL:    endpoint,
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Errors   []Error
      Included []*Entity `json:"included"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return result.Included, nil
}

func MovieRequest(token, showId string) ([]*Entity, error) {
   values := url.Values{}
   values.Set("page[items.size]", "1")
   parsedUrl := &url.URL{
      Path:     "/cms/routes/movie/" + showId,
      RawQuery: values.Encode(),
   }
   return entity_request(token, parsedUrl)
}

func SeasonRequest(token, showId string, seasonNumber int) ([]*Entity, error) {
   values := url.Values{}
   values.Set("pf[show.id]", showId)
   values.Set("pf[seasonNumber]", fmt.Sprint(seasonNumber))
   parsedUrl := &url.URL{
      Path:     "/cms/collections/generic-show-page-rail-episodes-tabbed-content",
      RawQuery: values.Encode(),
   }
   return entity_request(token, parsedUrl)
}

func SearchRequest(token, query string) ([]*Entity, error) {
   values := url.Values{}
   values.Set("contentFilter[query]", query)
   parsedUrl := &url.URL{
      Path:     "/cms/routes/search/result",
      RawQuery: values.Encode(),
   }
   return entity_request(token, parsedUrl)
}

func (p *Playback) WidevineRequest(body []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.Widevine.LicenseUrl, "application/x-protobuf",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}
