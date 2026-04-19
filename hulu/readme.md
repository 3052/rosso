# Hulu

## subscribe

1. hulu.com
2. choose Hulu
3. choose Hulu again (30 day free trial)
4. email
   - mail.tm (if you use this payment will fail)
   - mailsac.com
4. password
5. name
6. birthdate
7. gender
8. agree & continue
9. card number (NEW EACH TIME)
10. expiration
11. cvc
12. name on card
13. zip code
14. agree & subscribe

## tv

- https://apkmirror.com/apk/disney/hulu-android-tv
- https://play.google.com/store/apps/details?id=com.hulu.livingroomplus

## phone after login

- https://apkmirror.com/apk/disney/hulu-hulu
- https://play.google.com/store/apps/details?id=com.hulu.plus

create Android 7.1 device. install system certificate. do the login without
proxy, then kill app, start proxy and start app again, you will be able to
capture everything after the login

## phone login

- https://apkmirror.com/apk/disney/hulu-hulu
- https://play.google.com/store/apps/details?id=com.hulu.plus

create Android 7.1 device. install system certificate. enter credentials:

~~~
adb shell input text HELLO
~~~

after entering password, if you click LOG IN you get this:

> Hulu servers are not reachable. Check your internet connection and try again.

if we disable proxy? it works. next:

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
3. `DEBUG_MODE` to true

then:

~~~
frida -U `
-l config.js `
-l android/android-certificate-unpinning.js `
-l android/android-certificate-unpinning-fallback.js `
-f com.hulu.plus
~~~

if you get "Waiting for USB device to appear", you might need to reinstall
Python. now I get "Hulu has stopped". then:

~~~
adb logcat -T1 > hulu.txt
~~~

then:

~~~
[!] Matched class okhttp3.CertificatePinner but could not patch any methods
~~~

get `methodName`:

~~~
Thrown by okhttp3.internal.platform.android.AndroidCertificateChainCleaner->a
~~~

this worked a couple of times:

~~~diff
+++ b/android/android-certificate-unpinning.js
@@ -223,7 +223,7 @@ const PINNING_FIXES = {

     'okhttp3.CertificatePinner': [
         {
-            methodName: 'check',
+            methodName: 'a',
             overload: ['java.lang.String', 'java.util.List'],
             replacement: () => NO_OP
         },
~~~

but it seems to be a race condition or something, as it only works sometimes.
like it might fail the first time, but then if I restart the app it will work.
not sure.

https://github.com/httptoolkit/frida-interception-and-unpinning/issues/55
