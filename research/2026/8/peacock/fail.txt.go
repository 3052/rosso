package main

import (
	"bytes"
	"net/http"
	"net/url"
	"os"
)

func main() {
	client := &http.Client{}

	reqURL := &url.URL{
		Scheme: "https",
		Host:   "play.clients.peacocktv.com",
		Path:   "/drm/playready/acquirelicense",
	}

	q := url.Values{}
   
   //fail
	//q.Add("bt", "41-QfquYidBk3NF-Rph2MkjnnnFbMpea4I-I4DXihm5YbTH9kZ3aXQcPTiCJ76pC8Fl___gHb-y0rZk290xNPJ4aWqJwNP-9zR_c9q1UI6Wg6cOeYVxnMPhvIekzUO4dT1lINQ3CZbhV4MXAVRMmNx1xdY8R3papc2p0fp0GJcKlMU4-o97hrKj1jCBgiz_oK_nitU4eH8jSwijupOSGQGsArx-ZqyZfjbIirVmPwXR4m09QAnNLTtKc871heZ8PbU7Z2uGcV0pF2ohL2vKhnfn4EwFrT6rrBS8GfmU0j5lXN85yAVhJbYSRY8KrhFf-YySlu-jxdJo6uz6wn5jKhyGQzdt1R_20BpahDJfgeh8JN0R3fQfpP4th-HyRz9NCNaDWq8HtjJ6qkEaSVp32naB5jJw69LIYyQeSsptSfDmbSTxuAMvS2fO6axW2fItNS0TrqjMzbQ752aes17kcg0=")
   
   
	q.Add("bt", "42-Eceu3rYMJJXm2I55BW1hPZ71XhUwEGpmAajRsBD7nz9BCtP95cRdPilYKPFdprJTVYu9z0OhhXe1O4NS1VqoGbjIHhc7LNaRwhjMRJ0tQOeXVdq79JB1IkSpj3l8q7uTWkR_8JB0QSL1EZJ2hDJ4y3Qltr71TKAInF8WRTrQsFzupC_qTli4HySxoxPDiUlghtpKvFCH9XyZhWG44faG1Tyz5Hgx_iTQSeZncqY9Ez87aiRP3dAaSQg_Nnxv7q0fzy2YTmza8rY1g82kfNyhLeFrSTyRq7dqvcr-0QTzGrQenbnKE3AJN-TY-l1NmBcICZ7WOBxc9gDWqSFR4O5Gy4MB9B3CbPSBxekdQSoXMXmx6RtdubxPXPuNsmBCgMltra9H5l9wM5NDNqJE4AsnlD7JJKmnKAry7WV76uP2reopT3yAwszeOP0TdWXmgODN46Bu8uUq8Tmu")
   
   
	reqURL.RawQuery = q.Encode()

	bodyData := []byte(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><AcquireLicense xmlns="http://schemas.microsoft.com/DRM/2007/03/protocols"><challenge><Challenge xmlns="http://schemas.microsoft.com/DRM/2007/03/protocols/messages"><LA xmlns="http://schemas.microsoft.com/DRM/2007/03/protocols" Id="SignedData"><Version>1</Version><ContentHeader><WRMHEADER xmlns="http://schemas.microsoft.com/DRM/2007/03/PlayReadyHeader" version="4.0.0.0"><DATA><PROTECTINFO><KEYLEN>16</KEYLEN><ALGID>AESCTR</ALGID></PROTECTINFO><KID>deUaMAGpk/yfwEKTV4TfYw==</KID></DATA></WRMHEADER></ContentHeader><LicenseNonce>AQAAAAAAAAAAAAAAAAAAAA==</LicenseNonce><ClientTime>0</ClientTime><EncryptedData xmlns="http://www.w3.org/2001/04/xmlenc#" Type="http://www.w3.org/2001/04/xmlenc#Element"><EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes128-cbc"></EncryptionMethod><KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><EncryptedKey xmlns="http://www.w3.org/2001/04/xmlenc#"><EncryptionMethod Algorithm="http://schemas.microsoft.com/DRM/2007/03/protocols#ecc256"></EncryptionMethod><KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><KeyName>WMRMServer</KeyName></KeyInfo><CipherData><CipherValue>m78G2tmrWQXgVHHOFtUiLInCyqOfJiZ6wHRxKYhfvUQbzH+oTeEgo2dV2vMKb0fowNS93BUDbtKjRH36eh0+iIyH1qNb6HLCddPSNhZKatJrbNIYeQcZMurpQBftcq9Lu0nzf70404jG5+sAsh3o+w2A6a/6aLz80Q15x8kPNoU=</CipherValue></CipherData></EncryptedKey></KeyInfo><CipherData><CipherValue>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</CipherValue></CipherData></EncryptedData></LA><Signature><SignedInfo xmlns=""><Reference URI="#SignedData"><DigestValue>/GJxlD8kKD3AJD9qxiIdro5ikANSwu6Nqg3s2I67JeQ=</DigestValue></Reference></SignedInfo><SignatureValue>h1yQ/hlAvBkR09eevnlAk5APqlUHh+I9zYmAsZP9kluKRpKN/MU6v2avQbY5kTdHvb6YkxYT5bvjIjUc+zWmaA==</SignatureValue></Signature></Challenge></challenge></AcquireLicense></soap:Body></soap:Envelope>`)
	req, err := http.NewRequest("POST", reqURL.String(), bytes.NewBuffer(bodyData))

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if err := resp.Write(os.Stdout); err != nil {
		panic(err)
	}
}
