// step_get_episodes.go
package unext

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// GetEpisodeCodes fetches all episode codes (ED...) for a given title code (SID...)
// using the Mad_AllEpisodes operation.
func GetEpisodeCodes(titleCode string) ([]string, error) {
   body := map[string]any{
      "query": allEpisodesQuery,
      "variables": map[string]any{
         "titleCode": titleCode,
      },
   }
   bodyJSON, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("get_episodes: marshalling body: %w", err)
   }
   req, err := http.NewRequest("POST", "https://cc.unext.jp", bytes.NewReader(bodyJSON))
   if err != nil {
      return nil, fmt.Errorf("get_episodes: creating request: %w", err)
   }
   req.Header.Set("content-type", "application/json")
   resp, err := clientDo(req)
   if err != nil {
      return nil, fmt.Errorf("get_episodes: sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get_episodes: expected 200, got %d", resp.StatusCode)
   }

   var epResp EpisodesResponse
   if err := json.NewDecoder(resp.Body).Decode(&epResp); err != nil {
      return nil, fmt.Errorf("get_episodes: parsing response: %w", err)
   }

   if len(epResp.Errors) > 0 {
      return nil, fmt.Errorf("get_episodes: GraphQL error: %s", epResp.Errors[0].Message)
   }

   var codes []string
   for _, ep := range epResp.Data.WebfrontTitleTitleEpisodes.Episodes {
      codes = append(codes, ep.ID)
   }

   return codes, nil
}

// EpisodesResponse is the JSON envelope for the Mad_AllEpisodes query.
// Only webfront_title_titleEpisodes is decoded; extra fields are ignored.
type EpisodesResponse struct {
   Data struct {
      WebfrontTitleStage struct {
         TitleName string `json:"titleName"`
      } `json:"webfront_title_stage"`
      WebfrontTitleTitleEpisodes struct {
         Episodes []struct {
            ID string `json:"id"`
         } `json:"episodes"`
      } `json:"webfront_title_titleEpisodes"`
   } `json:"data"`
   Errors []GraphQLError `json:"errors"`
}
