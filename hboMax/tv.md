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
emulator -avd Pixel_5 -http-proxy http://127.0.0.1:8080 -no-snapshot-load
~~~

then:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell monkey -p com.wbd.stream -c android.intent.category.LEANBACK_LAUNCHER 1
~~~

Stop/Clear the app:

~~~
adb shell pm clear com.wbd.stream
~~~
