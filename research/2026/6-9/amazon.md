# amazon

- https://apkmirror.com/apk/amazon-mobile-llc/prime-video-android-tv-android-tv
- https://play.google.com/store/apps/details?id=com.amazon.amazonvideo.livingroom

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

old:

~~~
adb push C:/Users/Steven/.mitmproxy/mitmproxy-ca-cert.pem /data/local/tmp/cacerts/c8750f0d.0
adb shell su -c 'cp /data/local/tmp/cacerts/* /system/etc/security/cacerts'
adb shell mkdir -p /data/local/tmp/cacerts
adb shell cp /system/etc/security/cacerts/* /data/local/tmp/cacerts
adb shell su -c 'mount -t tmpfs tmpfs /system/etc/security/cacerts'
adb shell su -c 'chcon u:object_r:system_file:s0 /system/etc/security/cacerts/*'
~~~

new:

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
