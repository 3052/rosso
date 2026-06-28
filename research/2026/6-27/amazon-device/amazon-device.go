package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "io"
   "net/http"
   "os"
   "strings"
)

// InitConfig represents the outer JSON object assigned to ATVDeviceInitializationConfig
type InitConfig struct {
   DeviceProperties string `json:"DEVICE_PROPERTIES"`
   BlastOverride    string `json:"blastDeviceOverrideConfigJSON"`
}

// DeviceProps represents the unmarshaled DEVICE_PROPERTIES string
type DeviceProps struct {
   ClientName string `json:"CLIENT_NAME"`
}

// BlastOverride represents the unmarshaled blastDeviceOverrideConfigJSON string
type BlastOverride struct {
   Configs struct {
      ProductType map[string]interface{} `json:"PRODUCT_TYPE"`
   } `json:"configs"`
}

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

   // 1. Isolate the JSON block using strings.Cut
   _, after, ok := strings.Cut(htmlContent, `ATVDeviceInitializationConfig=`)
   if !ok {
      return fmt.Errorf("ATVDeviceInitializationConfig not found in the HTML payload")
   }

   // The JSON block ends right before the next variable assignment in the script
   jsonStr, _, ok := strings.Cut(after, `; injectedEnvValues`)
   if !ok {
      return fmt.Errorf("could not find the end boundary of the JSON payload")
   }

   // 2. Unmarshal the outer JSON payload
   var initConfig InitConfig
   if err := json.Unmarshal([]byte(jsonStr), &initConfig); err != nil {
      return fmt.Errorf("failed to unmarshal outer JSON: %w", err)
   }

   // 3. Unmarshal the nested DEVICE_PROPERTIES JSON
   var deviceProps DeviceProps
   if err := json.Unmarshal([]byte(initConfig.DeviceProperties), &deviceProps); err != nil {
      return fmt.Errorf("failed to unmarshal DEVICE_PROPERTIES: %w", err)
   }

   if deviceProps.ClientName == "" {
      return fmt.Errorf("manufacturer not found in the configuration payload")
   }

   // 4. Unmarshal the nested blastDeviceOverrideConfigJSON
   var blastOverride BlastOverride
   if err := json.Unmarshal([]byte(initConfig.BlastOverride), &blastOverride); err != nil {
      return fmt.Errorf("failed to unmarshal blastDeviceOverrideConfigJSON: %w", err)
   }

   // The model number is the dynamic key inside the PRODUCT_TYPE object
   var modelNumber string
   for key := range blastOverride.Configs.ProductType {
      modelNumber = key
      break // We just need the first/only key
   }

   if modelNumber == "" {
      return fmt.Errorf("model number not found in the configuration payload")
   }

   fmt.Printf("manufacturer name: %s\n", deviceProps.ClientName)
   fmt.Printf("model number: %s\n", modelNumber)

   return nil
}

func fetchAmazonConfig(dtid string) (string, error) {
   url := fmt.Sprintf("https://atv-ext.amazon.com/cdp/resources/app_host/index.html?deviceTypeID=%s", dtid)

   resp, err := http.Get(url)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("received HTTP status %d", resp.StatusCode)
   }

   var builder strings.Builder
   if _, err := io.Copy(&builder, resp.Body); err != nil {
      return "", err
   }

   return builder.String(), nil
}
