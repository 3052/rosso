# Peacock TV DRM License Capture via Frida

https://peacocktv.com/activate

## Why is this needed?

When analyzing the network traffic of the Peacock TV Android app
(`com.peacocktv.peacockandroid`) using a proxy like `mitmproxy`, certain
endpoints—specifically `play.clients.peacocktv.com` and
`tv.clients.peacocktv.com`—will return `403 Access Denied`. 

This is **not** standard Certificate Pinning. If it were pinning, the app would simply refuse to connect. Instead, the request goes through, but the server blocks it.

### The Real Problem: TLS / HTTP2 Fingerprinting

These endpoints are protected by Akamai. Akamai inspects the TLS ClientHello (JA3/JA4 fingerprint) and the HTTP/2 SETTINGS frame. 

When `mitmproxy` intercepts a TLS connection, it establishes its own TLS
session with the server using its own cryptographic library (usually OpenSSL).
The app's original network stack (Android's Conscrypt / OkHttp) has a very
specific TLS fingerprint. Because `mitmproxy`'s fingerprint does not match what
Akamai expects from an Android app, Akamai detects the proxy and returns a
`403` status code. 

**Mitmproxy cannot currently spoof Android/OkHttp TLS fingerprints.** There is no configuration flag or addon script that can bypass this.

### The Solution

Since we cannot proxy the connection without being blocked, we must intercept the traffic **inside the app's memory** before it gets encrypted. 

This Frida script hooks `okhttp3.internal.http.RealInterceptorChain.proceed()`.
Every HTTP request made by the app passes through this method. We check if the
request is destined for the Widevine license endpoint
(`https://play.clients.peacocktv.com/drm/widevine/acquirelicense`). If it is,
we extract the exact HTTP headers and the raw binary protobuf request body (the
Widevine challenge) and dump them to a text file.

Because the traffic is captured from memory and sent directly by the app,
Akamai sees the correct TLS fingerprint and allows the connection. We get to
see the plaintext request payload without fighting the TLS fingerprint.

## How to Use

### Prerequisites
* A rooted Android device or emulator (e.g., Android Studio Emulator with `adb root`)
* `frida-server` running on the device
* `frida-tools` installed on your PC (`pip install frida-tools`)
* `mitmproxy` installed on your PC

### Step 1: Capture non-blocked traffic with mitmproxy
Start mitmproxy, ignoring the hosts that block the proxy so they fall back to direct connections:

```powershell
mitmproxy --ignore-hosts play.clients.peacocktv.com --ignore-hosts tv.clients.peacocktv.com -w mitm_capture.mitm
```

### Step 2: Run the Frida script
In a separate terminal, spawn the app and load the capture script:

```powershell
frida -U -f com.peacocktv.peacockandroid -l license_capture.js
```

### Step 3: Trigger the request
Use the app and start playing a video. The app will make the Widevine license request. You will see the following in your Frida terminal:

```
[*] Captured: https://play.clients.peacocktv.com/drm/widevine/acquirelicense?bt=...
[*] Saved to: /storage/emulated/0/Android/data/com.peacocktv.peacockandroid/files/license_request.txt
```

### Step 4: Retrieve the file
Pull the raw HTTP request file from the device:

```powershell
adb pull /sdcard/Android/data/com.peacocktv.peacockandroid/files/license_request.txt
```

The `license_request.txt` file will contain the raw HTTP request, including the
binary protobuf body at the end, which can be used for replay attacks or
further analysis in tools like Go, Python, or Burp Suite.
