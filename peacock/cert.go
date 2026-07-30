// cert.go
package peacock

import (
   "crypto/tls"
   _ "embed"
   "net/http"
   "time"
)

//go:embed cert.pem
var certPEM []byte

//go:embed key.pem
var keyPEM []byte

// mtlsClient returns an *http.Client configured with the embedded mTLS
// certificate, ProxyFromEnvironment, and the given timeout.
func mtlsClient(timeout time.Duration) (*http.Client, error) {
   cert, err := tls.X509KeyPair(certPEM, keyPEM)
   if err != nil {
      return nil, err
   }

   return &http.Client{
      Timeout: timeout,
      Transport: &http.Transport{
         Proxy: http.ProxyFromEnvironment,
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion:   tls.VersionTLS12,
         },
      },
   }, nil
}
