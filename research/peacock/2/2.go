package main

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   "encoding/hex"
   "net/http"
   "net/url"
   "os"
)

const bodyData = `
{
  "device": {
    "capabilities": [
      {
        "acodec": "AAC",
        "container": "ISOBMFF",
        "protection": "WIDEVINE",
        "transport": "DASH",
        "vcodec": "H264"
      },
      {
        "acodec": "AAC",
        "container": "ISOBMFF",
        "protection": "NONE",
        "transport": "DASH",
        "vcodec": "H264"
      }
    ],
    "maxVideoFormat": "HD",
    "supportedColourSpaces": [
      "SDR"
    ],
    "model": "ANDROIDTV",
    "hdcpEnabled": false
  },
  "client": {
    "thirdParties": [
      "FREEWHEEL",
      "MEDIATAILOR",
      "CONVIVA"
    ],
    "variantCapable": true
  },
  "contentId": "GMO_00000000158234_02_HDSDR",
  "providerVariantId": "1cba422b-3533-33a4-84af-d57cb97bbfa1",
  "parentalControlPin": "",
  "personaParentalControlRating": "9"
}
`

func main() {
   pemData, err := os.ReadFile("play.clients.peacocktv.com.pem")
   if err != nil {
      panic(err)
   }

   cert, err := tls.X509KeyPair(pemData, pemData)
   if err != nil {
      panic(err)
   }

   client := &http.Client{
      Transport: &http.Transport{
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion:   tls.VersionTLS12,
         },
      },
   }

   reqURL := &url.URL{
      Scheme: "https",
      Host:   "play.clients.peacocktv.com",
      Path:   "/video/playouts/vod",
   }

   hash := md5.Sum([]byte(bodyData))
   contentMD5 := hex.EncodeToString(hash[:])

   req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBufferString(bodyData))
   if err != nil {
      panic(err)
   }

   req.Header.Add("content-md5", contentMD5)
   req.Header.Add("accept", "application/vnd.playvod.v1+json")
   req.Header.Add("accept-language", "en-US,en;q=0.9")
   req.Header.Add("content-type", "application/vnd.playvod.v1+json")
   req.Header.Add("origin", "https://tv.clients.peacocktv.com")
   req.Header.Add("referer", "https://tv.clients.peacocktv.com/")
   req.Header.Add("sec-fetch-dest", "empty")
   req.Header.Add("sec-fetch-mode", "cors")
   req.Header.Add("sec-fetch-site", "same-site")
   req.Header.Add("x-requested-with", "com.peacocktv.peacockandroid")
   req.Header.Add("x-skyott-ab-atom", "BrandIcons:variant_brand_icons;VisionDeCampoTileImagery:variation__new_tile_image_;pzEntitlementOrder:pzeoFree1")
   req.Header.Add("x-skyott-usertoken", "159-vKTw7CSV/i7pOPr7vjKYFqu6oCxiaouo+FIgSa1U0n34ulq0gE1/SLs6iniJccmLBgwhdv/G4I0oBrZBL8OtEVWzV5YCrCDiHuGAyVpzKPmMgMD0I40sLu8zjAj1uMMu3Gfmt5zt+gqXOH9ibtpC3ApX5jIqz0yU+gtgPHs4DVcmTmG19rYA3xf2nXhs+cSCuj7mHGondSRAUwzW/MvKh8K4UuCcyZVWQwEeAOEDEvAFrcfZ7lAK8r6qp7TSsdmVsK9LunAg8A7sGkGdlqWGE5g3iwDwRGiwkyGmc+mxqj7WjID0jHJB8tMu3r0gm1AyNQxGI6bdXwthobHEVRmOkQ2GcQFG/YyLJioRgz1ti3UK1r7waTW4rLB1Eu+1R4wnw5aiC9ROdI3/lplALBn9yL++r37RG1H/Fgf49uxuA+wLUMWzvFmocb9ajIGokPntNo0LyC/cQju2qM5BQwYqYQ==")
   req.Header.Add("x-skyott-ab-suggest", "olympicsSportsBoost:control")
   req.Header.Add("x-skyott-activeterritory", "US")
   req.Header.Add("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Add("x-skyott-coppa", "false")
   req.Header.Add("x-skyott-device", "TV")
   req.Header.Add("x-skyott-language", "en-US")
   req.Header.Add("x-skyott-pinoverride", "false")
   req.Header.Add("x-skyott-platform", "ANDROIDTV")
   req.Header.Add("x-skyott-proposition", "NBCUOTT")
   req.Header.Add("x-skyott-provider", "NBCU")
   req.Header.Add("x-skyott-territory", "US")

   resp, err := client.Do(req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()

   if err := resp.Write(os.Stdout); err != nil {
      panic(err)
   }
}
