package rakuten

import (
   "encoding/json"
   "fmt"
   "net/url"
   "strconv"
   "strings"

   "41.neocities.org/maya"
)

type Movie struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

type ViewOptions struct {
   Private Private `json:"private"`
}

type Private struct {
   Streams []Stream `json:"streams"`
}

type Stream struct {
   AudioLanguages []Language `json:"audio_languages"`
}

func FetchMovie(movieId string, userClassification Classification, targetMarket Market) (*Movie, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + movieId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(userClassification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", targetMarket.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data Movie `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (targetMovie *Movie) String() string {
   return formatPlayableDetails(targetMovie.Id, targetMovie.Title, targetMovie.ViewOptions.Private.Streams)
}

func formatPlayableDetails(identifier string, title string, playbackStreams []Stream) string {
   seenLanguages := make(map[string]bool)
   var availableLanguages []string
   for _, currentStream := range playbackStreams {
      for _, audioLanguage := range currentStream.AudioLanguages {
         if !seenLanguages[audioLanguage.Id] {
            seenLanguages[audioLanguage.Id] = true
            availableLanguages = append(availableLanguages, audioLanguage.Id)
         }
      }
   }
   formattedAudio := strings.Join(availableLanguages, ", ")
   return fmt.Sprintf("%s (%s) - Audio: %s", title, identifier, formattedAudio)
}
