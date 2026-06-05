# Netflix (Android TV)

## 4K

Android 12+, which means no intercept

## 1080p x86

~~~
com.netflix.ninja
~~~

2017:

https://apkmirror.com/apk/netflix-inc/netflix-android-tv

## 1080p x86 armeabi-v7a

https://play.google.com/store/apps/details?id=com.netflix.ninja

~~~
com.netflix.ninja
~~~

create Television 1080p

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

fails with Android 12:

~~~
Failure [INSTALL_FAILED_NO_MATCHING_ABIS: Failed to extract native libraries,
res=-113]
~~~

works with Android 13, but:

https://issuetracker.google.com/issues/331256113

## 1080p armeabi-v7a armeabi-v7a

only Android 6 or 12 are available, which means either the app wont run, or it
will run but we wont be able to intercept
