package paramount

var AppSecretProviders = []struct {
   url      string
   title    string
   id       string
   java     string
   uploaded string
}{
   {
      url:      "https://apkmirror.com/apk/cbs-interactive-inc/cbs",
      title:    "CBS",
      id:       "com.cbs.tve",
      uploaded: "December 17, 2025",
   },
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount",
      title:    "Paramount+",
      id:       "com.cbs.app",
      java:     "sources/com/cbs/app/config/UsaMobileAppConfigProvider.java",
      uploaded: "March 23, 2026",
   },
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-2",
      title:    "Paramount+ (Android TV)",
      id:       "com.cbs.ott",
      uploaded: "March 25, 2026",
   },
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-3",
      title:    "Paramount+",
      id:       "com.viacom.paramountplus",
      uploaded: "May 27, 2021",
   },
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-4",
      title:    "Paramount+",
      id:       "com.cbs.ca",
      uploaded: "March 23, 2026",
   },
   {
      url:      "https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv",
      title:    "Paramount+ (Android TV)",
      id:       "com.cbs.ca",
      java:     "sources/com/cbs/app/config/DefaultAppSecretProvider.java",
      uploaded: "March 24, 2026",
   },
}
