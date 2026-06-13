# amazon

https://apkmirror.com/apk/amazon-mobile-llc/prime-video-android-tv-android-tv

even though its a TV app - you can use it with a phone device - just install as
normal and start:

~~~
adb shell monkey -p com.amazon.amazonvideo.livingroom `
-c android.intent.category.LEANBACK_LAUNCHER 1
~~~

or stop/clear:

~~~
adb shell pm clear com.amazon.amazonvideo.livingroom
~~~

create Pixel 5, Android 11 device. install system certificate

~~~
pip install frida-tools
~~~

https://github.com/frida/frida/releases

~~~
frida-server-17.3.2-android-x86.xz 
~~~

install app, then push server:

~~~
$frida = 'frida-server-17.12.0-android-x86'
adb root
adb push $frida /data/app/frida-server
adb shell chmod +x /data/app/frida-server
adb shell /data/app/frida-server
~~~

https://github.com/httptoolkit/frida-interception-and-unpinning

update `config.js`:

1. `CERT_PEM` from `C:\Users\Steven\.mitmproxy\mitmproxy-ca-cert.pem`
2. `PROXY_PORT` to `8080`

~~~
python run_frida.py
~~~

1. https://github.com/httptoolkit/frida-interception-and-unpinning/issues/207
2. https://issuetracker.google.com/issues/331256113
3. https://issuetracker.google.com/issues/522344738
4. https://github.com/frida/frida-core/issues/1240
5. <https://gitlab.com/newbit/rootAVD/-/work_items/117>
