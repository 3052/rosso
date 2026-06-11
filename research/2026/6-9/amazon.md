# amazon

- https://apkmirror.com/apk/amazon-mobile-llc/prime-video-android-tv-android-tv
- https://play.google.com/store/apps/details?id=com.amazon.amazonvideo.livingroom

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb push C:/Users/Steven/.mitmproxy/mitmproxy-ca-cert.pem /data/local/tmp/c8750f0d.0
adb shell
su
mkdir -p /data/misc/user/0/cacerts-added
mv /data/local/tmp/c8750f0d.0 /data/misc/user/0/cacerts-added/
chown system:system /data/misc/user/0/cacerts-added/c8750f0d.0
chmod 644 /data/misc/user/0/cacerts-added/c8750f0d.0
~~~

https://github.com/pwnlogs/cert-fixer

~~~
adb push Cert-Fixer.zip /data/local/tmp/
adb shell su -c 'magisk --install-module /data/local/tmp/Cert-Fixer.zip'
adb reboot
adb shell su -c 'ls /apex/com.android.conscrypt/cacerts | grep c8750f0d'
~~~

then:

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
adb push frida-server-17.11.0-android-x86 /data/local/tmp/frida-server
adb shell chmod +x /data/local/tmp/frida-server
adb shell su -c /data/local/tmp/frida-server
~~~

https://github.com/httptoolkit/frida-interception-and-unpinning

update `config.js`:

1. `CERT_PEM` from `C:\Users\Steven\.mitmproxy\mitmproxy-ca-cert.pem`
2. `PROXY_PORT` to `8080`

## spawn

~~~
frida -U `
-l config.js `
-l android/android-certificate-unpinning.js `
-f com.amazon.amazonvideo.livingroom

frida -U `
-l ./config.js `
-l ./native-connect-hook.js `
-l ./native-tls-hook.js `
-l ./android/android-proxy-override.js `
-l ./android/android-system-certificate-injection.js `
-l ./android/android-certificate-unpinning.js `
-l ./android/android-certificate-unpinning-fallback.js `
-l ./android/android-disable-root-detection.js `
-f com.amazon.amazonvideo.livingroom
~~~

## attach

~~~
frida -U `
-l ./config.js `
-l ./native-connect-hook.js `
-l ./native-tls-hook.js `
-l ./android/android-proxy-override.js `
-l ./android/android-system-certificate-injection.js `
-l ./android/android-certificate-unpinning.js `
-l ./android/android-certificate-unpinning-fallback.js `
-l ./android/android-disable-root-detection.js `
com.amazon.amazonvideo.livingroom
~~~
