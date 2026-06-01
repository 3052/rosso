// get_enrich_item_metadata.go
package amazon

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

func GetEnrichItemMetadata(s *Session) error {
   reqUrl := fmt.Sprintf("https://www.amazon.com/gp/video/api/enrichItemMetadata?metadataToEnrich=%%7B%%22placement%%22%%3A%%22DETAIL_BTF%%22%%2C%%22playback%%22%%3Atrue%%7D&titleIDsToEnrich=%%5B%%22%s%%22%%5D", s.VideoID)

   req, err := http.NewRequest("GET", reqUrl, nil)
   if err != nil {
      return err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
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

   var result map[string]interface{}
   if err := json.Unmarshal(body, &result); err != nil {
      return err
   }

   enrichments, ok := result["enrichments"].(map[string]interface{})
   if !ok {
      return fmt.Errorf("enrichments not found in response")
   }

   vidData, ok := enrichments[s.VideoID].(map[string]interface{})
   if !ok {
      return fmt.Errorf("video data not found in enrichments")
   }

   playbackActions, ok := vidData["playbackActions"].([]interface{})
   if !ok || len(playbackActions) == 0 {
      return fmt.Errorf("playbackActions not found or empty")
   }

   firstAction := playbackActions[0].(map[string]interface{})
   titleID, ok := firstAction["titleID"].(string)
   if !ok {
      return fmt.Errorf("titleID missing from playback action")
   }
   s.TargetTitleID = titleID

   pem, ok := firstAction["playbackExperienceMetadata"].(map[string]interface{})
   if !ok {
      return fmt.Errorf("playbackExperienceMetadata missing")
   }

   envelope, ok := pem["playbackEnvelope"].(string)
   if !ok {
      return fmt.Errorf("playbackEnvelope missing")
   }
   s.PlaybackEnvelope = envelope

   return nil
}
