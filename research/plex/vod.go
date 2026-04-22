package plex

import (
   "encoding/json"
   urlpkg "net/url"

   "41.neocities.org/maya"
)

type VodAgeRating struct {
   Age    int    `json:"age"`
   Rating int    `json:"rating"`
   Type   string `json:"type"`
}

type VodCommonSense struct {
   Key       string         `json:"key"`
   OneLiner  string         `json:"oneLiner"`
   AgeRating []VodAgeRating `json:"AgeRating"`
}

type VodImage struct {
   Alt  string `json:"alt"`
   Type string `json:"type"`
   Url  string `json:"url"`
}

type VodGenre struct {
   Filter    string `json:"filter"`
   Id        string `json:"id"`
   Key       string `json:"key"`
   RatingKey string `json:"ratingKey"`
   Slug      string `json:"slug"`
   Tag       string `json:"tag"`
   Type      string `json:"type"`
   Context   string `json:"context"`
}

type VodGuid struct {
   Id string `json:"id"`
}

type VodRating struct {
   Image string  `json:"image"`
   Type  string  `json:"type"`
   Value float64 `json:"value"`
}

type VodCountry struct {
   Tag string `json:"tag"`
}

type VodRole struct {
   Id    string `json:"id"`
   Order int    `json:"order"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type VodDirector struct {
   Id    string `json:"id"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type VodProducer struct {
   Id    string `json:"id"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type VodWriter struct {
   Id    string `json:"id"`
   Slug  string `json:"slug"`
   Tag   string `json:"tag"`
   Thumb string `json:"thumb"`
   Role  string `json:"role"`
   Key   string `json:"key"`
   Type  string `json:"type"`
}

type VodStudio struct {
   Tag string `json:"tag"`
}

type VodSummary struct {
   Size int    `json:"size"`
   Tag  string `json:"tag"`
   Type string `json:"type"`
}

type VodAd struct {
   Internal bool   `json:"internal"`
   Url      string `json:"url"`
   Type     string `json:"type"`
}

type VodStream struct {
   Bitrate      int    `json:"bitrate"`
   Codec        string `json:"codec"`
   Height       int    `json:"height"`
   StreamType   int    `json:"streamType"`
   Width        int    `json:"width"`
   DisplayTitle string `json:"displayTitle"`
   Selected     bool   `json:"selected"`
   AudioTrackId string `json:"audioTrackId"`
   Channels     int    `json:"channels"`
   Language     string `json:"language"`
   Variant      string `json:"variant"`
   LanguageCode string `json:"languageCode"`
   Id           string `json:"id"`
}

type VodPart struct {
   Container   string      `json:"container"`
   Id          string      `json:"id"`
   Key         string      `json:"key"`
   Certificate string      `json:"certificate"`
   Drm         string      `json:"drm"`
   Indexes     string      `json:"indexes"`
   License     string      `json:"license"`
   Stream      []VodStream `json:"Stream"`
}

type VodMedia struct {
   Bitrate               int       `json:"bitrate"`
   Container             string    `json:"container"`
   Drm                   bool      `json:"drm"`
   Height                int       `json:"height"`
   Protocol              string    `json:"protocol"`
   Width                 int       `json:"width"`
   Id                    string    `json:"id"`
   OptimizedForStreaming bool      `json:"optimizedForStreaming"`
   VideoCodec            string    `json:"videoCodec"`
   VideoResolution       string    `json:"videoResolution"`
   Part                  []VodPart `json:"Part"`
}

type VodMetadata struct {
   Attribution           string           `json:"attribution"`
   Art                   string           `json:"art"`
   Banner                string           `json:"banner"`
   Guid                  string           `json:"guid"`
   Key                   string           `json:"key"`
   PrimaryExtraKey       string           `json:"primaryExtraKey"`
   Rating                int              `json:"rating"`
   RatingKey             string           `json:"ratingKey"`
   Studio                string           `json:"studio"`
   Summary               string           `json:"summary"`
   Tagline               string           `json:"tagline"`
   Type                  string           `json:"type"`
   AddedAt               int              `json:"addedAt"`
   AudienceRating        float64          `json:"audienceRating"`
   AudienceRatingImage   string           `json:"audienceRatingImage"`
   AvailabilityId        string           `json:"availabilityId"`
   Budget                int              `json:"budget"`
   ContentRating         string           `json:"contentRating"`
   Duration              int              `json:"duration"`
   ExpiresAt             int              `json:"expiresAt"`
   ImdbRatingCount       int              `json:"imdbRatingCount"`
   OriginallyAvailableAt string           `json:"originallyAvailableAt"`
   PublicPagesUrl        string           `json:"publicPagesURL"`
   RatingImage           string           `json:"ratingImage"`
   Revenue               int              `json:"revenue"`
   Slug                  string           `json:"slug"`
   StreamingMediaId      string           `json:"streamingMediaId"`
   Thumb                 string           `json:"thumb"`
   Title                 string           `json:"title"`
   ViewCount             int              `json:"viewCount"`
   ViewOffset            int              `json:"viewOffset"`
   Year                  int              `json:"year"`
   Media                 []VodMedia       `json:"Media"`
   Ads                   []VodAd          `json:"Ad"`
   CommonSenseMedia      []VodCommonSense `json:"CommonSenseMedia"`
   Image                 []VodImage       `json:"Image"`
   Genre                 []VodGenre       `json:"Genre"`
   Guids                 []VodGuid        `json:"Guid"`
   Ratings               []VodRating      `json:"Rating"`
   Country               []VodCountry     `json:"Country"`
   Role                  []VodRole        `json:"Role"`
   Director              []VodDirector    `json:"Director"`
   Producer              []VodProducer    `json:"Producer"`
   Writer                []VodWriter      `json:"Writer"`
   Studios               []VodStudio      `json:"Studio"`
   Summaries             []VodSummary     `json:"Summary"`
}

type VodMediaContainer struct {
   LibrarySectionId    string        `json:"librarySectionID"`
   LibrarySectionTitle string        `json:"librarySectionTitle"`
   Offset              int           `json:"offset"`
   TotalSize           int           `json:"totalSize"`
   Identifier          string        `json:"identifier"`
   Size                int           `json:"size"`
   Metadata            []VodMetadata `json:"Metadata"`
}

type VodResponse struct {
   MediaContainer VodMediaContainer `json:"MediaContainer"`
}

func GetVod(Metadata *DiscoverMetadata, authToken string) (*VodResponse, error) {
   endpoint := &urlpkg.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   "/library/metadata/" + Metadata.RatingKey,
   }

   headers := map[string]string{
      "Accept":       "application/json",
      "X-Plex-Token": authToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var out VodResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, err
   }

   return &out, nil
}
