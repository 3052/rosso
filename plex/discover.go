package plex

import (
   "encoding/json"
   urlpkg "net/url"

   "41.neocities.org/maya"
)

type DiscoverAgeRating struct {
   Age    int    `json:"age"`
   Rating int    `json:"rating"`
   Type   string `json:"type"`
}

type DiscoverCommonSense struct {
   Key       string              `json:"key"`
   OneLiner  string              `json:"oneLiner"`
   AgeRating []DiscoverAgeRating `json:"AgeRating"`
}

type DiscoverImage struct {
   Alt  string `json:"alt"`
   Type string `json:"type"`
   Url  string `json:"url"`
}

type DiscoverGenre struct {
   Filter    string `json:"filter"`
   Id        string `json:"id"`
   Key       string `json:"key"`
   RatingKey string `json:"ratingKey"`
   Slug      string `json:"slug"`
   Tag       string `json:"tag"`
   Type      string `json:"type"`
   Context   string `json:"context"`
}

type DiscoverGuid struct {
   Id string `json:"id"`
}

type DiscoverRating struct {
   Image string  `json:"image"`
   Type  string  `json:"type"`
   Value float64 `json:"value"`
}

type DiscoverCountry struct {
   Tag string `json:"tag"`
}

type DiscoverRole struct {
   Id    string `json:"id"`
   Order int    `json:"order"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type DiscoverDirector struct {
   Id    string `json:"id"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type DiscoverProducer struct {
   Id    string `json:"id"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type DiscoverWriter struct {
   Id    string `json:"id"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type DiscoverStudio struct {
   Tag string `json:"tag"`
}

type DiscoverSummary struct {
   Size int    `json:"size"`
   Tag  string `json:"tag"`
   Type string `json:"type"`
}

type DiscoverStream struct {
   Bitrate      int    `json:"bitrate"`
   Codec        string `json:"codec"`
   Height       int    `json:"height"`
   StreamType   int    `json:"streamType"`
   Width        int    `json:"width"`
   Id           string `json:"id"`
   DisplayTitle string `json:"displayTitle"`
   Selected     bool   `json:"selected"`
   AudioTrackId string `json:"audioTrackId"`
   Language     string `json:"language"`
   Variant      string `json:"variant"`
   LanguageCode string `json:"languageCode"`
}

type DiscoverPart struct {
   Container string           `json:"container"`
   Id        string           `json:"id"`
   Key       string           `json:"key"`
   Stream    []DiscoverStream `json:"Stream"`
}

type DiscoverMedia struct {
   Bitrate               int            `json:"bitrate"`
   Container             string         `json:"container"`
   Protocol              string         `json:"protocol"`
   Url                   string         `json:"url"`
   Height                int            `json:"height"`
   Width                 int            `json:"width"`
   OptimizedForStreaming bool           `json:"optimizedForStreaming"`
   VideoCodec            string         `json:"videoCodec"`
   VideoResolution       string         `json:"videoResolution"`
   Part                  []DiscoverPart `json:"Part"`
}

type DiscoverMetadata struct {
   Art                   string                `json:"art"`
   Banner                string                `json:"banner"`
   Guid                  string                `json:"guid"`
   Key                   string                `json:"key"`
   PrimaryExtraKey       string                `json:"primaryExtraKey"`
   Rating                int                   `json:"rating"`
   RatingKey             string                `json:"ratingKey"`
   Studio                string                `json:"studio"`
   Summary               string                `json:"summary"`
   Tagline               string                `json:"tagline"`
   Type                  string                `json:"type"`
   AddedAt               int                   `json:"addedAt"`
   AudienceRating        float64               `json:"audienceRating"`
   AudienceRatingImage   string                `json:"audienceRatingImage"`
   AvailabilityId        string                `json:"availabilityId"`
   Budget                int                   `json:"budget"`
   ContentRating         string                `json:"contentRating"`
   Duration              int                   `json:"duration"`
   ImdbRatingCount       int                   `json:"imdbRatingCount"`
   OriginallyAvailableAt string                `json:"originallyAvailableAt"`
   PlayableKey           string                `json:"playableKey"`
   PublicPagesUrl        string                `json:"publicPagesURL"`
   RatingImage           string                `json:"ratingImage"`
   Revenue               int                   `json:"revenue"`
   Slug                  string                `json:"slug"`
   StreamingMediaId      string                `json:"streamingMediaId"`
   Thumb                 string                `json:"thumb"`
   Title                 string                `json:"title"`
   UserState             bool                  `json:"userState"`
   Year                  int                   `json:"year"`
   Media                 []DiscoverMedia       `json:"Media"`
   CommonSenseMedia      []DiscoverCommonSense `json:"CommonSenseMedia"`
   Image                 []DiscoverImage       `json:"Image"`
   Genre                 []DiscoverGenre       `json:"Genre"`
   Guids                 []DiscoverGuid        `json:"Guid"`
   Ratings               []DiscoverRating      `json:"Rating"`
   Country               []DiscoverCountry     `json:"Country"`
   Role                  []DiscoverRole        `json:"Role"`
   Director              []DiscoverDirector    `json:"Director"`
   Producer              []DiscoverProducer    `json:"Producer"`
   Writer                []DiscoverWriter      `json:"Writer"`
   Studios               []DiscoverStudio      `json:"Studio"`
   Summaries             []DiscoverSummary     `json:"Summary"`
}

type DiscoverMediaContainer struct {
   Offset              int                `json:"offset"`
   TotalSize           int                `json:"totalSize"`
   LibrarySectionId    string             `json:"librarySectionID"`
   LibrarySectionTitle string             `json:"librarySectionTitle"`
   Identifier          string             `json:"identifier"`
   Size                int                `json:"size"`
   Metadata            []DiscoverMetadata `json:"Metadata"`
}

type DiscoverResponse struct {
   MediaContainer DiscoverMediaContainer `json:"MediaContainer"`
}

func GetMatches(url string, authToken string) (*DiscoverResponse, error) {
   endpoint := &urlpkg.URL{
      Scheme: "https",
      Host:   "discover.provider.plex.tv",
      Path:   "/library/metadata/matches",
   }

   query := urlpkg.Values{}
   query.Set("url", url)
   query.Set("x-plex-token", authToken)
   endpoint.RawQuery = query.Encode()

   headers := map[string]string{
      "Accept": "application/json",
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var out DiscoverResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, err
   }

   return &out, nil
}
