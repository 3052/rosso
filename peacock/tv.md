# tv

- https://apkmirror.com/apk/peacock-tv-llc/peacock-tv-android-tv
- https://play.google.com/store/apps/details?id=com.peacocktv.peacockandroid

create Pixel 5, APK lies you need at least API 31. install system certificate

~~~
emulator -avd Pixel_5 -http-proxy http://127.0.0.1:8080
~~~

then:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell monkey -p com.peacocktv.peacockandroid `
-c android.intent.category.LEANBACK_LAUNCHER 1
~~~

then:

~~~
adb shell pm clear com.peacocktv.peacockandroid
~~~
