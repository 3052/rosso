package crave

import (
   "41.neocities.org/drm/widevine"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "testing"
)

func TestLicense(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/crave-final.json")
   if err != nil {
      t.Fatal(err)
   }
   var final_tokens TokenResponse
   err = json.Unmarshal(data, &final_tokens)
   if err != nil {
      t.Fatal(err)
   }
   log.SetFlags(log.Ltime)
   username, err := run("credential", "-h=api.nordvpn.com", "-k=username")
   if err != nil {
      t.Fatal(err)
   }
   password, err := run("credential", "-h=api.nordvpn.com")
   if err != nil {
      t.Fatal(err)
   }
   proxy := url.URL{
      Scheme: "https",
      User:   url.UserPassword(username, password),
      Host:   "ca1103.nordvpn.com:89",
   }
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         if req.Method == "" {
            req.Method = "GET"
         }
         log.Println(req.Method, req.URL)
         return &proxy, nil
      },
   }
   publicUrl := "https://www.crave.ca/en/movie/goldeneye-38860"
   // Magic happens here
   mediaId, err := extractMediaId(publicUrl)
   if err != nil {
      t.Fatal(err)
   }
   contentId, err := GetContentId(mediaId)
   if err != nil {
      t.Fatal(err)
   }
   pkgId, destId, err := GetPlaybackDetails(contentId)
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(cache + "/L3/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   pem_bytes, err := os.ReadFile(cache + "/L3/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := widevine.DecodePrivateKey(pem_bytes)
   if err != nil {
      t.Fatal(err)
   }
   //23:06:49.012 WARN : PSSH(WV):
   data, err = base64.StdEncoding.DecodeString(
      "CAESEJJ16PwoqrI1cAURGjhZbgIaCWJlbGxtZWRpYSITZmYtMDFmODdmOTEtMTQxODIxNw==",
   )
   if err != nil {
      t.Fatal(err)
   }
   // 1. Create the PsshData struct
   pssh, err := widevine.DecodePsshData(data)
   if err != nil {
      t.Fatal(err)
   }
   // 2. Build the License Request directly from the pssh struct
   req_bytes, err := pssh.EncodeLicenseRequest(client_id)
   if err != nil {
      t.Fatal(err)
   }
   // 3. Sign the request
   signed_bytes, err := widevine.EncodeSignedMessage(req_bytes, private_key)
   if err != nil {
      t.Fatalf("Failed to create signed request: %v", err)
   }
   // 4. Send to License Server
   payload := base64.StdEncoding.EncodeToString(signed_bytes)
   session := PlaybackSession{
      ContentId: contentId,
      ContentPackageId: pkgId,
      DestinationId: destId,
   }
   data, err = final_tokens.GetWidevineLicense(&session, payload)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", data)
}
