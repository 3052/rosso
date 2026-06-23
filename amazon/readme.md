# amazon

> Everything you talk about, cut that<br>
> That door you trying to open, you could shut that
>
> [Logic (2017)](//youtube.com/watch?v=wH4kzAb4l0E)

https://apkmirror.com/apk/amazon-mobile-llc/prime-video-android-tv-android-tv

even though its a TV app - you can use it with a phone device - just install as
normal and start:

~~~
adb shell monkey -p com.amazon.amazonvideo.livingroom `
-c android.intent.category.LEANBACK_LAUNCHER 1
~~~

or stop/clear:

~~~
adb shell pm clear com.amazon.amazonvideo.livingroom
~~~

create Pixel 5, Android 11 device. install system certificate

1. https://github.com/httptoolkit/frida-interception-and-unpinning/issues/207
2. https://issuetracker.google.com/issues/331256113
3. https://issuetracker.google.com/issues/522344738
4. https://github.com/frida/frida-core/issues/1240
5. <https://gitlab.com/newbit/rootAVD/-/work_items/117>

## patch

### The binary

The app ships a **15MB ARM native library called `libignite.so`** inside `split_config.armeabi_v7a.apk`. This single binary contains the entire app engine including a statically linked copy of **curl 8.9.1-DEV** and **OpenSSL** (BoringSSL variant). Symbols are fully stripped — we confirmed this by finding zero exported curl/SSL symbols when enumerating the module. All the identification was done by scanning for embedded strings like `curl version is %s`, `curl/curl/lib/vtls/openssl.c`, `SSL certificate problem: %s`, etc.

### How the trusted CAs are loaded

The app bundles a CA certificate store inside `assets/ignite-assets.tar` in the base APK. This tar contains `bin/certs/` with ~168 individual PEM certificate files in OpenSSL's `c_rehash` hash-named format (`XXXXXXXX.0`). At first launch, the tar is extracted to `/data/data/com.amazon.amazonvideo.livingroom/files/bin/certs/`.

However, **the certs are NOT loaded from disk at verification time**. We proved this by hooking both `openat()` and `fopen()` in libc — no cert files were ever opened during an SSL connection. The only file access was `/usr/local/ssl/openssl.cnf` (which doesn't exist on Android).

Instead, the app reads the cert files into memory at startup and passes them to curl via **`CURLOPT_CAINFO_BLOB`** — an in-memory blob. We confirmed this by:
- Finding the string `CURLOPT_SSLCERT_BLOB` in the binary
- Finding the `CApath` string reference and the path `/data/user/0/com.amazon.amazonvideo.livingroom/files/bin/certs/` in writable heap memory at runtime
- Finding zero PEM certificate bundles in any memory-mapped region (the blob is likely DER or loaded/freed before we scanned)
- Proving that adding certs to the on-disk directory had no effect on validation

### How verification is configured

The statically linked curl is configured with `CURLOPT_SSL_VERIFYPEER = 1` via `curl_easy_setopt()`. We found three call sites in the `.text` section of `libignite.so` with this exact ARM Thumb instruction pattern:

```
MOVS R1, #64    (CURLOPT_SSL_VERIFYPEER = 64)
MOVS R2, #1     (value = true)
BL curl_easy_setopt
```

When verification is enabled, curl's OpenSSL backend calls `SSL_get_verify_result()` after the TLS handshake. If the server's certificate chain doesn't validate against the in-memory CA blob, it returns `CURLE_PEER_FAILED_VERIFICATION` (error code 60) with the message `SSL certificate problem: %s`.

### Why it was hard to intercept

The ARM code runs on an **x86 emulator via `libndk_translation.so`** (Google's ARM-to-x86 binary translator). This means:
- The translated ARM code doesn't appear as normal x86 modules — Frida can't enumerate its exports
- ARM OpenSSL's functions are invisible to Frida's `Module.findExportByName()`
- The translated code doesn't call through the x86 system `libssl.so` — it has its own OpenSSL compiled in
- File I/O from translated ARM code goes through syscalls differently than x86 libc wrappers
- The Frida build on this device (17.11.0) didn't have the `Java` bridge available, ruling out Java-layer SSL hooks

### The fix

Patching 3 bytes in the APK — changing `MOVS R2, #1` to `MOVS R2, #0` at each `curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 1)` call site — disables SSL peer verification entirely. The patch is done in-place on the APK file to preserve zip offsets and alignment, since Android loads `.so` files directly from the APK via memory mapping.
