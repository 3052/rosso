package amazon

import (
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// LicenseResponse maps the JSON returned by the Amazon license endpoint.
type LicenseResponse struct {
   Widevine2License struct {
      License string `json:"license"`
   } `json:"widevine2License"`
   ErrorsByResource map[string]struct {
      ErrorCode string `json:"errorCode"`
      Message   string `json:"message"`
      Type      string `json:"type"`
   } `json:"errorsByResource"`
   Error struct {
      ErrorCode string `json:"errorCode"`
      Message   string `json:"message"`
      Type      string `json:"type"`
   } `json:"error"`
}

// GetWidevineLicense wraps the raw protobuf challenge, makes the request to Amazon,
// and unwraps the response returning the final Widevine license bytes.
func GetWidevineLicense(
   accessToken string,
   asin string,
   marketplaceID string,
   customerID string, // Obtained from the ManifestResponse SelectedEntitlement["grantedByCustomerId"]
   challenge []byte,
) ([]byte, error) {
   reqUrl := url.URL{
      Scheme: "https",
      Host:   "atv-ps.amazon.com",
      Path:   "/cdp/catalog/GetPlaybackResources",
   }

   gascEnabled := "false"

   // Python script logic: OS is Linux/unknown for SD, Windows/10.0 for higher.
   osName := "Windows"
   osVersion := "10.0"
   if DefaultPlaybackOptions.VideoQuality == "SD" {
      osName = "Linux"
      osVersion = "unknown"
   }

   q := reqUrl.Query()
   q.Set("asin", asin)
   q.Set("consumptionType", "Streaming")
   q.Set("desiredResources", "Widevine2License") // Requests Widevine instead of PlaybackUrls/PlayReady
   q.Set("deviceID", defaultDevice["device_serial"])
   q.Set("deviceTypeID", defaultDevice["device_type"])
   q.Set("firmware", "1")
   q.Set("gascEnabled", gascEnabled)
   q.Set("marketplaceID", marketplaceID)
   q.Set("resourceUsage", "ImmediateConsumption")
   q.Set("videoMaterialType", "Feature")
   q.Set("operatingSystemName", osName)
   q.Set("operatingSystemVersion", osVersion)
   q.Set("customerID", customerID)
   q.Set("deviceDrmOverride", "CENC")
   q.Set("deviceStreamingTechnologyOverride", "DASH")
   q.Set("deviceVideoQualityOverride", DefaultPlaybackOptions.VideoQuality)
   q.Set("deviceHdrFormatsOverride", DefaultPlaybackOptions.HDRFormat)

   reqUrl.RawQuery = q.Encode()

   // Widevine Challenge goes in the x-www-form-urlencoded body as a base64 string
   form := url.Values{}
   form.Set("widevine2Challenge", base64.StdEncoding.EncodeToString(challenge))
   form.Set("includeHdcpTestKeyInLicense", "true")

   req, err := http.NewRequest(http.MethodPost, reqUrl.String(), strings.NewReader(form.Encode()))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Accept", "application/json")
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("Authorization", "Bearer "+accessToken)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var licResp LicenseResponse
   if err := json.Unmarshal(bodyBytes, &licResp); err != nil {
      return nil, fmt.Errorf("failed to decode license response: %v\nBody: %s", err, string(bodyBytes))
   }

   // 1. Check for global errors (e.g. VPN block)
   if licResp.Error.ErrorCode != "" || licResp.Error.Type != "" {
      errCode := licResp.Error.ErrorCode
      if errCode == "" {
         errCode = licResp.Error.Type
      }
      if errCode == "PRS.NoRights.AnonymizerIP" {
         return nil, fmt.Errorf("amazon detected a Proxy/VPN and refused to return a license")
      }
      return nil, fmt.Errorf("amazon API error: %s [%s]", licResp.Error.Message, errCode)
   }

   // 2. Check for resource specific errors
   if resErr, ok := licResp.ErrorsByResource["widevine2License"]; ok {
      errCode := resErr.ErrorCode
      if errCode == "" {
         errCode = resErr.Type
      }
      if errCode == "PRS.NoRights.AnonymizerIP" {
         return nil, fmt.Errorf("amazon detected a Proxy/VPN and refused to return a license")
      }
      return nil, fmt.Errorf("amazon resource error: %s [%s]", resErr.Message, errCode)
   }

   // 3. Extract License
   if licResp.Widevine2License.License == "" {
      return nil, fmt.Errorf("license response is empty. Raw body: %s", string(bodyBytes))
   }

   // Decode the base64 wrapped license back to pure protobuf bytes
   licenseData, err := base64.StdEncoding.DecodeString(licResp.Widevine2License.License)
   if err != nil {
      return nil, fmt.Errorf("failed to decode base64 license: %v", err)
   }

   return licenseData, nil
}
