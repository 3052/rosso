# Paramount+

## How to get secret\_key?

~~~
sources/com/cbs/app/androiddata/retrofit/util/RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

## com.cbs.app

https://apkmirror.com/apk/cbs-interactive-inc/paramount

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

## com.cbs.ca

https://apkmirror.com/apk/viacomcbs-streaming/paramount-4

## com.cbs.tve

https://apkmirror.com/apk/cbs-interactive-inc/cbs

Create Android 7 device. Install system certificate
