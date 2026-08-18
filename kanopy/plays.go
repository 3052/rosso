package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
)

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Caption struct {
   Language string `json:"language"`
   Files    []File `json:"files"`
   Label    string `json:"label"`
}

type File struct {
   Type string
   Url  string
}

type Login struct {
   Jwt               string `json:"jwt"`
   VisitorId         string `json:"visitorId"`
   UserId            int    `json:"userId"`
   KanopyKidsEnabled bool   `json:"kanopyKidsEnabled"`
   WebshopId         int    `json:"webshopId"`
   WebshopCode       string `json:"webshopCode"`
   UserRole          string `json:"userRole"`
}

type Manifest struct {
   Url            string
   ManifestType   string `json:"manifestType"`
   DrmType        string `json:"drmType"`
   StorageService string `json:"storageService"`
   Cdn            string `json:"cdn"`
   DrmLicenseId   string `json:"drmLicenseID"`
}

type Membership struct {
   IdentityId         int    `json:"identityId"`
   DomainId           int    `json:"domainId"`
   UserId             int    `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type PlayResponse struct {
   Captions  []Caption
   Manifests []Manifest
   PlayId    string

   HttpCode     int    `json:"httpCode"`       // 2026-08-17: error status code from API
   ErrorMsgLong string `json:"error_msg_long"` // 2026-08-17: error description from API
}

func CreatePlay(loginData *Login, membershipData *Membership, videoData *Video) (*PlayResponse, error) {
   body, err := json.Marshal(PlayRequest{
      DomainId: membershipData.DomainId,
      UserId:   loginData.UserId,
      VideoId:  videoData.VideoId,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", "https://www.kanopy.com/kapi/plays", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+loginData.Jwt)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-version", "web/undefined/undefined/undefined")

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var play PlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&play); err != nil {
      return nil, err
   }

   // 2026-08-17: check error response from API
   if play.HttpCode != 0 {
      return nil, &play
   }

   return &play, nil
}

func (p *PlayResponse) Error() string {
   return fmt.Sprintf("kanopy: plays request failed: %d: %s",
      p.HttpCode, p.ErrorMsgLong)
}

type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHtml string `json:"descriptionHtml"`
   DurationSeconds int    `json:"durationSeconds"`
   Alias           string `json:"alias"`
}

// plays.go
