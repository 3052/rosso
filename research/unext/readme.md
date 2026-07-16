# unext

https://wikipedia.org/wiki/Shuhari

## jp.unext.tv.player

https://play.google.com/store/apps/details?id=jp.unext.tv.player

## jp.unext.mediaplayer

- https://play.google.com/store/apps/details?id=jp.unext.mediaplayer
- https://u-next.en.uptodown.com

~~~
package: name='jp.unext.mediaplayer' versionCode='57100' versionName='5.71.0'
compileSdkVersion='37' compileSdkVersionCodename='17'
sdkVersion:'32'
  uses-feature: name='android.hardware.faketouch'
~~~

Create Pixel 5, API 32 device. install system certificate

Black Button:
Log in

~~~
adb shell pm clear jp.unext.mediaplayer
~~~

then:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell input text HELLO
~~~

1. Launch ProtonVPN, pick your exit country, wait until it shows "Connected"
2. In Android Studio → Device Manager → ⋮ next to your AVD → Cold Boot Now

~~~
adb shell settings put global http_proxy 10.0.2.2:8080
adb shell settings delete global http_proxy
~~~
