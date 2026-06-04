package amazon

import (
   "encoding/json"
   "errors"
   "strings"
)

// HydrationData strictly defines the JSON path down to the required fields.
// This completely avoids interface{} and type assertions.
type HydrationData struct {
   Init struct {
      Preparations struct {
         Atf struct {
            State struct {
               Action struct {
                  Atf map[string]TitleData `json:"atf"`
               } `json:"action"`
            } `json:"state"`
         } `json:"atf"`
      } `json:"preparations"`
   } `json:"init"`
}

type TitleData struct {
   PrimaryActions []ActionNode `json:"primaryActions"`
}

type ActionNode struct {
   Payload struct {
      Playback struct {
         IsTrailer        bool   `json:"isTrailer"`
         PlaybackEnvelope string `json:"playbackEnvelope"`
      } `json:"playback"`
   } `json:"payload"`
}

// ExtractFeaturePlaybackEnvelope parses the HTML string, extracts the embedded JSON
// using strings.Cut, and returns the playback envelope for the main feature.
func ExtractFeaturePlaybackEnvelope(htmlBody string) (string, error) {
   // 1. Locate the specific script tag identifier
   _, afterId, ok := strings.Cut(htmlBody, `id="dv-web-page-hydration-data"`)
   if !ok {
      return "", errors.New("hydration data script tag id not found")
   }

   // 2. Find the end of the opening <script> tag
   _, afterTag, ok := strings.Cut(afterId, ">")
   if !ok {
      return "", errors.New("malformed script tag")
   }

   // 3. Extract the JSON content up to the closing </script> tag
   jsonData, _, ok := strings.Cut(afterTag, "</script>")
   if !ok {
      return "", errors.New("closing script tag not found")
   }

   // 4. Unmarshal the extracted JSON string directly into our concrete struct
   var data HydrationData
   if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
      return "", errors.New("failed to parse hydration JSON: " + err.Error())
   }

   // 5. Traverse the strongly-typed data to find the main feature envelope.
   // As demonstrated in the source data, the main feature is within PrimaryActions.
   for _, titleData := range data.Init.Preparations.Atf.State.Action.Atf {
      for _, action := range titleData.PrimaryActions {
         pb := action.Payload.Playback
         if !pb.IsTrailer && pb.PlaybackEnvelope != "" {
            return pb.PlaybackEnvelope, nil
         }
      }
   }

   return "", errors.New("non-trailer playback envelope not found")
}
