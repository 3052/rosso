# amazon

## test 1

with no proxy everything is fine

## test 2

~~~
adb shell monkey -p com.topjohnwu.magisk -c android.intent.category.LAUNCHER 1
~~~

do you want to proceed and reboot? OK

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

set proxy, start app

Error code: 0.60

## test 3

start app, set proxy

error code: 9345

## test 4

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

~~~
frida -U `
-l ./android/android-certificate-unpinning-fallback.js `
-l ./android/android-certificate-unpinning.js `
-l ./android/android-disable-root-detection.js `
-l ./android/android-proxy-override.js `
-l ./android/android-system-certificate-injection.js `
-l ./config.js `
-l ./native-connect-hook.js `
-l ./native-tls-hook.js `
-f com.amazon.amazonvideo.livingroom
~~~

Failed to spawn: unable to find a front-door activity

## test 5

~~~
> python run_frida.py
[*] Combining scripts...
[*] Spawning com.amazon.amazonvideo.livingroom via com.amazon.ignition.IgnitionActivity...
[*] Attaching to PID 3219...
[-] An unexpected error occurred: target terminated with signal 31
~~~

## test 6

~~~
frida -U `
-l ./android/android-certificate-unpinning-fallback.js `
-l ./android/android-certificate-unpinning.js `
-l ./android/android-disable-root-detection.js `
-l ./android/android-proxy-override.js `
-l ./android/android-system-certificate-injection.js `
-l ./config.js `
-l ./native-connect-hook.js `
-l ./native-tls-hook.js `
com.amazon.amazonvideo.livingroom
~~~

Failed to attach: target terminated with signal 31

- https://github.com/frida/frida-core/issues/1240
- https://github.com/httptoolkit/frida-interception-and-unpinning/issues/206

## test 7

~~~
frida -U com.amazon.amazonvideo.livingroom
~~~

Failed to attach: target terminated with signal 31

## test 8

armeabi-v7a

~~~
> frida-ps -U -a
 PID  Name               Identifier
----  -----------------  -----------------------------
1940  Google Play Store  com.android.vending
2831  Magisk             com.topjohnwu.magisk
3761  Peacock TV         com.peacocktv.peacockandroid
2174  YouTube            com.google.android.youtube.tv

> frida -U 3761
Failed to attach: target terminated with signal 31
~~~

## test 9

noarch

~~~
frida -U kr.co.captv.pooq.tv
~~~

Failed to attach: target terminated with signal 31

## test 10

Android 9

~~~
frida -U kr.co.captv.pooq.tv
~~~

pass
