# Paramount+

## how to get cmsAccountId

its in the HTML response body:

<https://paramountplus.com/shows/video/8PO2sBBr6lFb7J4nklXuzNZRhUR_V9dd>

## How to get secret\_key?

~~~
com\cbs\app\androiddata\retrofit\util\RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

## how to get app secret?

us:

~~~
sources\com\cbs\app\config\UsaMobileAppConfigProvider.java
~~~

- https://apkmirror.com/apk/cbs-interactive-inc/paramount
- https://play.google.com/store/apps/details?id=com.cbs.app

international:

~~~
sources/com/cbs/app/config/DefaultAppSecretProvider.java
~~~

- https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv
- https://play.google.com/store/apps/details?id=com.cbs.ca

## paypal.com US

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
