// step_get_episodes_detail.go
package unext

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetEpisodeCodesViaDetail fetches all episode codes (ED...) for a given title
// code (SID...) using the Mad_VideoDetail operation.
func GetEpisodeCodesViaDetail(accessToken, titleCode string) ([]string, error) {
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "cc.unext.jp",
      Path:   "/",
   }

   body := map[string]any{
      "operationName": "Mad_VideoDetail",
      "variables":     map[string]string{"titleCode": titleCode},
      "query":         minVideoDetailQuery,
   }

   bodyJSON, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("get_episodes_detail: marshalling body: %w", err)
   }

   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewReader(bodyJSON))
   if err != nil {
      return nil, fmt.Errorf("get_episodes_detail: creating request: %w", err)
   }

   req.Header.Set("accept", "multipart/mixed;deferSpec=20220824, application/graphql-response+json, application/json")
   req.Header.Set("content-type", "application/json")
   req.Header.Set("apollo-require-preflight", "true")
   req.Header.Set("apollographql-client-name", "mad_for_mobile_jp.unext.mediaplayer")
   req.Header.Set("apollographql-client-version", "5.73.1")
   req.Header.Set("filmratingcode", "")
   req.Header.Set("u-device-id", "466d0fcd-79f5-3fb6-b580-cb34999f49dc")
   req.Header.Set("u-device-type", "920")
   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.73.1 sdk_gphone64_x86_64")
   req.Header.Set("x-apollo-operation-name", "Mad_VideoDetail")
   req.Header.Set("x-forwarded-for", "159.26.119.122")
   req.Header.Set("authorization", "Bearer "+accessToken)

   resp, err := clientDo(req)
   if err != nil {
      return nil, fmt.Errorf("get_episodes_detail: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("get_episodes_detail: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get_episodes_detail: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   var vdResp VideoDetailResponse
   if err := json.Unmarshal(respBody, &vdResp); err != nil {
      return nil, fmt.Errorf("get_episodes_detail: parsing response: %w (body starts with: %q)", err, string(respBody[:min(len(respBody), 50)]))
   }

   if len(vdResp.Errors) > 0 {
      return nil, fmt.Errorf("get_episodes_detail: GraphQL error: %s", vdResp.Errors[0].Message)
   }

   var codes []string
   for _, ep := range vdResp.Data.WebfrontTitleTitleEpisodes.Episodes {
      codes = append(codes, ep.ID)
   }

   return codes, nil
}

// VideoDetailResponse is the JSON envelope for the Mad_VideoDetail query.
// Only webfront_title_titleEpisodes is decoded; extra fields are ignored.
type VideoDetailResponse struct {
   Data struct {
      WebfrontTitleTitleEpisodes struct {
         Episodes []struct {
            ID string `json:"id"`
         } `json:"episodes"`
      } `json:"webfront_title_titleEpisodes"`
   } `json:"data"`
   Errors []GraphQLError `json:"errors"`
}
