// plays.go
package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type PlayRequest struct {
   DomainID int `json:"domainId"`
   UserID   int `json:"userId"`
   VideoID  int `json:"videoId"`
}

type StudioDrm struct {
   AuthXML      string `json:"authXml"`
   DRMLicenseID string `json:"drmLicenseId"`
}

type Manifest struct {
   ManifestType   string    `json:"manifestType"`
   URL            string    `json:"url"`
   DRMType        string    `json:"drmType"`
   StudioDRM      StudioDrm `json:"studioDrm"`
   StorageService string    `json:"storageService"`
   CDN            string    `json:"cdn"`
   DRMLicenseID   string    `json:"drmLicenseID"`
}

type PlayResponse struct {
   PlayID    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
}

func CreatePlay(req *PlayRequest, authorization string) (*PlayResponse, error) {
   targetURL, err := url.Parse("https://www.kanopy.com/kapi/plays")
   if err != nil {
      return nil, err
   }

   bodyBytes, err := json.Marshal(req)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type":  "application/json",
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
      "authorization": authorization,
   }

   resp, err := maya.Post(targetURL, headers, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlayResponse
   if err := json.Unmarshal(respBody, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
