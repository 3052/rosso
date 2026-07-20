"""
com.wbd.stream (Max) SSL Certificate Unpinning Script

Patches CURLOPT_SSL_VERIFYPEER/HOST/STATUS in the statically-linked cURL
inside libhbomax.so to disable SSL certificate verification, allowing
mitmproxy to intercept HTTPS traffic.

Usage:
    python patch-max-ssl.py

Prerequisites:
    - ADB connected to device/emulator
    - App installed (com.wbd.stream) and run at least once
    - mitmproxy running on 127.0.0.1:8080
    - mitmproxy CA cert installed in system trust store
    - Device rooted (adb root)

Note: Re-run this script after reinstalling/updating the app.
"""
import zipfile
import subprocess
import os
import sys
import time

PACKAGE_NAME = "com.wbd.stream"
ACTIVITY = "com.wbd.stream/com.wbd.beam.BeamActivity"

ARCH_NAME = "armeabi-v7a"
ARCH_KEYWORD = "armeabi_v7a"
LIB_PATH_IN_APK = f"lib/{ARCH_NAME}/libhbomax.so"

# ARM Thumb-2 patterns to patch:
#   movs r1, #OPT   (0xNN 0x21)
#   movs r2, #0x01  (0x01 0x22)  -> patch 0x01 to 0x00
PATCH_PATTERNS = [
    (0x40, "CURLOPT_SSL_VERIFYPEER"),
    (0x51, "CURLOPT_SSL_VERIFYHOST"),
    (0x66, "CURLOPT_SSL_VERIFYSTATUS"),
]

def shell(cmd):
    """Run an adb shell command and return stdout."""
    return subprocess.run(
        ["adb", "shell"] + cmd.split(),
        capture_output=True, text=True
    ).stdout.strip()

def find_extracted_so():
    """Find the path where libhbomax.so is extracted on the filesystem."""
    result = shell("find /data/app -name libhbomax.so")
    if result:
        return result.split("\n")[0].strip()
    return None

def find_apk_path():
    """Find the armeabi-v7a split APK."""
    output = shell(f"pm path {PACKAGE_NAME}")
    if not output:
        print(f"[-] App {PACKAGE_NAME} not found. Is it installed?")
        sys.exit(1)

    paths = [p.replace("package:", "") for p in output.split("\n")]
    for p in paths:
        if ARCH_KEYWORD in p:
            return p

    print(f"[-] No matching {ARCH_KEYWORD} APK found among:")
    for p in paths:
        print(f"    {p}")
    sys.exit(1)

def pull_apk(device_path, local_name):
    """Pull an APK from the device."""
    if os.path.exists(local_name):
        print(f"[*] Using cached {local_name}")
        return True

    print(f"[*] Pulling {device_path}...")
    result = subprocess.run(
        ["adb", "pull", device_path, local_name],
        capture_output=True, text=True
    )
    return result.returncode == 0

def extract_so(apk_path, lib_zip_path):
    """Extract libhbomax.so from an APK."""
    print(f"[*] Extracting {lib_zip_path} from {apk_path}...")

    with zipfile.ZipFile(apk_path, "r") as z:
        if lib_zip_path not in z.namelist():
            print(f"[-] {lib_zip_path} not found in APK")
            return None

        so_data = z.read(lib_zip_path)

    if so_data[:4] != b"\x7fELF":
        print(f"[-] Not an ELF file (got {so_data[:4]})")
        return None

    print(f"[+] Extracted {len(so_data)} bytes (ELF confirmed)")
    return so_data

def patch_so(so_data):
    """Patch CURLOPT_SSL_VERIFYPEER/HOST/STATUS in the .so binary."""
    so = bytearray(so_data)
    total_patched = 0

    print("[*] Patching cURL SSL verification flags...")

    for opt_val, opt_name in PATCH_PATTERNS:
        patched = 0
        for i in range(len(so) - 3):
            # Pattern: movs r1, #OPT (0xNN 0x21) + movs r2, #0x01 (0x01 0x22)
            if (so[i] == opt_val and so[i+1] == 0x21 and
                so[i+2] == 0x01 and so[i+3] == 0x22):
                so[i + 2] = 0x00  # Change value from 0x01 to 0x00
                patched += 1
                print(f"  [+] Patched {opt_name} at offset 0x{i:x}")

        if patched == 0:
            print(f"  [-] {opt_name}: not found")
        else:
            print(f"  [+] {opt_name}: patched {patched} location(s)")
        total_patched += patched

    if total_patched == 0:
        print("[-] No cURL SSL verification patterns found!")
        return None

    print(f"[+] Total: {total_patched} patch(es)")
    return bytes(so)

def push_patched_so(patched_data, target_path):
    """Push the patched .so to the device."""
    temp_file = "libhbomax_patched.so"
    with open(temp_file, "wb") as f:
        f.write(patched_data)

    print(f"\n[*] Rooting device...")
    shell("root")
    time.sleep(1)

    print(f"[*] Pushing patched .so to {target_path}...")
    subprocess.run(["adb", "push", temp_file, target_path], capture_output=True)

    print(f"[*] Setting permissions...")
    shell(f"chmod 755 {target_path}")

    print(f"[*] Clearing code cache...")
    shell(f"rm -rf /data/data/{PACKAGE_NAME}/code_cache")

def main():
    print("=" * 60)
    print("  Max (com.wbd.stream) SSL Unpinning Script")
    print("  Patches CURLOPT_SSL_VERIFYPEER in libhbomax.so")
    print("=" * 60)

    # Step 1: Find where the .so is extracted
    print("\n[*] Step 1: Finding extracted libhbomax.so...")
    extracted_path = find_extracted_so()
    if not extracted_path:
        print("[-] Could not find extracted libhbomax.so on filesystem")
        print("    Is the app installed? Has it been run at least once?")
        sys.exit(1)
    print(f"  Found: {extracted_path}")

    # Step 2: Find the matching APK
    print("\n[*] Step 2: Finding APK...")
    apk_path = find_apk_path()
    print(f"  Using: {apk_path}")

    # Step 3: Pull the APK
    print(f"\n[*] Step 3: Pulling APK...")
    local_apk = os.path.basename(apk_path)
    if not pull_apk(apk_path, local_apk):
        print("[-] Failed to pull APK")
        sys.exit(1)

    # Step 4: Extract libhbomax.so
    print(f"\n[*] Step 4: Extracting libhbomax.so from APK...")
    so_data = extract_so(local_apk, LIB_PATH_IN_APK)
    if not so_data:
        sys.exit(1)

    # Step 5: Patch
    print(f"\n[*] Step 5: Patching binary...")
    patched_data = patch_so(so_data)
    if not patched_data:
        sys.exit(1)

    # Step 6: Push to device
    print(f"\n[*] Step 6: Replacing .so on device...")
    push_patched_so(patched_data, extracted_path)

    print("\n" + "=" * 60)
    print("  DONE! The app should now work through mitmproxy.")
    print("  No Frida hooks needed — the patch is static.")
    print("=" * 60)
    print("\nNote: Re-run this script after reinstalling/updating the app.")

if __name__ == "__main__":
    main()
