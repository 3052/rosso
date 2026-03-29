package paramount

var AppSecretProviders = []struct {
   uploaded string
   version  string
   url      string
   id       string
   java     string
}{
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv",
      id:       "com.cbs.ca",
      java:     "sources/com/cbs/app/config/DefaultAppSecretProvider.java",
      uploaded: "March 24, 2026",
      version:  "Paramount+ (Android TV) 16.8.0",
   },
   {
      uploaded: "March 23, 2026",
      id:       "com.cbs.app",
      java:     "sources/com/cbs/app/config/UsaMobileAppConfigProvider.java",
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount",
      version:  "Paramount+ 16.8.0",
   },
   {
      url:      "https://apkmirror.com/apk/cbs-interactive-inc/cbs",
      id:       "com.cbs.tve",
      version:  "CBS 15.6.0",
      uploaded: "December 17, 2025",
   },
   {
      url:      "https://apkmirror.com/apk/cbs-interactive-inc/cbs-android-tv",
      id:       "com.cbs.tve",
      version:  "CBS (Android TV) 15.6.0",
      uploaded: "December 17, 2025",
   },
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-2",
      uploaded: "March 25, 2026",
      id:       "com.cbs.ott",
      version:  "Paramount+ (Android TV) 16.8.0",
   },
   {
      id:       "com.cbs.ca",
      uploaded: "March 23, 2026",
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-4",
      version:  "Paramount+ 16.8.0",
   },
}
