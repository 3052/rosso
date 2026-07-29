// der2pem.go
package main

import (
   "crypto/ecdsa"
   "crypto/rsa"
   "crypto/x509"
   "encoding/pem"
   "flag"
   "fmt"
   "os"
   "path/filepath"
)

// derKeyToPEM attempts to parse a DER private key in multiple formats
// (PKCS1 RSA, PKCS8, EC SEC1) and returns the PEM-encoded version.
func derKeyToPEM(der []byte) ([]byte, error) {
   // Try PKCS1 RSA first
   if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
      return pem.EncodeToMemory(&pem.Block{
         Type:  "RSA PRIVATE KEY",
         Bytes: x509.MarshalPKCS1PrivateKey(key),
      }), nil
   }

   // Try PKCS8 (can wrap RSA or EC keys)
   if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
      switch k := key.(type) {
      case *rsa.PrivateKey:
         return pem.EncodeToMemory(&pem.Block{
            Type:  "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(k),
         }), nil
      case *ecdsa.PrivateKey:
         keyBytes, err := x509.MarshalECPrivateKey(k)
         if err != nil {
            return nil, fmt.Errorf("failed to marshal EC key: %w", err)
         }
         return pem.EncodeToMemory(&pem.Block{
            Type:  "EC PRIVATE KEY",
            Bytes: keyBytes,
         }), nil
      default:
         return nil, fmt.Errorf("unsupported PKCS8 key type: %T", key)
      }
   }

   // Try EC SEC1
   if key, err := x509.ParseECPrivateKey(der); err == nil {
      keyBytes, err := x509.MarshalECPrivateKey(key)
      if err != nil {
         return nil, fmt.Errorf("failed to re-marshal EC key: %w", err)
      }
      return pem.EncodeToMemory(&pem.Block{
         Type:  "EC PRIVATE KEY",
         Bytes: keyBytes,
      }), nil
   }

   return nil, fmt.Errorf("could not parse key as PKCS1 RSA, PKCS8, or EC SEC1")
}

func main() {
   certPath := flag.String("cert", "", "Path to DER-encoded certificate file")
   keyPath := flag.String("key", "", "Path to DER-encoded private key file")
   host := flag.String("host", "", "Hostname")
   flag.Parse()

   if *certPath == "" || *keyPath == "" || *host == "" {
      fmt.Fprintf(os.Stderr, "Usage: %s -cert <cert.der> -key <key.der> -host <hostname>\n", os.Args[0])
      fmt.Fprintf(os.Stderr, "\nFlags:\n")
      flag.PrintDefaults()
      fmt.Fprintf(os.Stderr, "\nExample:\n")
      fmt.Fprintf(os.Stderr, "  %s -cert extracted_client.pem -key extracted_client.key -host play.clients.peacocktv.com\n", os.Args[0])
      os.Exit(1)
   }

   // --- Read & parse the DER certificate ---
   certDER, err := os.ReadFile(*certPath)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error reading cert DER: %v\n", err)
      os.Exit(1)
   }

   cert, err := x509.ParseCertificate(certDER)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error parsing certificate: %v\n", err)
      os.Exit(1)
   }

   certPEM := pem.EncodeToMemory(&pem.Block{
      Type:  "CERTIFICATE",
      Bytes: certDER,
   })

   fmt.Printf("Certificate loaded:\n")
   fmt.Printf("  Subject:   %s\n", cert.Subject.String())
   fmt.Printf("  Issuer:    %s\n", cert.Issuer.String())
   fmt.Printf("  Not After: %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))
   fmt.Println()

   // --- Read & parse the DER private key ---
   keyDER, err := os.ReadFile(*keyPath)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error reading key DER: %v\n", err)
      os.Exit(1)
   }

   keyPEM, err := derKeyToPEM(keyDER)
   if err != nil {
      fmt.Fprintf(os.Stderr, "Error converting key to PEM: %v\n", err)
      os.Exit(1)
   }
   fmt.Printf("Private key loaded and converted to PEM.\n\n")

   // --- Write combined PEM to current directory ---
   outPath := filepath.Join(".", *host+".pem")
   combined := append(certPEM, keyPEM...)

   if err := os.WriteFile(outPath, combined, 0600); err != nil {
      fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
      os.Exit(1)
   }
   fmt.Printf("✓ Wrote %s\n", outPath)

   fmt.Println("\nDone! Now launch mitmproxy:")
   fmt.Printf("  mitmproxy --set client_certs=.\n")
}
