# Paramount+

## how to get cmsAccountId

its in the HTML response body:

<https://paramountplus.com/shows/video/8PO2sBBr6lFb7J4nklXuzNZRhUR_V9dd>

## How to get secret\_key?

~~~
sources/com/cbs/app/androiddata/retrofit/util/RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

## how to get app secret?

us:

~~~
sources/com/cbs/app/config/UsaMobileAppConfigProvider.java
~~~

- https://apkmirror.com/apk/cbs-interactive-inc/paramount
- https://play.google.com/store/apps/details?id=com.cbs.app

international:

~~~
sources/com/cbs/app/config/DefaultAppSecretProvider.java
~~~

- https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv
- https://play.google.com/store/apps/details?id=com.cbs.ca

## paramount-4

international

https://apkmirror.com/apk/viacomcbs-streaming/paramount-4

## paramount-3

old

https://apkmirror.com/apk/viacomcbs-streaming/paramount-3

## paramount-2

android TV

https://apkmirror.com/apk/viacomcbs-streaming/paramount-2

## paramount

US

https://apkmirror.com/apk/viacomcbs-streaming/paramount

APK lies, you need at least Android 12 (level 31)

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
~/.android/avd/Pixel_XL.avd/emulator-user.ini
~~~

to:

~~~
window.x = 0
window.y = 0
~~~

install system certificate

## cbs

https://apkmirror.com/apk/cbs-interactive-inc/cbs

Create Android 7 device. Install system certificate

## paypal.com United Kingdom

1. United Kingdom VPN
2. address
   - https://paramountplus.com/gb/movies
      - if you go here it says free
   - https://paramountplus.com/gb
      - $2.99
3. full name
4. email
5. password
6. By pressing “Agree & Continue”, you confirm that you have read and agree to
   the Terms of Use
7. agree & continue
8. continue
9. premium
   - select plan
10. monthly
   - continue
11. paypal
12. DELETE ALL PAYPAL COOKIES
13. continue to paypal
14. agree and continue
15. buy now
16. paypal.com/myaccount/autopay
17. paramount
18. stop paying with paypal
19. stop paying with paypal (again)
20. stop paying with paypal (again again)

## paypal.com United States

1. about:config
2. general.useragent.override
3. string
4. add
5. Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:147.0) Gecko/20100101 Firefox/147.0
6. paramountplus.com
7. get started
8. paramount+ premium
   - continue
9. full name
10. email
   - mail.tm
11. password
12. zip code
13. birthdate
14. gender
15. agree & continue
16. paypal
17. continue to paypal
18. agree and continue
19. subscribe
20. paypal.com/myaccount/autopay
21. paramount
22. stop paying with paypal
