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

## web

you can get `x-skyott-usertoken` with web client via `/auth/tokens`, but it
need `idsession` cookie. Looks like Android is the same.

## tv

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
