package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	playReady "41.neocities.org/diana/playReady"
)

// ---------------------------------------------------------------------------
// Config (mirrors the YAML)
// ---------------------------------------------------------------------------

type Profile struct {
	Platform  string
	Device    string
	ClientSDK string
	HMACKey   string
}

var profiles = map[string]Profile{
	"tv":      {"ANDROIDTV", "TV", "NBCU-ANDROID-v3", "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"},
	"android": {"ANDROID", "TABLET", "NBCU-ANDROID-v3", "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"},
	"web":     {"PC", "COMPUTER", "NBCU-WEB-v4", "FvT9VtwvhtSZvqnExMsvDDTEvBqR3HdsMcBFtWYV"},
}

const (
	clientTerritory   = "US"
	clientProvider    = "NBCU"
	clientProposition = "NBCUOTT"
	clientAuthScheme  = "MESSO"
	clientProfile     = "tv"
	clientDeviceID    = "PC"
	clientDrmDeviceID = "UNKNOWN"

	endpointPersonas = "https://persona.id.peacocktv.com/persona-store/personas"
	endpointTokens   = "https://play.ovp.peacocktv.com/auth/tokens"
	endpointNode     = "https://atom.peacocktv.com/adapter-calypso/v3/query/node"
	endpointVOD      = "https://play.ovp.peacocktv.com/video/playouts/vod"

	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"
)

var (
	prof    = profiles[clientProfile]
	hmacKey = []byte(prof.HMACKey)
)

// ---------------------------------------------------------------------------
// HTTP (honors HTTPS_PROXY env var automatically; MITM_INSECURE=1 disables verify)
// ---------------------------------------------------------------------------

var httpClient = func() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	if os.Getenv("MITM_INSECURE") != "" {
		t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // local proxy debugging only
	}
	return &http.Client{Transport: t}
}()

func doHTTP(method, rawURL string, body []byte, headers map[string]string, cookies string) (int, []byte, error) {
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, rawURL, rdr)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return resp.StatusCode, data, err
}

// ---------------------------------------------------------------------------
// Cookie file (Netscape format)
// ---------------------------------------------------------------------------

func loadCookies(path, domainFilter string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var pairs []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#HttpOnly_") {
			line = strings.TrimPrefix(line, "#HttpOnly_")
		} else if strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 7 || !strings.Contains(f[0], domainFilter) || f[6] == "" {
			continue
		}
		pairs = append(pairs, f[5]+"="+f[6])
	}
	return strings.Join(pairs, "; "), nil
}

// ---------------------------------------------------------------------------
// SkyOTT signing
// ---------------------------------------------------------------------------

func skyHeaders(extra map[string]string) map[string]string {
	h := map[string]string{
		"X-SkyOTT-Device":      prof.Device,
		"X-SkyOTT-Platform":    prof.Platform,
		"X-SkyOTT-Proposition": clientProposition,
		"X-SkyOTT-Provider":    clientProvider,
		"X-SkyOTT-Territory":   clientTerritory,
	}
	for k, v := range extra {
		h[k] = v
	}
	return h
}

func md5Hex(b []byte) string {
	sum := md5.Sum(b)
	return hex.EncodeToString(sum[:])
}

func md5Headers(headers map[string]string) string {
	var lines []string
	for k, v := range headers {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-skyott-") {
			lines = append(lines, lk+": "+v)
		}
	}
	sort.Strings(lines)
	text := strings.Join(lines, "\n")
	if len(lines) > 0 {
		text += "\n"
	}
	return md5Hex([]byte(text))
}

func sign(method, path string, headers map[string]string, body []byte) string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sdk := prof.ClientSDK
	msg := strings.Join([]string{
		strings.ToUpper(method), path, "", sdk, "1.0",
		md5Headers(headers), ts, md5Hex(body),
	}, "\n") + "\n"
	mac := hmac.New(sha1.New, hmacKey)
	mac.Write([]byte(msg))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf(`SkyOTT client="%s",signature="%s",timestamp="%s",version="1.0"`, sdk, sig, ts)
}

