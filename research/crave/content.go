package crave

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// extractMediaID parses the trailing ID from a URL (e.g. .../movie/goldeneye-38860 -> 38860)
func extractMediaID(rawURL string) (string, error) {
   u, err := url.Parse(rawURL)
   if err != nil {
      return "", err
   }
   parts := strings.Split(strings.TrimSuffix(u.Path, "/"), "-")
   if len(parts) == 0 {
      return "", fmt.Errorf("invalid url format")
   }
   return parts[len(parts)-1], nil
}

const (
   graphqlURL  = "https://rte-api.bellmedia.ca/graphql"
   playbackURL = "https://playback.rte-api.bellmedia.ca/contents/%s"
   manifestURL = "https://stream.video.9c9media.com/meta/content/%s/contentpackage/%d/destination/%d/platform/1"
)

// GetContentID queries the GraphQL API to translate a Media ID to a Content ID
func (t *TokenResponse) GetContentID(mediaID string) (string, error) {
   query := `query GetShowpage($sessionContext: SessionContext!, $ids: [String!]!) { medias(sessionContext: $sessionContext, ids: $ids) { firstContent { id } } }`
   
   payload := map[string]interface{}{
      "query": query,
      "variables": map[string]interface{}{
         "ids": []string{mediaID},
         "sessionContext": map[string]string{
            "userMaturity": "ADULT",
            "userLanguage": "EN",
         },
      },
   }

   body, _ := json.Marshal(payload)
   req, _ := http.NewRequest(http.MethodPost, graphqlURL, bytes.NewBuffer(body))
   
   // The GraphQL endpoint uses a base64 encoded JSON string that includes the access token
   authData := map[string]string{
      "platform":    "platform_web",
      "accessToken": t.AccessToken,
   }
   authBytes, _ := json.Marshal(authData)
   encodedAuth := base64.StdEncoding.EncodeToString(authBytes)

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Authorization", "Bearer "+encodedAuth)
   req.Header.Set("x-client-platform", "platform_web")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   var result struct {
      Data struct {
         Medias[]struct {
            FirstContent struct {
               ID string `json:"id"`
            } `json:"firstContent"`
         } `json:"medias"`
      } `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   if len(result.Data.Medias) == 0 || result.Data.Medias[0].FirstContent.ID == "" {
      return "", fmt.Errorf("content ID not found in GraphQL response")
   }

   return result.Data.Medias[0].FirstContent.ID, nil
}

// GetPlaybackDetails retrieves the ContentPackage ID and Destination ID
func (t *TokenResponse) GetPlaybackDetails(contentID string) (int, int, error) {
   targetURL := fmt.Sprintf(playbackURL, contentID)
   req, _ := http.NewRequest(http.MethodGet, targetURL, nil)
   req.Header.Set("Authorization", "Bearer "+ t.AccessToken)
   req.Header.Set("x-client-platform", "platform_jasper_web")
   req.Header.Set("x-playback-language", "EN")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return 0, 0, err
   }
   defer resp.Body.Close()

   var result struct {
      ContentPackage struct {
         ID            int `json:"id"`
         DestinationID int `json:"destinationId"`
      } `json:"contentPackage"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return 0, 0, err
   }

   if result.ContentPackage.ID == 0 {
      return 0, 0, fmt.Errorf("invalid content package ID received")
   }

   return result.ContentPackage.ID, result.ContentPackage.DestinationID, nil
}

// GetManifest retrieves the .mpd playback manifest URL from the 9c9media metadata API
func (t *TokenResponse) GetManifest(contentID string, contentPackageID, destinationID int) (string, error) {
   targetURL := fmt.Sprintf(manifestURL, contentID, contentPackageID, destinationID)
   
   req, _ := http.NewRequest(http.MethodGet, targetURL, nil)
   
   // Append requested query parameters
   q := req.URL.Query()
   q.Add("format", "mpd")
   q.Add("filter", "fe")
   q.Add("uhd", "false")
   q.Add("hd", "true")
   q.Add("mcv", "false")
   q.Add("mca", "false")
   q.Add("mta", "true")
   q.Add("stt", "true")
   req.URL.RawQuery = q.Encode()

   req.Header.Set("Authorization", "Bearer "+ t.AccessToken)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   var result struct {
      Playback string `json:"playback"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   if result.Playback == "" {
      return "", fmt.Errorf("playback URL missing in manifest response")
   }

   return result.Playback, nil
}

// GetManifestFromURL orchestrates the entire flow from a public URL to the manifest URL.
func (t *TokenResponse) GetManifestFromURL(publicURL string) (string, error) {
   mediaID, err := extractMediaID(publicURL)
   if err != nil {
      return "", fmt.Errorf("failed to extract media ID: %w", err)
   }
   contentID, err := t.GetContentID(mediaID)
   if err != nil {
      return "", fmt.Errorf("failed to get content ID: %w", err)
   }
   pkgID, destID, err := t.GetPlaybackDetails(contentID)
   if err != nil {
      return "", fmt.Errorf("failed to get playback details: %w", err)
   }
   return t.GetManifest(contentID, pkgID, destID)
}
