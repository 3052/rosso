// step_get_episodes.go
package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// allEpisodesQuery is the full Mad_AllEpisodes query (required by server safelist).
const allEpisodesQuery = `query Mad_AllEpisodes(
  $titleCode: ID!
  $episodePage: Int
  $episodePageSize: Int
) {
  webfront_title_stage(id: $titleCode) {
    titleName
    nextUpdateDateTime
  }
  webfront_title_titleEpisodes(
    id: $titleCode
    page: $episodePage
    pageSize: $episodePageSize
  ) {
    pageInfo {
      page
      pageSize
      pages
      results
    }
    episodes {
      __typename
      ...ListVideoEpisodeInfo
    }
  }
}

fragment BlockVideoEpisodeInfo on Episode {
  id
  episodeTitleInfo {
    id
    name
  }
  episodeName
  thumbnail {
    standard
  }
  hasSubtitle
  hasDub
  duration
  displayNo
  interruption
  completeFlag
  publishStyleCode
  chromecastFlag
  productLineupCodeList
  hasPackRights
}

fragment ListVideoEpisodeInfo on Episode {
  __typename
  ...BlockVideoEpisodeInfo
  purchaseEpisodeLimitday
  endrollPosition
  downloadFlag
  chromecastFlag
  maxResolutionCode
  no
  saleTypeCode
  displayDurationText
  introduction
  nodCatchupPlanCode
  nodSpecialPlanCode
  movieTypeCode
  maxResolutionCode
  saleText
  isNew
  paymentBadgeList {
    name
    code
  }
  isPurchased
  purchaseEpisodeLimitday
  publicEndDate
  minimumPrice
  hasMultiplePrices
  episodeNotices
  playButtonName
}`

// GetEpisodeCodes fetches all episode codes (ED...) for a given title code (SID...)
// using the Mad_AllEpisodes operation.
func GetEpisodeCodes(client *http.Client, accessToken, titleCode string) ([]string, error) {
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "cc.unext.jp",
      Path:   "/",
   }

   variables := map[string]interface{}{
      "titleCode":       titleCode,
      "episodePage":     1,
      "episodePageSize": 1100,
   }

   variablesJSON, err := json.Marshal(variables)
   if err != nil {
      return nil, fmt.Errorf("get_episodes: marshalling variables: %w", err)
   }

   q := url.Values{}
   q.Add("operationName", "Mad_AllEpisodes")
   q.Add("variables", string(variablesJSON))
   q.Add("query", allEpisodesQuery)
   reqURL.RawQuery = q.Encode()

   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      return nil, fmt.Errorf("get_episodes: creating request: %w", err)
   }

   req.Header.Set("accept", "multipart/mixed;deferSpec=20220824, application/graphql-response+json, application/json")
   req.Header.Set("apollo-require-preflight", "true")
   req.Header.Set("apollographql-client-name", "mad_for_mobile_jp.unext.mediaplayer")
   req.Header.Set("apollographql-client-version", "5.73.1")
   req.Header.Set("filmratingcode", "")
   req.Header.Set("u-device-id", "466d0fcd-79f5-3fb6-b580-cb34999f49dc")
   req.Header.Set("u-device-type", "920")
   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.73.1 sdk_gphone64_x86_64")
   req.Header.Set("x-apollo-operation-name", "Mad_AllEpisodes")
   req.Header.Set("x-forwarded-for", "159.26.119.122")
   req.Header.Set("authorization", "Bearer "+accessToken)

   resp, err := clientDo(client, req)
   if err != nil {
      return nil, fmt.Errorf("get_episodes: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("get_episodes: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get_episodes: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   var epResp EpisodesResponse
   if err := json.Unmarshal(respBody, &epResp); err != nil {
      return nil, fmt.Errorf("get_episodes: parsing response: %w (body starts with: %q)", err, string(respBody[:min(len(respBody), 50)]))
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
