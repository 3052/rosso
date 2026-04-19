package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type Manifest struct {
   ManifestType string `json:"manifestType"`
   Url          string `json:"url"`
   DrmLicenseID string `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func CreatePlay(domainId int, userId int, videoId int, jwt string) (*PlayResponse, error) {
   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/plays")
   if err != nil {
      return nil, err
   }

   requestPayload := PlayRequest{
      DomainId: domainId,
      UserId:   userId,
      VideoId:  videoId,
   }

   bodyBytes, err := json.Marshal(requestPayload)
   if err != nil {
      return nil, err
   }

   requestHeaders := map[string]string{
      "content-type":  "application/json",
      "authorization": "Bearer " + jwt,
   }

   response, err := maya.Post(targetUrl, requestHeaders, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   responseBytes, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   var playResponse PlayResponse
   err = json.Unmarshal(responseBytes, &playResponse)
   if err != nil {
      return nil, err
   }

   return &playResponse, nil
}
