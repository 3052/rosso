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

   initConfig, err := extractInitConfig(htmlContent)
   if err != nil {
      return err
   }

   // Unmarshal the nested DEVICE_PROPERTIES JSON to extract manufacturer
   var deviceProps struct {
      ClientName string `json:"CLIENT_NAME"`
   }
   if initConfig.DeviceProperties != "" {
      if err := json.Unmarshal([]byte(initConfig.DeviceProperties), &deviceProps); err != nil {
         return fmt.Errorf("failed to unmarshal DEVICE_PROPERTIES: %w", err)
      }
   }

   // Unmarshal the nested blastDeviceOverrideConfigJSON to extract the model
   var modelNumber string
   if initConfig.BlastOverride != "" {
      var blastOverride struct {
         Configs struct {
            ProductType map[string]interface{} `json:"PRODUCT_TYPE"`
         } `json:"configs"`
      }
      if err := json.Unmarshal([]byte(initConfig.BlastOverride), &blastOverride); err != nil {
         return fmt.Errorf("failed to unmarshal blastDeviceOverrideConfigJSON: %w", err)
      }

      // The model number is dynamically set as the key inside PRODUCT_TYPE
      for key := range blastOverride.Configs.ProductType {
         modelNumber = key
         break
      }
   }

   // Output ONLY the fields that were actually found
   if deviceProps.ClientName != "" {
      fmt.Printf("manufacturer name: %s\n", deviceProps.ClientName)
   }
   if modelNumber != "" {
      fmt.Printf("model number: %s\n", modelNumber)
   }

   return nil
}

// InitConfig represents the common properties we care about in the config
type InitConfig struct {
   DeviceProperties string `json:"DEVICE_PROPERTIES"`
   BlastOverride    string `json:"blastDeviceOverrideConfigJSON"`
}

// extractInitConfig checks known HTML injection patterns and extracts the config
func extractInitConfig(htmlContent string) (*InitConfig, error) {
   // Format 1: PS4 style -> `injectedEnvValues = {"ATVDeviceInitializationConfig": {...}}; </script>`
   if _, after, ok := strings.Cut(htmlContent, "injectedEnvValues = "); ok {
      if jsonStr, _, ok := strings.Cut(after, "; </script>"); ok {
         var wrapper struct {
            ATV InitConfig `json:"ATVDeviceInitializationConfig"`
         }
         // If this fails (e.g. invalid JSON), just fall through to the next format check
         if err := json.Unmarshal([]byte(jsonStr), &wrapper); err == nil && wrapper.ATV.DeviceProperties != "" {
            return &wrapper.ATV, nil
         }
      }
   }

   // Format 2: Hisense style -> `ATVDeviceInitializationConfig={...}; injectedEnvValues`
   if _, after, ok := strings.Cut(htmlContent, "ATVDeviceInitializationConfig="); ok {
      if jsonStr, _, ok := strings.Cut(after, "; injectedEnvValues"); ok {
         var initConfig InitConfig
         if err := json.Unmarshal([]byte(jsonStr), &initConfig); err == nil && initConfig.DeviceProperties != "" {
            return &initConfig, nil
         }
      }
   }

   return nil, fmt.Errorf("could not find valid configuration block in HTML payload")
}
