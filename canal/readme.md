# canal

https://wikipedia.org/wiki/Shuhari

1. canalplus.cz/en
2. subscribe now
3. czech republic VPN
4. canal+ komplet, subscribe
5. e-mail
   - just use real email since no trial
6. first name
7. last name
8. I agree to the general terms and conditions
9. payment
10. enable JavaScript
11. card number
12. expiry date
13. security code
14. name
15. pay
16. password
   - Your password must meet following criteria (min. 8 characters, capital
     letter, digit, special character
17. confirm password
18. next

## web

~~~py
from mitmproxy import http

data = '''
console.log('_0xb40f61', _0xb40f61);
console.log('_0xffbd34', _0xffbd34);
console.log('_0x44b887', _0x44b887);
console.log('_0x5bdf04', _0x5bdf04);
console.log('_0x5430bb', _0x5430bb);
console.log('_0x4ab337', _0x4ab337);
return'Client'''

def response(f: http.HTTPFlow) -> None:
   if f.request.path.startswith('/static/js/main.4c582264.js'):
      f.response.text = f.response.text.replace("return'Client", data)
~~~

## com.canal.android.canal

- https://apkmirror.com/apk/groupe-canal/mycanal-vos-programmes-en-live-ou-en-replay
- https://boutique.canalplus.com/plan-de-site
- https://play.google.com/store/apps/details?id=com.canal.android.canal

create Android 7 device. install system certificate. need residential proxy

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

then:

~~~xml
action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = http
data.scheme = https
data.host = www.mycanal.fr

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = tvchannels
data.host = com.canal.android.canal

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = http
data.scheme = https
data.host = www.canalplus.com

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = https
data.host = mycanal.onelink.me
data.pathPrefix = /1424707377

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = https
data.scheme = http
data.host = mycan.al

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.host = com.canal.android.canal
data.pathPrefix = /content
data.pathPrefix = /startapp
data.scheme = mycanaltvlauncher
~~~

## com.canalplus.canalplus

- https://canalplus.cz
- https://play.google.com/store/apps/details?id=com.canalplus.canalplus

first:

<https://apk.gold/download?file_id=3155967/canalplus-app>

~~~
The APK failed to install.<br/> Error: INSTALL_FAILED_MISSING_SPLIT: Missing
split for com.canalplus.canalplus
~~~

https://apkpure.com/canal-app/com.canalplus.canalplus/download/11.3

create Android 7 device

1. select your region, cesko
2. submit
3. english
4. sign in
5. email
6. password
   - adb shell input text PASSWORD
7. log in

~~~xml
com.canalplus.canalplus.xml
action.name = android.intent.action.VIEW
category.name = android.intent.category.DEFAULT
category.name = android.intent.category.BROWSABLE
data.scheme = @string/application_scheme
data.scheme = @string/additional_scheme
~~~

then:

~~~
resources\res\values\strings.xml
<string name="additional_scheme">m7cp</string>
<string name="application_scheme">https</string>
~~~

if you do this:

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

it opens the app, but it stays stuck on the splash screen. same for this:

~~~
adb shell am start -a android.intent.action.VIEW `
-d m7cp://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

also no requests are made
