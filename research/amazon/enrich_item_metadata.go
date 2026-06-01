// get_enrich_item_metadata.go
package amazon

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

type PlaybackExperienceMetadata struct {
   PlaybackEnvelope string `json:"playbackEnvelope"`
}

type PlaybackAction struct {
   TitleID                    string                     `json:"titleID"`
   PlaybackExperienceMetadata PlaybackExperienceMetadata `json:"playbackExperienceMetadata"`
}

type EnrichmentData struct {
   PlaybackActions []PlaybackAction `json:"playbackActions"`
}

type EnrichItemMetadataResponse struct {
   Enrichments map[string]EnrichmentData `json:"enrichments"`
}

func GetEnrichItemMetadata(s *Session) error {
   u := &url.URL{
      Scheme: "https",
      Host:   "www.amazon.com",
      Path:   "/gp/video/api/enrichItemMetadata",
   }
   q := u.Query()
   q.Set("metadataToEnrich", `{"placement":"DETAIL_BTF","playback":true}`)
   q.Set("titleIDsToEnrich", fmt.Sprintf(`["%s"]`, s.VideoID))
   u.RawQuery = q.Encode()
   req, err := http.NewRequest("GET", u.String(), nil)
   if err != nil {
      return err
   }
   req.Header.Set("X-Requested-With", "XMLHttpRequest")
   resp, err := s.Client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return err
   }

   var result EnrichItemMetadataResponse
   if err := json.Unmarshal(body, &result); err != nil {
      return err
   }

   vidData, ok := result.Enrichments[s.VideoID]
   if !ok {
      return fmt.Errorf("video data not found in enrichments for ID: %s", s.VideoID)
   }

   if len(vidData.PlaybackActions) == 0 {
      return fmt.Errorf("playbackActions not found or empty")
   }

   firstAction := vidData.PlaybackActions[0]
   if firstAction.TitleID == "" {
      return fmt.Errorf("titleID missing from playback action")
   }
   if firstAction.PlaybackExperienceMetadata.PlaybackEnvelope == "" {
      return fmt.Errorf("playbackEnvelope missing from playback experience metadata")
   }

   s.TargetTitleID = firstAction.TitleID
   s.PlaybackEnvelope = firstAction.PlaybackExperienceMetadata.PlaybackEnvelope

   return nil
}
