package peacock

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// AtomResponse holds the raw JSON body returned by the Atom API.
type AtomResponse struct {
   Body []byte
}

// GetProviderVariant fetches node data for the given provider variant ID
// from the Atom adapter-calypso endpoint.
func (c *Client) GetProviderVariant(providerVariantID string) (*AtomResponse, error) {
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "atom.peacocktv.com",
      Path:   "/adapter-calypso/v3/query/node/provider_variant_id/" + providerVariantID,
   }
   q := url.Values{}
   q.Add("features", "upcoming")
   q.Add("represent", "(trailers[take=1],recs[take=1],collections(items),campaigns)")
   q.Add("contentSegments", "D2C,ESSENTIALS,Free")
   reqURL.RawQuery = q.Encode()

   req, err := c.newRequest(http.MethodGet, reqURL.String(), nil)
   if err != nil {
      return nil, err
   }

   req.Header.Add("accept", "*/*")
   req.Header.Add("accept-language", "en-US,en;q=0.9")
   req.Header.Add("origin", "https://tv.clients.peacocktv.com")
   req.Header.Add("referer", "https://tv.clients.peacocktv.com/")
   req.Header.Add("sec-fetch-dest", "empty")
   req.Header.Add("sec-fetch-mode", "cors")
   req.Header.Add("sec-fetch-site", "same-site")
   req.Header.Add("x-requested-with", "com.peacocktv.peacockandroid")
   req.Header.Add("x-skyott-ab-suggest", "olympicsSportsBoost:control")
   req.Header.Add("x-skyott-ab-atom", "BrandIcons:variant_brand_icons;VisionDeCampoTileImagery:variation__new_tile_image_;pzEntitlementOrder:pzeoFree1")
   req.Header.Add("x-skyott-ab-clip", "browseNTilesV1:variation_2;clipsearchmigration:variation_1;languageMigration:variation_1;newLabelsServiceV1:variation_1")
   req.Header.Add("x-skyott-subbouquetid", "0")

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return nil, fmt.Errorf("get provider variant: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get provider variant: bad status %d", resp.StatusCode)
   }

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("get provider variant: read body: %w", err)
   }

   return &AtomResponse{Body: body}, nil
}
