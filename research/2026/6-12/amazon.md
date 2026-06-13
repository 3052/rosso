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

you also must use `x86_64` otherwise you get this with Frida?

Failed to attach: target terminated with signal 31

1. https://github.com/httptoolkit/frida-interception-and-unpinning/issues/206
2. https://issuetracker.google.com/issues/331256113
3. https://issuetracker.google.com/issues/522344738
4. https://github.com/frida/frida-core/issues/1240
5. <https://gitlab.com/newbit/rootAVD/-/work_items/117>
