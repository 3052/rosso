package main

import (
   "flag"
   "fmt"
   "io"
   "net/http"
   "os"
   "regexp"
   "strings"
)

func main() {
   var dtid string
   flag.StringVar(&dtid, "dtid", "", "The Amazon DeviceTypeID (e.g., A3GTP8TAF8V3YG)")
   flag.Parse()

   if dtid == "" {
      flag.Usage()
      os.Exit(1)
   }

   if err := printDeviceInfo(dtid); err != nil {
      fmt.Fprintf(os.Stderr, "Error: %v\n", err)
      os.Exit(1)
   }
}

func printDeviceInfo(dtid string) error {
   htmlContent, err := fetchAmazonConfig(dtid)
   if err != nil {
      return fmt.Errorf("failed to fetch data for DTID %s: %w", dtid, err)
   }

   regexManufacturer := regexp.MustCompile(`\\*"(?:CLIENT_NAME|CLIENT_PRIME_SIGNUP_NAME)\\*"\s*:\s*\\*"([^\\"]+)\\*"`)
   regexModel := regexp.MustCompile(`\\*"PRODUCT_TYPE\\*"\s*:\s*\{\s*\\*"([^\\"]+)\\*"`)
   fallbackModelRegex := regexp.MustCompile(`deviceIdsForLogging["\\]*\s*:\s*["\\]*.*?([A-Z0-9]{8,})`)

   var manufacturer string
   if m := regexManufacturer.FindStringSubmatch(htmlContent); len(m) > 1 {
      manufacturer = strings.ReplaceAll(m[1], " TV", "")
   }
   if manufacturer == "" {
      return fmt.Errorf("manufacturer not found in the configuration payload")
   }

   var modelNumber string
   if m := regexModel.FindStringSubmatch(htmlContent); len(m) > 1 {
      modelNumber = m[1]
   } else if fb := fallbackModelRegex.FindStringSubmatch(htmlContent); len(fb) > 1 {
      modelNumber = fb[1]
   }
   if modelNumber == "" {
      return fmt.Errorf("model number not found in the configuration payload")
   }

   fmt.Printf("manufacturer name: %s\n", manufacturer)
   fmt.Printf("model number: %s\n", modelNumber)

   return nil
}

func fetchAmazonConfig(dtid string) (string, error) {
   endpoints := []string{
      fmt.Sprintf("https://atv-ext.amazon.com/cdp/resources/app_host/index.html?deviceTypeID=%s", dtid),
      fmt.Sprintf("https://atv-ext.amazon.com/blast-app-hosting/html5/index.html?deviceTypeID=%s", dtid),
   }

   client := &http.Client{}

   for _, url := range endpoints {
      req, err := http.NewRequest("GET", url, nil)
      if err != nil {
         continue
      }

      req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

      resp, err := client.Do(req)
      if err != nil {
         continue
      }

      if resp.StatusCode == http.StatusOK {
         bodyBytes, err := io.ReadAll(resp.Body)
         resp.Body.Close()
         if err == nil {
            return string(bodyBytes), nil
         }
      } else {
         resp.Body.Close()
      }
   }

   return "", fmt.Errorf("could not find valid configuration on any known endpoint (HTTP 404/403)")
}
