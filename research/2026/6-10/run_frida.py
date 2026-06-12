import frida
import sys

# --- CONFIGURATION ---
package_name = "com.amazon.amazonvideo.livingroom"
activity_name = "com.amazon.ignition.IgnitionActivity" # <-- UPDATE THIS FROM THE ADB COMMAND

# List of your scripts in the exact order you want them loaded
script_files = [
    "./android/android-certificate-unpinning-fallback.js",
    "./android/android-certificate-unpinning.js",
    "./android/android-disable-root-detection.js",
    "./android/android-proxy-override.js",
    "./android/android-system-certificate-injection.js",
    "./config.js",
    "./native-connect-hook.js",
    "./native-tls-hook.js"
]

def on_message(message, data):
    """Handles console.log and errors from the injected JavaScript"""
    if message['type'] == 'send':
        print(f"[*] {message['payload']}")
    elif message['type'] == 'error':
        print(f"[-] {message['stack']}")
    else:
        print(message)

def main():
    # 1. Combine all JavaScript files into a single string
    print("[*] Combining scripts...")
    combined_script_code = ""
    for file_path in script_files:
        try:
            with open(file_path, "r", encoding="utf-8") as f:
                # Adding comments to help trace errors back to the original file
                combined_script_code += f"\n\n// --- START OF {file_path} ---\n"
                combined_script_code += f.read()
                combined_script_code += f"\n// --- END OF {file_path} ---\n"
        except FileNotFoundError:
            print(f"[-] Error: Could not find {file_path}. Are you in the right directory?")
            sys.exit(1)

    try:
        # 2. Connect to USB device
        device = frida.get_usb_device()
        
        # 3. Force spawn the app using the specific activity
        print(f"[*] Spawning {package_name} via {activity_name}...")
        pid = device.spawn([package_name], activity=activity_name)
        
        # 4. Attach to the new process
        print(f"[*] Attaching to PID {pid}...")
        session = device.attach(pid)
        
        # 5. Create and load the combined script
        print("[*] Injecting scripts (Early Instrumentation)...")
        script = session.create_script(combined_script_code)
        script.on('message', on_message) # Attach the message handler so we see console.log
        script.load()
        
        # 6. Resume the app execution now that bypasses are in place
        print("[*] Resuming app execution...")
        device.resume(pid)
        
        print("[*] App is running! Press Enter to detach and quit.")
        sys.stdin.read()
        
    except frida.ServerNotRunningError:
        print("[-] Frida server is not running on the Android device.")
    except frida.ExecutableNotFoundError:
        print(f"[-] Could not find the package {package_name} on the device.")
    except Exception as e:
        print(f"[-] An unexpected error occurred: {e}")

if __name__ == "__main__":
    main()