// ---------------------------------------------------------------------------
// Metadata: watch URL -> contentId / providerVariantId
// ---------------------------------------------------------------------------

var slugRe = regexp.MustCompile(`/(?:movies|tv|news|sports)/[a-z0-9_./-]+/(?:[a-f0-9-]{36}|\d+)`)

func getNodeMetadata(titleURL, cookies string) (title, contentID, variantID string, err error) {
	slug := slugRe.FindString(titleURL)
	if slug == "" {
		return "", "", "", fmt.Errorf("couldn't extract asset slug from %q", titleURL)
	}

	q := url.Values{"slug": {slug}, "represent": {"(items(items))"}}
	status, raw, err := doHTTP("GET", endpointNode+"?"+q.Encode(), nil, merge(skyHeaders(nil), map[string]string{
		"Accept":           "*",
		"Referer":          "https://www.peacocktv.com/watch/asset" + slug,
		"X-SkyOTT-Language": "en",
	}), cookies)
	if err != nil {
		return "", "", "", err
	}
	if status >= 400 {
		return "", "", "", fmt.Errorf("node query failed [%d]: %s", status, raw[:min(300, len(raw))])
	}

	var node struct {
		Attributes struct {
			Title             string `json:"title"`
			ProviderVariantID string `json:"providerVariantId"`
			Formats           map[string]struct {
				ContentID string `json:"contentId"`
			} `json:"formats"`
		} `json:"attributes"`
	}
	if err := json.Unmarshal(raw, &node); err != nil {
		return "", "", "", err
	}

	attrs := node.Attributes
	for _, pref := range []string{"UHD", "HD"} {
		if f, ok := attrs.Formats[pref]; ok {
			return attrs.Title, f.ContentID, attrs.ProviderVariantID, nil
		}
	}
	for _, f := range attrs.Formats {
		return attrs.Title, f.ContentID, attrs.ProviderVariantID, nil
	}
	return "", "", "", fmt.Errorf("no formats in node response")
}

// ---------------------------------------------------------------------------
// Personas + tokens
// ---------------------------------------------------------------------------

func getPersonaID(cookies string) string {
	h := skyHeaders(map[string]string{
		"Accept":              "application/vnd.persona.v1+json",
		"X-SkyOTT-TokenType":  clientAuthScheme,
	})
	status, raw, err := doHTTP("GET", endpointPersonas, nil, h, cookies)
	if err != nil || status >= 400 {
		return ""
	}
	var resp struct {
		Personas []struct {
			PersonaID string `json:"personaId"`
		} `json:"personas"`
	}
	if json.Unmarshal(raw, &resp) == nil && len(resp.Personas) > 0 {
		return resp.Personas[0].PersonaID
	}
	return ""
}

