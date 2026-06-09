# magisk

### Step 1: Prepare the Files on Windows
1. On your Windows Desktop, create an empty folder named `Patch`.
2. Download the official Magisk APK (e.g., v30.7) and change the extension to `.zip`. Extract it.
3. Open the extracted Magisk folder and navigate to `lib\x86\` (Do NOT use `x86_64`).
4. Copy these 4 files into your `Patch` folder and rename them exactly like this:
   * `libmagiskboot.so` -> `magiskboot`
   * `libmagiskinit.so` -> `magiskinit`
   * `libmagisk.so` -> `magisk32`
   * `libinit-ld.so` -> `init-ld`
5. Navigate to the `assets\` folder inside the extracted Magisk zip, copy `stub.apk`, and paste it into your `Patch` folder.
6. Open your TV emulator's SDK folder (usually
   `%LOCALAPPDATA%\Android\Sdk\system-images\android-34\android-tv\x86\`). Copy
   the `ramdisk.img` file and paste it into your `Patch` folder.

You should now have exactly 6 files in your Patch folder.

### Step 2: Push Files to the Running Emulator
1. Boot up your API 34 Android TV emulator normally in Android Studio.
2. Open a Windows PowerShell window inside your `Patch` folder.
3. Run these exact commands one by one to send the files to a temporary folder on the TV:

`adb shell mkdir -p /data/local/tmp/patch`

`adb push magiskboot magiskinit magisk32 init-ld stub.apk ramdisk.img /data/local/tmp/patch/`

`adb shell chmod +x /data/local/tmp/patch/magiskboot`

### Step 3: Execute the Manual CPIO Injection

Run these commands one at a time. (The multi-line injection command maps the
32-bit binary to all 64-bit aliases to bypass the kernel architecture check.
You can copy and paste the entire block exactly as it is).

Decompress the archive:
`adb shell "cd /data/local/tmp/patch && ./magiskboot decompress ramdisk.img ramdisk.cpio"`

Create the Magisk configuration file:
`adb shell 'cd /data/local/tmp/patch && printf "KEEPVERITY=true\nKEEPFORCEENCRYPT=true\nRECOVERYMODE=false\n" > config'`

Inject the binaries and hijack the bootloader:
```powershell
adb shell 'cd /data/local/tmp/patch && ./magiskboot cpio ramdisk.cpio' `
  '"mkdir 0000 .backup"' `
  '"mv init .backup/init"' `
  '"add 0000 .backup/.magisk config"' `
  '"mkdir 0750 overlay.d"' `
  '"mkdir 0750 overlay.d/sbin"' `
  '"add 0750 init magiskinit"' `
  '"add 0755 overlay.d/sbin/magisk magisk32"' `
  '"add 0755 overlay.d/sbin/magisk32 magisk32"' `
  '"add 0755 overlay.d/sbin/magisk64 magisk32"' `
  '"add 0755 overlay.d/sbin/init-ld init-ld"' `
  '"add 0644 overlay.d/sbin/stub.apk stub.apk"'
```

Compress the archive back to LZ4:
`adb shell "cd /data/local/tmp/patch && ./magiskboot compress=lz4_legacy ramdisk.cpio magisk_patched.img"`

Pull the finished file to Windows:
`adb pull /data/local/tmp/patch/magisk_patched.img`

### Step 4: Flash the Emulator & Wipe Data
1. Completely close the running TV emulator.
2. Go to your Windows SDK folder (`%LOCALAPPDATA%\Android\Sdk\system-images\android-34\android-tv\x86\`).
3. Rename the original `ramdisk.img` to `ramdisk.img.stock` (to keep as a backup).
4. Move your new `magisk_patched.img` from your `Patch` folder into this SDK folder.
5. Rename `magisk_patched.img` to exactly `ramdisk.img`.
6. Open Android Studio Device Manager. Click the **three dots** next to your TV and click **Wipe Data**. (Crucial: Skipping this causes a permanent black screen).
7. Click the **three dots** again and click **Cold Boot Now**.

### Step 5: Finalize Magisk
1. Once the TV reaches the home screen, install the Magisk app from your PowerShell window (rename your zip back to apk first):
   `adb install Magisk-v30.7.apk`
2. Open the Magisk app on the TV. It will say Installed 30.7. 
3. A pop-up will appear saying "Requires additional setup". Click **CANCEL**. (Clicking OK will crash the emulator).
4. Navigate to the **Superuser** tab (the shield icon at the bottom).
5. Toggle the switch ON for **[SharedUID] Shell (com.android.shell)**.

You now have unrestricted `su` root access on an API 34 Android TV emulator.
Install your ARM APK, launch it, and use the Superuser tab to grant it root
permissions.
