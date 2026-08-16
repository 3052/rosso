# plex

## android

https://play.google.com/store/apps/details?id=com.plexapp.android

1. create Android 6 device
2. install app
3. install user certificate

~~~
adb root
adb push frida-server-16.2.1-android-x86 /data/app/frida-server
adb shell chmod +x /data/app/frida-server
adb shell /data/app/frida-server

frida -U `
-l config.js `
-l android/android-certificate-unpinning.js `
-f com.plexapp.android
~~~

## movies only

~~~
GET https://vod.provider.plex.tv/library/metadata/movie:cruel-intentions HTTP/2.0
accept-encoding: identity
accept-language: en-US,en;q=0.5
accept: application/json
content-length: 0
content-type: application/json
origin: https://watch.plex.tv
referer: https://watch.plex.tv/
sec-fetch-dest: empty
sec-fetch-mode: cors
sec-fetch-site: same-site
te: trailers
x-plex-client-identifier: ff8a91f5-8f93-4dba-b61b-e0f286101d29
x-plex-language: en
x-plex-product: Plex Mediaverse
x-plex-provider-version: 6.5.0
~~~

## missing license

~~~
GET /library/metadata/matches?url=/movie/cruel-intentions&X-Plex-Token=aREUTWtbGNN8p_ChaGpv&X-Plex-Language=en-us HTTP/1.1
Host: discover.provider.plex.tv
Accept-Language: en-us
Accept: application/json
Connection: Keep-Alive
X-Plex-Advertising-DoNotTrack: 0
X-Plex-Advertising-Identifier: 7d57058e-e508-435d-98ff-8b6aa0cd9a9b
X-Plex-Client-Identifier: 429675a6d0ceebc9-com-plexapp-android
X-Plex-Client-Platform: Android
X-Plex-Device-Screen-Density: 420
X-Plex-Device-Screen-Resolution: 1920x1080 (Mobile)
X-Plex-Device-Vendor: unknown
X-Plex-Device: Android SDK built for x86
X-Plex-Features: external-media,indirect-media,hub-style-list
X-Plex-Model: generic_x86
X-Plex-Platform-Version: 6.0
X-Plex-Platform: Android
X-Plex-Product: Plex for Android (Mobile)
X-Plex-Provider-Version: 6.5.0
X-Plex-Provides: controller,sync-target
X-Plex-Session-Id: 2c509611-b10a-4dbe-bdf6-b17eebdddb3b
X-Plex-Version: 10.12.1.370
~~~

## includes license but needs metadata ID

~~~
GET https://vod.provider.plex.tv/library/metadata/5d7768b8594b2b001e693c52?checkFiles=1&includeReviews=1&includeExtras=1&asyncAugmentMetadata=1&X-Plex-Language=en-us HTTP/2.0
accept-language: en-us
accept: application/json
x-plex-advertising-donottrack: 0
x-plex-advertising-identifier: 621bba5a-646d-41d4-9e4a-202479d59127
x-plex-client-identifier: 1adb7ab9dee363ed-com-plexapp-android
x-plex-client-platform: Android
x-plex-device-screen-density: 420
x-plex-device-screen-resolution: 1920x1080 (Mobile)
x-plex-device-vendor: unknown
x-plex-device: Android SDK built for x86
x-plex-features: external-media,indirect-media,hub-style-list
x-plex-marketplace: googlePlay
x-plex-model: generic_x86
x-plex-platform-version: 6.0
x-plex-platform: Android
x-plex-product: Plex for Android (Mobile)
x-plex-provider-version: 6.5.0
x-plex-provides: controller,sync-target
x-plex-session-id: ffdd0a71-3725-4d2d-ba63-5989080912df
x-plex-token: V3KoARMyL631Hfjst8iP
x-plex-version: 10.12.1.370
~~~
