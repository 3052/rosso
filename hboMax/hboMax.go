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

func (l *Login) fetch_playback(edit_id, drm string) (*Playback, error) {
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
   req.Header.Set("authorization", "Bearer "+l.Token)
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
   return &result, nil
}

const (
   disco_client = "!:!:beam:!"
   device_info  = "!/!(!/!;!/!;!/!)"
)

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func FetchLogin(st *http.Cookie) (*Login, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com", // Refactored
         Path:   "/authentication/linkDevice/login",
      },
      Header: http.Header{},
   }
   req.AddCookie(st)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Attributes Login
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}

type Login struct {
   Token string
}

// Resource represents a relationship pointer in the JSON:API graph.
type Resource struct {
   Id   string `json:"id"`
   Type string `json:"type"`
}

// Entity represents a single unified node in the Max API response.
type Entity struct {
   Attributes struct {
      Name          string `json:"name"`
      Alias         string `json:"alias"`
      ShowType      string `json:"showType"`
      VideoType     string `json:"videoType"`
      MaterialType  string `json:"materialType"`
      Description   string `json:"description"`
      SeasonNumber  int    `json:"seasonNumber"`
      EpisodeNumber int    `json:"episodeNumber"`
      AirDate       string `json:"airDate"`
   }
   Id            string `json:"id"`
   Relationships struct {
      Edit struct {
         Data Resource `json:"data"`
      } `json:"edit"`
      Items struct {
         Data []Resource `json:"data"`
      } `json:"items"`
      Show struct {
         Data Resource `json:"data"`
      } `json:"show"`
      Video struct {
         Data Resource `json:"data"`
      } `json:"video"`
   } `json:"relationships"`
   Type string `json:"type"`
}

func FetchInitiate(st *http.Cookie, market string) (*Initiate, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   fmt.Sprintf("default.beam-%v.prd.api.discomax.com", market),
         Path:   "/authentication/linkDevice/initiate",
      },
      Header: http.Header{},
   }
   req.AddCookie(st)
   req.Header.Set("x-device-info", device_info)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Attributes Initiate
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Fallback struct {
      Manifest struct {
         Url string // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

func (l *Login) FetchPlayReady(editId string) (*Playback, error) {
   return l.fetch_playback(editId, "playready")
}

func (l *Login) FetchWidevine(editId string) (*Playback, error) {
   return l.fetch_playback(editId, "widevine")
}

// Search queries the API and returns the root entity slice
func (l Login) Search(query string) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("page[items.size]", "10")
   queryParams.Set("contentFilter[query]", query)
   parsedUrl := &url.URL{
      Path:     "/cms/routes/search/result",
      RawQuery: queryParams.Encode(),
   }
   return l.fetch_entities(parsedUrl)
}

// GetMovie fetches the CMS data for a movie ID and returns the parsed entities
func (l Login) FetchMovie(movieRouteId string) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("page[items.size]", "1")
   parsedUrl := &url.URL{
      Path:     "/cms/routes/movie/" + movieRouteId,
      RawQuery: queryParams.Encode(),
   }
   return l.fetch_entities(parsedUrl)
}

func (l Login) FetchSeason(showId string, seasonNumber int) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("pf[show.id]", showId)
   queryParams.Set("pf[seasonNumber]", fmt.Sprint(seasonNumber))
   parsedUrl := &url.URL{
      Path:     "/cms/collections/generic-show-page-rail-episodes-tabbed-content",
      RawQuery: queryParams.Encode(),
   }
   return l.fetch_entities(parsedUrl)
}

func (l Login) fetch_entities(endpoint *url.URL) ([]*Entity, error) {
   // Scheme
   endpoint.Scheme = "https"
   // Host
   endpoint.Host = "default.prd.api.hbomax.com"
   // RawQuery
   queryParams := endpoint.Query()
   queryParams.Set("include", "default")
   endpoint.RawQuery = queryParams.Encode()
   req := http.Request{
      URL:    endpoint,
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
   }
   var result struct {
      Included []*Entity `json:"included"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Included, nil
}

func SearchResults(entities []*Entity) ([]*Entity, error) {
   entitiesMap := make(map[string]*Entity)
   for _, entity := range entities {
      entitiesMap[entity.Id] = entity
   }

   var searchResultsCollection *Entity
   for _, entity := range entities {
      if entity.Type == "collection" && entity.Attributes.Alias == "search-page-rail-results" {
         searchResultsCollection = entity
         break
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

func EpisodeResults(entities []*Entity) []*Entity {
   var episodes []*Entity
   for _, item := range entities {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         episodes = append(episodes, item)
      }
   }
   // Sort episodes by EpisodeNumber using the modern slices.SortFunc
   slices.SortFunc(episodes, func(entityA, entityB *Entity) int {
      return cmp.Compare(entityA.Attributes.EpisodeNumber, entityB.Attributes.EpisodeNumber)
   })
   return episodes
}

// String implements the fmt.Stringer interface to provide a clean visual output for the Entity.
func (e *Entity) String() string {
   data := &strings.Builder{}

   // 1. print episode number if material type is episode
   if e.Attributes.MaterialType == "EPISODE" {
      fmt.Fprintf(data, "Episode: %d\n", e.Attributes.EpisodeNumber)
   }

   // 2. print attributes name
   fmt.Fprintf(data, "Name: %s\n", e.Attributes.Name)

   // 3 & 4. print edit ID if type is video, otherwise print ID
   if e.Type == "video" {
      fmt.Fprintf(data, "Edit ID: %s\n", e.Relationships.Edit.Data.Id)
   } else {
      fmt.Fprintf(data, "ID: %s\n", e.Id)
   }

   // 5. print either show type or video type
   if e.Attributes.ShowType != "" {
      fmt.Fprintf(data, "Show Type: %s\n", e.Attributes.ShowType)
   } else if e.Attributes.VideoType != "" {
      fmt.Fprintf(data, "Video Type: %s\n", e.Attributes.VideoType)
   }

   return strings.TrimSpace(data.String())
}

// SL2000 max 1080p
// SL3000 max 2160p
func (p *Playback) FetchPlayReady(body []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.PlayReady.LicenseUrl, "text/xml",
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

///

func (p *Playback) FetchWidevine(body []byte) ([]byte, error) {
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

func (p *Playback) FetchDash() (*Dash, error) {
   resp, err := http.Get(
      strings.Replace(p.Fallback.Manifest.Url, "_fallback", "", 1),
   )
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

type Scheme struct {
   LicenseUrl string
}

var Markets = []string{
   "amer",
   "apac",
   "emea",
   "latam",
}

func FetchSt() (*http.Cookie, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "default.prd.api.hbomax.com", // Refactored
         Path:     "/token",
         RawQuery: "realm=bolt",
      },
      Header: http.Header{},
   }
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}

type Dash struct {
   Body []byte
   Url  *url.URL
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