func getTokens(cookies string) (string, error) {
	skyH := skyHeaders(nil)

	auth := map[string]any{
		"authScheme":        clientAuthScheme,
		"provider":          clientProvider,
		"providerTerritory": clientTerritory,
		"proposition":       clientProposition,
	}
	if pid := getPersonaID(cookies); pid != "" {
		auth["personaId"] = pid
		fmt.Println("[+] personaId:", pid)
	}

	body, _ := json.Marshal(map[string]any{ // json.Marshal is compact; body MD5 is in the signature
		"auth": auth,
		"device": map[string]any{
			"type":        prof.Device,
			"platform":    prof.Platform,
			"id":          clientDeviceID,
			"drmDeviceId": clientDrmDeviceID,
		},
	})

	status, raw, err := doHTTP("POST", endpointTokens, body, merge(skyH, map[string]string{
		"Accept":           "application/vnd.tokens.v1+json",
		"Content-Type":     "application/vnd.tokens.v1+json",
		"X-Sky-Signature":  sign("POST", "/auth/tokens", skyH, body),
	}), cookies)
	if err != nil {
		return "", err
	}
	var resp struct {
		UserToken   string `json:"userToken"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", err
	}
	if status >= 400 || resp.UserToken == "" {
		return "", fmt.Errorf("token fetch failed [%d]: %s", status, resp.Description)
	}
	return resp.UserToken, nil
}

// ---------------------------------------------------------------------------
// Playout
// ---------------------------------------------------------------------------

func getPlayout(contentID, variantID, userToken, vcodec, colourSpace string) (dashURL, licenseURL string, err error) {
	skyH := map[string]string{
		"X-SkyOTT-Agent": strings.ToLower(strings.Join([]string{
			clientProposition, prof.Device, prof.Platform,
		}, ".")),
		"X-SkyOTT-PinOverride": "false",
		"X-SkyOTT-Provider":    clientProvider,
		"X-SkyOTT-Territory":   clientTerritory,
		"X-SkyOTT-UserToken":   userToken,
	}

	body, _ := json.Marshal(map[string]any{
		"device": map[string]any{
			"capabilities": []map[string]any{{
				"protection": "PLAYREADY",
				"container":  "ISOBMFF",
				"transport":  "DASH",
				"acodec":     "AAC",
				"vcodec":     vcodec,
			}},
			"maxVideoFormat":        map[bool]string{true: "UHD", false: "HD"}[vcodec == "H265"],
			"supportedColourSpaces": []string{colourSpace},
			"model":                 prof.Platform,
			"hdcpEnabled":           "true",
		},
		"client":                         map[string]any{"thirdParties": []string{"FREEWHEEL", "YOSPACE"}},
		"contentId":                      contentID,
		"providerVariantId":              variantID,
		"parentalControlPin":             "null",
		"personaParentalControlRating":   9,
	})

	status, raw, err := doHTTP("POST", endpointVOD, body, merge(skyH, map[string]string{
		"Accept":          "application/vnd.playvod.v1+json",
		"Content-Type":    "application/vnd.playvod.v1+json",
		"X-Sky-Signature": sign("POST", "/video/playouts/vod", skyH, body),
	}), "")
	if err != nil {
		return "", "", err
	}

	var resp struct {
		ErrorCode   string `json:"errorCode"`
		Description string `json:"description"`
		Protection  struct {
			LicenceAcquisitionURL string `json:"licenceAcquisitionUrl"`
		} `json:"protection"`
		Asset struct {
			Endpoints []struct {
				CDN string `json:"cdn"`
				URL string `json:"url"`
			} `json:"endpoints"`
		} `json:"asset"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", "", err
	}
	if status >= 400 || resp.ErrorCode != "" {
		return "", "", fmt.Errorf("playout error [%d]: %s (%s)", status, resp.Description, resp.ErrorCode)
	}

	for _, e := range resp.Asset.Endpoints {
		if strings.EqualFold(e.CDN, "FASTLY") {
			return e.URL, resp.Protection.LicenceAcquisitionURL, nil
		}
	}
	if len(resp.Asset.Endpoints) > 0 {
		return resp.Asset.Endpoints[0].URL, resp.Protection.LicenceAcquisitionURL, nil
	}
	return "", "", fmt.Errorf("no DASH endpoint in playout response")
}

// ---------------------------------------------------------------------------
// PSSH -> KID (via PlayReady Object in the MPD)
// ---------------------------------------------------------------------------

var playReadySystemID = []byte{0x9a, 0x04, 0xf0, 0x79, 0x98, 0x40, 0x42, 0x86,
	0xab, 0x92, 0xe6, 0x5b, 0xe0, 0x88, 0x5f, 0x95}

var psshRe = regexp.MustCompile(`<[a-zA-Z0-9:._-]*pssh[^>]*>([^<]+)</[a-zA-Z0-9:._-]*pssh>`)

