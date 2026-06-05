# amazon

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
