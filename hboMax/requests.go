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

func StRequest() (*Cookie, error) {
   req, err := http.NewRequest(
      http.MethodGet,
      "https://default.prd.api.hbomax.com/token?realm=bolt",
      nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.Header.Set("x-disco-params", disco_params)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   for _, each := range resp.Cookies() {
      if each.Name == "st" {
         return &Cookie{Name: each.Name, Value: each.Value}, nil
      }
   }
   return nil, errors.New("named cookie not present")
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func LoginRequest(st *Cookie) (*Login, error) {
   req, err := http.NewRequest(http.MethodPost, "https://default.prd.api.hbomax.com/authentication/linkDevice/login", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("cookie", st.String())
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.Header.Set("x-disco-params", disco_params)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Errors APIErrors `json:"errors"`
      Data   struct {
         Attributes Login `json:"attributes"`
      } `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if len(result.Errors) > 0 {
      return nil, result.Errors
   }
   return &result.Data.Attributes, nil
}

func InitiateRequest(st *Cookie, market string) (*Initiate, error) {
   address := url.URL{
      Scheme: "https",
      Host:   fmt.Sprintf("default.any-%v.prd.api.discomax.com", market),
      Path:   "/authentication/linkDevice/initiate",
   }
   req, err := http.NewRequest(http.MethodPost, address.String(), nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("cookie", st.String())
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.Header.Set("x-disco-params", disco_params)

   resp, err := doReq(req)
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

func MovieRequest(token, movieId string) ([]*Entity, error) {
   values := url.Values{}
   values.Set("page[items.size]", "1")
   parsedUrl := &url.URL{
      Path:     "/cms/routes/movie/" + movieId,
      RawQuery: values.Encode(),
   }
   return entity_request(token, parsedUrl)
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

func PlayReadyRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "playready")
}

func WidevineRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "widevine")
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
      http.MethodPost,
      "https://default.prd.api.hbomax.com/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+token)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-disco-params", disco_params)
   req.Header.Set("x-disco-client", disco_client)
   req.Header.Set("x-device-info", device_info)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) > 0 {
      return nil, result.Errors
   }
   return &result, nil
}

func (*Playback) CachePath() string {
   return "rosso/hboMax/Playback"
}

func (p *Playback) GetManifest() (*url.URL, error) {
   manifest, err := url.Parse(p.Fallback.Manifest.Url)
   if err != nil {
      return nil, err
   }
   manifest.Path = strings.Replace(manifest.Path, "_fallback", "", 1)
   return manifest, nil
}

// SL2000 max 1080p
// SL3000 max 2160p
func (p *Playback) PlayReadyRequest(body []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, p.Drm.Schemes.PlayReady.LicenseUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "text/xml")

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (p *Playback) WidevineRequest(body []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, p.Drm.Schemes.Widevine.LicenseUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-protobuf")

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// requests.go
