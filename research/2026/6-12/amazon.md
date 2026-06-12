# amazon

you must use at least Android 14 (API 34) to get ARM support

but you also must use `x86_64` otherwise you get this with Frida:

Failed to attach: target terminated with signal 31

which means you must use Android 16 (API 36)

however this image only supports `x86_64` and arm64-v8a, and Amazon only supports
armeabi-v7a, which means we are fucked

- <https://gitlab.com/newbit/rootAVD/-/work_items/117>
- https://apkmirror.com/apk/amazon-mobile-llc/prime-video-android-tv-android-tv
- https://github.com/frida/frida-core/issues/1240
- https://github.com/httptoolkit/frida-interception-and-unpinning/issues/206
