# magisk

## Step 1: Boot the emulator from the stock image
```bash
emulator -avd Television_1080p
```

## Step 2: Check `init` bitness
```bash
adb shell getprop ro.product.cpu.abilist64
```
- If the output is **empty** → use `-init-bits 32`
- If the output contains an ABI (like `x86_64` or `arm64-v8a`) → use `-init-bits 64`

## Step 3: Download the Magisk APK and the Go script, then run the patch script
Download the Magisk APK from the official releases page: <https://github.com/topjohnwu/magisk/releases>

Download the Go patch script from: <https://github.com/3052/verde/tree/main/2026/7/magisk>

With the emulator running, execute the script (replace `Magisk-v30.7.apk` with your downloaded filename):
```bash
magisk -apk Magisk-v30.7.apk -img ramdisk.img -init-bits 32
```

## Step 4: Replace the ramdisk in your system images directory
1. Close the emulator
2. Backup the original ramdisk: `cp ramdisk.img ramdisk.img.orig`
3. Locate your system image directory (e.g., `C:\Users\Steven\AppData\Local\Android\Sdk\system-images\android-34\android-tv\x86\`)
4. Copy the patched image to that directory, replacing the original:
   ```bash
   cp magisk_patched.img "C:\Users\Steven\AppData\Local\Android\Sdk\system-images\android-34\android-tv\x86\ramdisk.img"
   ```

## Step 5: Cold boot the emulator
```bash
emulator -avd Television_1080p -no-snapshot-load -show-kernel
```

## Step 6: Install the Magisk APK
```bash
adb install Magisk-v30.7.apk
```

## Step 7: Open the Magisk app
```bash
adb shell monkey -p com.topjohnwu.magisk -c android.intent.category.LAUNCHER 1
```
- It should now say **Installed: 30.7**
- It should prompt: *"Your device needs additional setup..."*
- Tap **OK** — emulator reboots

## Step 8: Verify root
```bash
adb shell su -c id
# uid=0(root) gid=0(root)
```

## Step 9: Configure proxy and custom CA certificates (optional)

**1. Download Cert-Fixer from its official repository (<https://github.com/pwnlogs/cert-fixer>) and install the module:**
```bash
adb push Cert-Fixer.zip /data/local/tmp/
adb shell su -c 'magisk --install-module /data/local/tmp/Cert-Fixer.zip'
```

**2. Push your certificate:**
```bash
adb push C:/Users/Steven/.mitmproxy/mitmproxy-ca-cert.pem /data/local/tmp/c8750f0d.0
```

**3. Move the certificate to the Cert-Fixer directory:**
```bash
adb shell su -c "mkdir -p /data/misc/user/0/cacerts-added"
adb shell su -c "cp /data/local/tmp/c8750f0d.0 /data/misc/user/0/cacerts-added/"
adb shell su -c "chown system:system /data/misc/user/0/cacerts-added/c8750f0d.0"
adb shell su -c "chmod 644 /data/misc/user/0/cacerts-added/c8750f0d.0"
```

**4. Reboot the emulator with the proxy:**
```bash
emulator -avd Television_1080p -no-snapshot-load -http-proxy http://127.0.0.1:8080
```

**5. Verify the certificate is in the APEX namespace:**
```bash
adb shell su -c 'ls /apex/com.android.conscrypt/cacerts | grep c8750f0d'
```
