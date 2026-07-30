# Peacock TV

- https://apkmirror.com/apk/peacock-tv-llc/peacock-tv-android-tv
- https://play.google.com/store/apps/details?id=com.peacocktv.peacockandroid

create Pixel 5, APK lies you need at least API 31. install system certificate

~~~
emulator -avd Pixel_5 -no-snapshot-load -http-proxy http://127.0.0.1:8080
~~~

then:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell am force-stop com.peacocktv.peacockandroid
adb shell pm clear com.peacocktv.peacockandroid
~~~

## Sign-In Methods

The `peacock` Go package supports two sign-in methods:

- **Companion Device Activation** — Start a journey on the TV, get a 6-digit
  code, visit `peacocktv.com/activate` on another device, and enter the code
  to link your account. The TV never sees your email or password. Implemented
  via `StartJourney` + `Activate`.

- **Email OTP** — Enter your email on the device, receive a 6-digit code by
  email, and enter it to authenticate. Implemented via `InitiateOTP` +
  `VerifyOTP`.

## Background

The endpoints `play.clients.peacocktv.com` and `tv.clients.peacocktv.com`
enforce **mutual TLS (mTLS)** at the CDN edge layer. The Peacock Android TV
app presents a client certificate during the TLS handshake that is validated
against a specific CA ("NOWTV INTL CERTIFICATE AUTHORITY"). Without this
certificate, the CDN rejects the request before it reaches the origin —
regardless of HTTP headers.

The client certificate is **not stored as a file** in the APK. It is loaded at
runtime from a native library (`libwebview.so`).

## Extract the cert at runtime with Frida

Since the cert is generated/retrieved inside native code, use **Frida** to hook
the Java method (`com.sky.core.webview.TVWebView.load()`) and dump the objects
after the native call returns.

**`extract_cert.js`**

```javascript
Java.perform(function() {
    var TVWebView = Java.use("com.sky.core.webview.TVWebView");
    var OpenSSLRSAPrivateKey = Java.use("com.android.org.conscrypt.OpenSSLRSAPrivateKey");
    var OpenSSLX509Certificate = Java.use("com.android.org.conscrypt.OpenSSLX509Certificate");

    TVWebView.load.implementation = function() {
        var result = this.load();

        if (result !== null && result.length >= 2) {
            var privateKey = Java.cast(result[0], OpenSSLRSAPrivateKey);
            var cert = Java.cast(result[1], OpenSSLX509Certificate);

            var keyBytes = privateKey.getEncoded();
            var certBytes = cert.getEncoded();

            // Convert Java byte[] to JavaScript ArrayBuffer
            var keyBuffer = new ArrayBuffer(keyBytes.length);
            var keyView = new Uint8Array(keyBuffer);
            for (var i = 0; i < keyBytes.length; i++) {
                keyView[i] = keyBytes[i] & 0xff;
            }

            var certBuffer = new ArrayBuffer(certBytes.length);
            var certView = new Uint8Array(certBuffer);
            for (var i = 0; i < certBytes.length; i++) {
                certView[i] = certBytes[i] & 0xff;
            }

            // Write to app's private files directory
            var ActivityThread = Java.use("android.app.ActivityThread");
            var context = ActivityThread.currentApplication().getApplicationContext();
            var filesDir = context.getFilesDir().getAbsolutePath();

            var fKey = new File(filesDir + "/key.der", "wb");
            fKey.write(keyBuffer);
            fKey.flush();
            fKey.close();

            var fCert = new File(filesDir + "/cert.der", "wb");
            fCert.write(certBuffer);
            fCert.flush();
            fCert.close();

            console.log("[+] Cert and key saved to: " + filesDir);
        }

        return result;
    };
});
```

1. Rooted Android emulator with Frida server running
2. Install the Peacock APK
3. Run:

   ```powershell
   frida -U -f com.peacocktv.peacockandroid -l extract_cert.js
   ```

4. Pull the files:

   ```powershell
   adb pull /data/user/0/com.peacocktv.peacockandroid/files/key.der
   adb pull /data/user/0/com.peacocktv.peacockandroid/files/cert.der
   ```

The extracted files are raw DER-encoded binary.

## Convert DER to PEM

The raw DER files need to be converted to PEM format for use with mitmproxy
and curl. Use the `der2pem` Go program (see `der2pem.go`):

```powershell
der2pem -cert cert.der -key key.der
```

This writes two files to the current directory:

- `cert.pem` — the certificate
- `key.pem` — the private key

> **Note:** mitmproxy expects a **single combined PEM file** (cert + key
> concatenated) named `<hostname>.pem` per target host. The same cert
> is used for all Peacock endpoints. Use the following PowerShell to
> combine the two files for each host:

```powershell
Get-Content cert.pem, key.pem | Set-Content tv.clients.peacocktv.com.pem
Get-Content cert.pem, key.pem | Set-Content play.clients.peacocktv.com.pem
```

The combined file looks like:

```
-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----
```

## Make the request with the client cert

**With mitmproxy**

mitmproxy requires the combined `<hostname>.pem` files generated in
the previous section:

```powershell
mitmproxy --set client_certs=.
```

> mitmproxy matches each outgoing connection's target hostname to a
> `<hostname>.pem` file in the `client_certs` directory.

**With curl**

curl accepts the separate cert and key files directly via `--cert` and `--key`, so no combining step is needed:

```powershell
curl --cert cert.pem --key key.pem `
    "https://tv.clients.peacocktv.com/cvsdk/android/18.0.4/bundle.sdk-ext-peacock.js"

curl --cert cert.pem --key key.pem `
    "https://play.clients.peacocktv.com/..."
```
