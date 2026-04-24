package plex

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type VodMetadata struct {
   MediaContainer VodContainer `json:"MediaContainer"`
}

type VodContainer struct {
   Metadata []MetadataItem `json:"Metadata"`
}

type MetadataItem struct {
   Guid  string     `json:"guid"`
   Title string     `json:"title"`
   Media []VodMedia `json:"Media"`
}

type VodMedia struct {
   Id       string    `json:"id"`
   Protocol string    `json:"protocol"`
   Part     []VodPart `json:"Part"`
}

type VodPart struct {
   Id  string `json:"id"`
   Key string `json:"key"`
}

func GetVodMetadata(match *MatchItem, anonymous *AnonymousUser) (*VodMetadata, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   "/library/metadata/" + match.RatingKey,
   }

   headers := map[string]string{
      "x-plex-token": anonymous.AuthToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var vod VodMetadata
   if err := json.NewDecoder(resp.Body).Decode(&vod); err != nil {
      return nil, err
   }

   return &vod, nil
}
