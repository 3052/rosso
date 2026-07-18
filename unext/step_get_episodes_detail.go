// step_get_episodes_detail.go
package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// videoDetailQuery is the full Mad_VideoDetail query (required by server safelist).
const videoDetailQuery = `query Mad_VideoDetail(
  $titleCode: ID!
  $episodeCode: ID
  $featureCode: ID
) {
  webfront_title_stage(id: $titleCode, featureCode: $featureCode) {
    __typename
    ...TitleStageCard
  }
  webfront_title_relatedTitles(id: $titleCode) {
    titles {
      __typename
      ...BlockVideoTitleInfo
    }
  }
  videoRelatedLives(titleCode: $titleCode, includePositionPlayLive: true) {
    liveList {
      __typename
      ...BlockLiveInfo
    }
  }
  webfront_title_titleEpisodes(id: $titleCode) {
    episodes {
      __typename
      ...ListVideoEpisodeInfo
    }
  }
  webfront_title_relatedBooks(id: $titleCode) {
    books {
      __typename
      ...BlockBookTitleInfo
    }
  }
  webfront_title_credits(id: $titleCode, page: 1, pageSize: 100) {
    hasEpisodeCredits
    titleCredits {
      __typename
      ...Cast
    }
  }
  webfront_title_recommendedTitles(id: $titleCode) {
    titles {
      __typename
      ...BlockVideoTitleInfo
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
}

fragment VideoExclusiveInfo on ExclusiveInfo {
  isOnlyOn
  typeCode
}

fragment SakuhinEventNotices on SakuhinEventNotice {
  id
  heading
  text
  url
  publicStartDatetime
  publicEndDatetime
}

fragment TitleStageCard on TitleStage {
  id
  titleName
  bookmarkStatus
  rate
  userRate
  productionYear
  country
  catchphrase
  attractions
  story
  check
  seriesCode
  seriesName
  publicStartDate
  publicEndDate
  publishStyleCode
  displayPublicEndDate
  sakuhinNotices
  publicMainEpisodeCount
  restrictedCode
  copyright
  bookmarkStatus
  thumbnail {
    secondary
    standard
  }
  mainGenreId
  mainGenreName
  isKids
  isNew
  hasDub
  hasSubtitle
  lastEpisode
  updateOfWeek
  nextUpdateDate
  nextUpdateDateTime
  nodSpecialFlag
  nodCatchupFlag
  hasMultiprice
  minimumPrice
  country
  productionYear
  paymentBadgeList {
    name
    code
  }
  nfreeBadge
  saleText
  missingAlertText
  episode(id: $episodeCode) {
    __typename
    ...ListVideoEpisodeInfo
  }
  keyEpisodes {
    current {
      __typename
      ...ListVideoEpisodeInfo
    }
    latest {
      __typename
      ...ListVideoEpisodeInfo
    }
  }
  comingSoonMainEpisodeCount
  feature {
    featureName
    titleComment
  }
  exclusive {
    __typename
    ...VideoExclusiveInfo
    isOnlyOn
    typeCode
  }
  isOriginal
  productLineupCodeList
  hasPackRights
  sakuhinEventNotices {
    __typename
    ...SakuhinEventNotices
  }
}

fragment BlockVideoTitleInfo on Title {
  id
  titleName
  isNew
  paymentBadgeList {
    code
    name
  }
  thumbnail {
    standard
    secondary
  }
  exclusive {
    __typename
    ...VideoExclusiveInfo
    isOnlyOn
    typeCode
  }
  isOriginal
  productLineupCodeList
  hasPackRights
  hasMultiprice
  minimumPrice
  isMaxService
}

fragment BlockLiveInfo on Live {
  id
  name
  saleTypeCode
  liveTypeCode
  parentLiveCode
  positionPlayLiveStartStatus
  deliveryStartDateTime
  deliveryEndDateTime
  programStartDateTime
  tickets {
    id
    name
    isSelling
    price
    saleEndDateTime
    saleStartDateTime
  }
  notices {
    thumbnail {
      standard
      secondary
    }
    name
    publicStartDateTime
  }
  isOnlyOn
  subContentList {
    typeCode
    publicStartDateTime
    publicEndDateTime
  }
  paymentBadgeList {
    name
    code
  }
  allowsTimeshiftedPlayback
  hasPackRights
  productLineupCodeList
}

fragment BlockBookProductInfo on Book {
  code
  name
  thumbnail {
    standard
  }
  credits {
    personCode
    penName
    penNameCode
    bookAuthorType
    unextPublishingDetail {
      thumbnail {
        standard
      }
      introduction
    }
  }
}

fragment BlockBookTitleInfo on BookSakuhin {
  sakuhinCode: code
  name
  freeBookNum
  isNew
  bookViewCode
  isUnextOriginal
  isChapter
  isBookPlusTicketAvailable
  isBookSakuhinTicketAvailable
  featurePieceCode
  bookLabel {
    code
    name
  }
  paymentBadgeList {
    code
    name
  }
  book {
    __typename
    mediaType {
      code
      name
    }
    ...BlockBookProductInfo
  }
  detail {
    catchSentence
  }
}

fragment Cast on Credit {
  castTypeName
  characterName
  personCode
  personName
  personNameCode
}`

// GetEpisodeCodesViaDetail fetches all episode codes (ED...) for a given title
// code (SID...) using the Mad_VideoDetail operation.
func GetEpisodeCodesViaDetail(client *http.Client, accessToken, titleCode string) ([]string, error) {
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "cc.unext.jp",
      Path:   "/",
   }

   variables := map[string]interface{}{
      "titleCode": titleCode,
   }

   variablesJSON, err := json.Marshal(variables)
   if err != nil {
      return nil, fmt.Errorf("get_episodes_detail: marshalling variables: %w", err)
   }

   q := url.Values{}
   q.Add("operationName", "Mad_VideoDetail")
   q.Add("variables", string(variablesJSON))
   q.Add("query", videoDetailQuery)
   reqURL.RawQuery = q.Encode()

   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      return nil, fmt.Errorf("get_episodes_detail: creating request: %w", err)
   }

   req.Header.Set("accept", "multipart/mixed;deferSpec=20220824, application/graphql-response+json, application/json")
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

   resp, err := clientDo(client, req)
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
