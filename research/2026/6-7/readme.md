# amazon

https://wikipedia.org/wiki/Shuhari

country: GB
name: United Kingdom
monetization: ADS

amazon.co.uk/gp/video/detail/0OZ8RZLMBWRZNF1DGCKHX07107

amazon.com/mytv

amazon.com/gp/video/settings/your-devices

amazon.co.uk/gp/video/detail/B07XZHJ25H

create Pixel 5, Android 11 device. install system certificate

adb shell input text PASSWORD

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

- https://apkmirror.com/apk/amazon-mobile-llc/prime-video-android-tv-android-tv
- https://play.google.com/store/apps/details?id=com.amazon.amazonvideo.livingroom

~~~
pip install frida-tools
~~~

download and extract server:

https://github.com/frida/frida/releases

for example:

~~~
frida-server-17.3.2-android-x86.xz 
~~~

install app, then push server:

~~~
$frida = 'frida-server-17.3.2-android-x86'
adb root
adb push $frida /data/app/frida-server
adb shell chmod +x /data/app/frida-server
adb shell /data/app/frida-server
~~~

then:

https://github.com/httptoolkit/frida-interception-and-unpinning

update `config.js`:

1. `CERT_PEM` from `C:\Users\Steven\.mitmproxy\mitmproxy-ca-cert.pem`
2. `PROXY_PORT` to `8080`

~~~
frida -U `
-l config.js `
-l android/android-certificate-unpinning.js `
-f com.amazon.avod.thirdpartyclient 
~~~