func getPlayReadyKID(mpdURL string) ([]byte, error) {
	status, raw, err := doHTTP("GET", mpdURL, nil, nil, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("MPD fetch failed [%d]", status)
	}

	for _, m := range psshRe.FindAllStringSubmatch(string(raw), -1) {
		box, err := base64.StdEncoding.DecodeString(strings.TrimSpace(m[1]))
		if err != nil || len(box) < 32 {
			continue
		}
		if string(box[4:8]) != "pssh" || !bytes.Equal(box[12:28], playReadySystemID) {
			continue
		}
		// skip box header to the PlayReady Object (PRO) payload
		var pro []byte
		if box[8] == 0 { // version 0: [28:32]=data size, [32:]=data
			pro = box[32:]
		} else { // version 1: KID list first
			n := int(binary.BigEndian.Uint32(box[28:32]))
			off := 32 + 16*n + 4
			if off > len(box) {
				continue
			}
			pro = box[off:]
		}
		wrm, err := playReady.ParsePro(pro)
		if err != nil {
			continue
		}
		// WRMHEADER KID is GUID byte order; LicenseRequestBytes base64s it
		// verbatim into the challenge, so pass it through unmodified.
		return wrm.Data.Kid, nil
	}
	return nil, fmt.Errorf("no PlayReady PSSH found in MPD")
}

// ---------------------------------------------------------------------------
// PlayReady license request
// ---------------------------------------------------------------------------

func getPlayReadyLicense(licenseURL string, challenge []byte) ([]byte, error) {
	u, err := url.Parse(licenseURL)
	if err != nil {
		return nil, err
	}
	status, raw, err := doHTTP("POST", licenseURL, challenge, map[string]string{
		"Content-Type":    "text/xml; charset=utf-8",
		"X-Sky-Signature": sign("POST", u.Path, map[string]string{}, challenge),
	}, "")
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("license request failed [%d]: %s", status, raw[:min(300, len(raw))])
	}
	return raw, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func merge(base, extra map[string]string) map[string]string {
	for k, v := range extra {
		base[k] = v
	}
	return base
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	const (
		cookieFile = "cookies.txt"
		titleURL   = "https://www.peacocktv.com/watch/asset/movies/obsession/210b50e4-5ccd-3548-a1c5-0160585a4a67"
	)

	cookies, err := loadCookies(cookieFile, "peacocktv.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] loaded %d cookies\n", strings.Count(cookies, "="))

	// --- metadata ---
	title, contentID, variantID, err := getNodeMetadata(titleURL, cookies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] %s\n    contentId:         %s\n    providerVariantId: %s\n", title, contentID, variantID)

	// --- tokens ---
	userToken, err := getTokens(cookies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] userToken: %s...\n", userToken[:min(24, len(userToken))])

	// --- playout ---
	dashURL, licenseURL, err := getPlayout(contentID, variantID, userToken, "H264", "SDR")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] DASH manifest:", dashURL)
	fmt.Println("[+] License URL:  ", licenseURL)

	// --- KID from MPD ---
	kid, err := getPlayReadyKID(dashURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] KID (guid): %x\n", kid)

	// --- device ---
	chainBytes, err := os.ReadFile("bdevcert.dat")
	if err != nil {
		log.Fatal(err)
	}
	chain, err := playReady.ParseChain(chainBytes)
	if err != nil {
		log.Fatal(err)
	}
	sigBytes, err := os.ReadFile("zprivsig.dat")
	if err != nil {
		log.Fatal(err)
	}
	signingKey, err := playReady.ParseRawPrivateKey(sigBytes)
	if err != nil {
		log.Fatal(err)
	}
	encBytes, err := os.ReadFile("zprivencr.dat")
	if err != nil {
		log.Fatal(err)
	}
	encryptKey, err := playReady.ParseRawPrivateKey(encBytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(&chain.Certificates[0]) // manufacturer + security level

	// --- license ---
	challenge, err := chain.LicenseRequestBytes(signingKey, kid, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] Challenge built (%d bytes)\n", len(challenge))

	licenseResp, err := getPlayReadyLicense(licenseURL, challenge)
	if err != nil {
		log.Fatal(err)
	}
	license, err := playReady.ParseLicense(licenseResp)
	if err != nil {
		log.Fatal(err)
	}
	key, err := license.Decrypt(encryptKey)
	if err != nil {
		log.Fatal(err)
	}

	kidOut := append([]byte(nil), license.ContainerOuter.ContainerKeys.ContentKey.GuidKeyID...)
	playReady.UuidOrGuid(kidOut) // GUID -> UUID byte order for display
	fmt.Printf("[KEY] %x:%x\n", kidOut, key)
}
