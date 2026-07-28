# mTLS

> I've just given this a quick test, and on my machine (with an emulator and a
> US VPN) I can start the app and start intercepting, but I quickly get a 403
> response from tv.clients.peacocktv.com, which seems to be an Akamai endpoint.
> It's hard to know for sure, but I'd best this is TLS fingerprinting (some
> context: https://httptoolkit.com/blog/tls-fingerprinting-node-js/) which
> Akamai uses to block all sorts of slightly unusual clients. 

The endpoint `tv.clients.peacocktv.com` enforces **mutual TLS (mTLS)** at the
CDN edge layer. The Peacock Android TV app presents a client certificate during
the TLS handshake that is validated against a specific CA ("NOWTV INTL
CERTIFICATE AUTHORITY"). Without this certificate, the CDN rejects the request
before it reaches the origin — regardless of HTTP headers.

The client certificate is **not stored as a file** in the APK. It is loaded at runtime from a native library (`libwebview.so`).

### 2. Extract the cert at runtime with Frida

Since the cert is generated/retrieved inside native code, use **Frida** to hook
the Java method (`com.sky.core.webview.TVWebView.load()`) and dump the objects
after the native call returns.

**Frida script (`extract_cert.js`):**

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

            var fKey = new File(filesDir + "/extracted_client.key", "wb");
            fKey.write(keyBuffer);
            fKey.flush();
            fKey.close();

            var fCert = new File(filesDir + "/extracted_client.pem", "wb");
            fCert.write(certBuffer);
            fCert.flush();
            fCert.close();

            console.log("[+] Cert and key saved to: " + filesDir);
        }

        return result;
    };
});
```

**Steps:**
1. Rooted Android emulator with Frida server running
2. Install the Peacock APK
3. Run: `frida -U -f com.peacocktv.peacockandroid -l extract_cert.js`
4. Pull the files:
   ```bash
   adb pull /data/user/0/com.peacocktv.peacockandroid/files/extracted_client.key extracted_client.key
   adb pull /data/user/0/com.peacocktv.peacockandroid/files/extracted_client.pem extracted_client.pem
   ```

The extracted files are raw DER-encoded binary.

### 3. Convert DER to PEM and use the cert

The raw DER files need to be wrapped in PEM headers. This was done with a small Go program:

```go
keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
os.WriteFile("client.key", keyPEM, 0644)
os.WriteFile("client.pem", certPEM, 0644)
```

### 4. Make the request with the client cert

```bash
curl --cert client.pem --key client.key "https://tv.clients.peacocktv.com/cvsdk/android/18.0.4/bundle.sdk-ext-peacock.js"
```
