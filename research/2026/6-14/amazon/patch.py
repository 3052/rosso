import zipfile
import subprocess
import os

PACKAGE_NAME = "com.amazon.amazonvideo.livingroom"

def shell(cmd):
    return subprocess.run(["adb", "shell"] + cmd.split(),
                         capture_output=True, text=True).stdout.strip()

def main():
    # Step 1: Find and pull APK
    print("[*] Finding APK path...")
    paths = shell(f"pm path {PACKAGE_NAME}").split("\n")
    arm_path = [p.replace("package:", "") for p in paths if "armeabi" in p]
    
    if not arm_path:
        print("[-] Could not find armeabi APK split")
        return
    
    apk_device_path = arm_path[0]
    print(f"[*] APK: {apk_device_path}")
    
    if not os.path.exists("app.apk"):
        print("[*] Pulling APK...")
        subprocess.run(["adb", "pull", apk_device_path, "app.apk"])
    
    # Step 2: Find libignite.so inside APK
    with zipfile.ZipFile("app.apk", "r") as z:
        info = z.getinfo("lib/armeabi-v7a/libignite.so")
    
    with open("app.apk", "rb") as f:
        f.seek(info.header_offset)
        lh = f.read(30)
        data_offset = info.header_offset + 30 + int.from_bytes(lh[26:28], "little") + int.from_bytes(lh[28:30], "little")
        f.seek(data_offset)
        assert f.read(4) == b"\x7fELF"
    
    # Step 3: Find and patch CURLOPT_SSL_VERIFYPEER = 1
    with open("app.apk", "rb") as f:
        apk = bytearray(f.read())
    
    so_data = apk[data_offset:data_offset + info.file_size]
    
    patched = 0
    for i in range(len(so_data) - 3):
        if (so_data[i] == 0x40 and so_data[i+1] == 0x21 and
            so_data[i+2] == 0x01 and so_data[i+3] == 0x22):
            apk[data_offset + i + 2] = 0x00
            patched += 1
            print(f"[*] Patched CURLOPT_SSL_VERIFYPEER at 0x{i:x}")
    
    if patched == 0:
        print("[-] No CURLOPT_SSL_VERIFYPEER found!")
        return
    
    print(f"[+] Patched {patched} locations")
    
    with open("app_patched.apk", "wb") as f:
        f.write(apk)
    
    # Step 4: Backup and push
    print("[*] Backing up original...")
    shell(f"cp {apk_device_path} {apk_device_path}.bak")
    
    print("[*] Pushing patched APK...")
    subprocess.run(["adb", "push", "app_patched.apk", apk_device_path])
    
    shell(f"am force-stop {PACKAGE_NAME}")
    shell(f"rm -rf /data/data/{PACKAGE_NAME}/code_cache")
    
    print("[*] Starting app...")
    shell(f"am start -n {PACKAGE_NAME}/com.amazon.ignition.IgnitionActivity")
    
    print("[+] Done! Check mitmproxy for traffic.")
    print(f"[*] To restore: adb shell cp {apk_device_path}.bak {apk_device_path}")

if __name__ == "__main__":
    main()
