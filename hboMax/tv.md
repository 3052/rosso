# tv

## com.wbd.hbomax 

THIS IS OLD ONE

https://play.google.com/store/apps/details?id=com.wbd.hbomax 

## com.wbd.stream

THIS IS NEW ONE

https://apkmirror.com/apk/warnermedia-direct-llc/max-stream-hbo-tv-movies-android-tv

Even though it's a TV app, you can use it with a phone device - just install as normal

Create Pixel 5, Android 11 device. Install system certificate

~~~
adb shell monkey -p com.wbd.stream -c android.intent.category.LEANBACK_LAUNCHER 1
~~~

Stop/Clear the app:

~~~
adb shell pm clear com.wbd.stream
~~~

then:

~~~
emulator -avd Pixel_5 -http-proxy http://127.0.0.1:8080 -no-snapshot-load
~~~

then:

~~~
pip install frida-tools
~~~

https://github.com/frida/frida/releases

~~~
frida-server-17.3.2-android-x86.xz 
~~~

install app, then push server:

~~~
$frida = 'frida-server-17.16.2-android-x86'
adb root
adb push $frida /data/app/frida-server
adb shell chmod +x /data/app/frida-server
adb shell /data/app/frida-server
~~~

https://github.com/httptoolkit/frida-interception-and-unpinning

update `config.js`:

1. `CERT_PEM` from `C:\Users\Steven\.mitmproxy\mitmproxy-ca-cert.pem`
2. `PROXY_PORT` to `8080`
3. `DEBUG_MODE` to true

~~~
adb shell am start -n com.wbd.stream/com.wbd.beam.BeamActivity
frida -U --realm=emulated -n com.wbd.stream --eval "var m = Process.findModuleByName('libhbomax.so'); console.log('Path: ' + m.path);"

frida -U --realm=emulated `
-l ./config.js `
-l ./native-connect-hook.js `
-l ./native-tls-hook.js `
-l ./youi-ssl-bypass.js `
-n com.wbd.stream
~~~
