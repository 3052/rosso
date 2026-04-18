# peacock

## account

1. https://peacocktv.com
2. get started
3. PREMIUM
   1. monthly
   2. choose
4. Email
5. Password
6. Re-enter Password
7. First Name
8. Last Name
9. Gender
10. Birth Year
11. Zip Code
12. CREATE ACCOUNT
13. debit card
14. first name
15. last name
16. address
17. city
18. state
19. zip
20. card number
21. expiry date
22. CVC
23. by checking the box, you agree to pay
24. SUBSCRIBE
25. by checking the box, you agree to pay (again)
26. subscribe (again)

## how to get SkyOTT key?

use sky.js

## android

https://play.google.com/store/apps/details?id=com.peacocktv.peacockandroid

If you start the app and Sign In, this request:

~~~
POST https://rango.id.peacocktv.com/signin/service/international HTTP/2.0
content-type: application/x-www-form-urlencoded
x-skyott-device: MOBILE
x-skyott-proposition: NBCUOTT
x-skyott-provider: NBCU
x-skyott-territory: US

userIdentifier=MY_EMAIL&password=MY_PASSWORD
~~~

will fail:

~~~
HTTP/2.0 429
~~~

You can fix this problem by removing this request header before starting the
app:

~~~
set modify_headers '/~u signin.service.international/x-skyott-device/'
~~~

Header needs to be removed from that request only, since other requests need the
header.

## tv

~~~
$env:path = 'C:\windows\system32'
.\rootAVD.bat system-images\android-29\android-tv\x86\ramdisk.img

adb shell mkdir -p /data/local/tmp/cacerts
adb push C:/Users/Steven/.mitmproxy/mitmproxy-ca-cert.pem /data/local/tmp/cacerts/c8750f0d.0
adb shell cp /system/etc/security/cacerts/* /data/local/tmp/cacerts

adb shell su -c 'mount -t tmpfs tmpfs /system/etc/security/cacerts'
adb shell su -c 'cp /data/local/tmp/cacerts/* /system/etc/security/cacerts'
adb shell su -c 'chcon u:object_r:system_file:s0 /system/etc/security/cacerts/*'

adb push frida-server-16.2.1-android-x86 /data/local/tmp/frida-server
adb shell chmod +x /data/local/tmp/frida-server
adb shell su -c /data/local/tmp/frida-server
~~~

https://issuetracker.google.com/issues/331256113

## web

you can get `x-skyott-usertoken` with web client via `/auth/tokens`, but it
need `idsession` cookie. Looks like Android is the same.
