# Android Studio

First download the package [1]. Start the program, and click More Actions,
AVD Manager, Create Virtual Device. On the Select Hardware screen, click Next.
On the System Image screen, click x86 Images. If the APK you are using
supports `x86`, then you can use versions down to API 23:

<https://docs.mitmproxy.org/stable/howto-install-system-trusted-ca-android>

Once you have chosen, click Download. Then click Next. On the Android Virtual
Device screen, click Finish. On the Your Virtual Devices screen, click Launch
this AVD in the emulator.

1. https://developer.android.com/studio#downloads

If you need to configure a proxy, in the emulator click More. On the
Extended Controls screen, click Settings, Proxy. Uncheck Use Android Studio
HTTP proxy settings. Click Manual proxy configuration. Enter Host name and
Port number as determined by the proxy program you are using. Click Apply,
and you should see Proxy status Success.
