// FILE: vod_metadata.go
package plex

import (
   "encoding/json"
   "errors"
   "net/url"

   "41.neocities.org/maya"
)

type VodMetadata struct {
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
   Id      string `json:"id"`
   Key     string `json:"key"`
   License string `json:"license"`
}

func GetVodMetadata(match *MatchItem, anonymous *AnonymousUser) (*VodMetadata, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   match.Key,
   }

   headers := map[string]string{
      "x-plex-token": anonymous.AuthToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      MediaContainer VodMetadata `json:"MediaContainer"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result.MediaContainer, nil
}

func (vod *VodMetadata) GetDashMedia() (*VodMedia, error) {
   for _, item := range vod.Metadata {
      for _, media := range item.Media {
         if media.Protocol == "dash" {
            return &media, nil
         }
      }
   }
   return nil, errors.New("dash media not found")
}

func (media *VodMedia) GetMpdUrl(anonymous *AnonymousUser) (*url.URL, error) {
   if len(media.Part) == 0 {
      return nil, errors.New("no media parts found")
   }

   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   media.Part[0].Key,
   }

   query := url.Values{}
   query.Set("x-plex-token", anonymous.AuthToken)
   endpoint.RawQuery = query.Encode()

   return endpoint, nil
}
