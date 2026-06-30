# molotov.tv

> I’m not really sure what the point of this video is, but I guess just be
> generous. be kind to people, because you never know how much they might need
> it, or how far it’ll go.
>
> [NakeyJakey](//youtube.com/watch?v=Cr0UYNKmrUs) (2018)

https://justwatch.com/fr/plateforme/molotov-tv

## subscribe

1. FRANCE VPN
2. molotov.tv
3. e-mail
   - mail.tm

## android

- https://play.google.com/store/apps/details?id=tv.molotov.app
- https://apkmirror.com/apk/molotov/molotov-tv-en-direct-et-en-replay

~~~
action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = https
data.scheme = http
data.host = www.molotov.tv
data.pathPrefix = /deeplink

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = molotov

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.host = app.molotov.tv
data.pathPattern = /.*
data.scheme = https
~~~

- <https://app.molotov.tv/p/program-details/program/VOD_314017>
- <https://molotov.tv/deeplink?channel_id=374&id=233268&type=program>

~~~
adb shell am start -a android.intent.action.VIEW `
-d molotov://deeplink?page=coupe-du-monde-fifa-2026
~~~

APK lies - you need at least API 31

~~~
adb shell input text HELLO
~~~
