// main.go
package main

import (
   "fmt"
   "log"
)

const (
   // TODO: Replace with your actual valid Bearer access token
   AccessToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI2OTk2N2RhOWM5M2VlZjVkZjIwZjg3MTIiLCJzY29wZSI6ImFjY291bnQ6d3JpdGUgZGVmYXVsdCBtYXR1cml0eTphZHVsdCIsImlzcyI6Imh0dHBzOi8vYWNjb3VudC5iZWxsbWVkaWEuY2EiLCJjb250ZXh0Ijp7InByb2ZpbGVfaWQiOiI2OTk3MGVmYTczMTc2ZDJiMmU1M2E1YTMiLCJicmFuZF9pZHMiOlsiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTExIiwiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTE0IiwiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTE1Il19LCJleHAiOjE3NzU5NTgyNjAsImlhdCI6MTc3NTk0Mzg2MCwidmVyc2lvbiI6IlYyIiwianRpIjoiZTJjNThlMGEtNGFmNy00MzQ0LWFlNjgtNTRhZDU4M2ViZjY4IiwiYXV0aG9yaXRpZXMiOlsiUkVHVUxBUl9VU0VSIl0sImNsaWVudF9pZCI6ImNyYXZlLXdlYiJ9.zMEL5wKq5VuCg9URCYY8vvQH5k4Fmal6wJGyMlBFiDbWsXWxBtUw--TYVDzL0yLCPkTQ__q6UzXxD3JvOHQ4ZkE0JMSQjBL5Za-V4EnOLY3uNq2gB3VIDCCGo8YdjVpD0oRvtG4TEHD1OUhQ16YpA_FlVSSZYVt-4MA2iHqeE3HW99YRN7Bh0WJ1ndhwURsrJjx69uiJV6CT9W6h4ZDZYX1HE06DfemHOtdQxG-MwrgmJo3pcovOQPfBuF5lTUaLl0QRKi_PVFs6eU8bswkHZdamUDuuuKVK53RwMxj6i5y4XCVnGcTQOAFtr9Jwlle-jBPcp95I91NdAYia4hfAxw"
   ContentID   = "986962"
)

func main() {
   fmt.Printf("Fetching playback info for Content ID: %s...\n", ContentID)
   playbackInfo, err := GetPlaybackInfo(ContentID, AccessToken)
   if err != nil {
      log.Fatalf("Error getting playback info: %v", err)
   }
   if len(playbackInfo.AvailableContentPackages) == 0 {
      log.Fatalf("No content packages available for this content.")
   }
   // We typically use the first available package matching our criteria (e.g., English)
   pkg := playbackInfo.AvailableContentPackages[0]
   packageID := fmt.Sprintf("%d", pkg.ID)
   destID := fmt.Sprintf("%d", pkg.DestinationID)
   fmt.Printf("Found Package ID: %s, Destination ID: %s\n", packageID, destID)
   fmt.Println("Fetching Stream Metadata...")
   streamInfo, err := GetStreamMeta(ContentID, packageID, destID, AccessToken)
   if err != nil {
      log.Fatalf("Error getting stream meta: %v", err)
   }
   fmt.Println("\n=== SUCCESS ===")
   fmt.Printf("MPD Playback URL:\n%s\n", streamInfo.Playback)
}
