# Paramount+

## how to get cmsAccountId

its in the HTML response body:

<https://paramountplus.com/shows/video/8PO2sBBr6lFb7J4nklXuzNZRhUR_V9dd>

## How to get secret\_key?

~~~
sources/com/cbs/app/androiddata/retrofit/util/RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

## how to get app secret?

us:

~~~
sources/com/cbs/app/config/UsaMobileAppConfigProvider.java
~~~

- https://apkmirror.com/apk/cbs-interactive-inc/paramount
- https://play.google.com/store/apps/details?id=com.cbs.app

international:

~~~
sources/com/cbs/app/config/DefaultAppSecretProvider.java
~~~

- https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv
- https://play.google.com/store/apps/details?id=com.cbs.ca

## paramount-4

international

https://apkmirror.com/apk/viacomcbs-streaming/paramount-4

## paramount-3

old

https://apkmirror.com/apk/viacomcbs-streaming/paramount-3

## paramount-2

android TV

https://apkmirror.com/apk/viacomcbs-streaming/paramount-2

## paramount

US

https://apkmirror.com/apk/viacomcbs-streaming/paramount

APK lies, you need at least Android 12 (level 31)

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
~/.android/avd/Pixel_XL.avd/emulator-user.ini
~~~

to:

~~~
window.x = 0
window.y = 0
~~~

install system certificate

## cbs

https://apkmirror.com/apk/cbs-interactive-inc/cbs

Create Android 7 device. Install system certificate
