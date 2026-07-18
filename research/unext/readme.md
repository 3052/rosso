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

Osaka is blocked but Tokyo works - I was tricked because the Osaka error is weird

~~~
mitmproxy --mode upstream:http://isp.decodo.com:10001 --set upstream_auth=USERNAME:PASSWORD

emulator -avd Pixel_5 -http-proxy http://127.0.0.1:8080

adb root
adb shell "sysctl -w net.ipv6.conf.all.disable_ipv6=1"
~~~
