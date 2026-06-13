import frida
import sys

# The exact list of scripts you provided
script_files = [
    "./config.js",
    "./native-connect-hook.js",
    "./native-tls-hook.js",
    "./android/android-proxy-override.js",
    "./android/android-system-certificate-injection.js",
    "./android/android-certificate-unpinning.js",
    "./android/android-certificate-unpinning-fallback.js",
    "./android/android-disable-root-detection.js"
]

def main():
    print("[*] Connecting to device...")
    device = frida.get_usb_device()

    print("[*] Spawning Amazon Video TV app...")
    pid = device.spawn(
        "com.amazon.amazonvideo.livingroom", 
        activity="com.amazon.ignition.IgnitionActivity"
    )
    session = device.attach(pid)

    print("[*] Compiling script payload...")
    combined_script = ""
    
    # PHASE 1: Load config.js and Native Hooks immediately
    for file_path in script_files:
        if "android" not in file_path.lower():
            try:
                with open(file_path, "r", encoding="utf-8") as f:
                    combined_script += f"\n// --- {file_path} (LOADED IMMEDIATELY) ---\n"
                    combined_script += f.read()
            except FileNotFoundError:
                print(f"[-] ERROR: Could not find {file_path}")
                device.kill(pid)
                sys.exit(1)

    # PHASE 2: Wrap Android hooks in a function so they don't execute right away
    combined_script += "\n\n// --- DELAYED ANDROID HOOKS WRAPPER ---\n"
    combined_script += "function startAndroidHooks() {\n"
    
    for file_path in script_files:
        if "android" in file_path.lower():
            try:
                with open(file_path, "r", encoding="utf-8") as f:
                    combined_script += f"\n  console.log('[*] Injecting {file_path}...');\n"
                    # Wrap each script in an anonymous function to prevent variable scope conflicts
                    combined_script += "  try {\n    (function() {\n"
                    combined_script += f.read()
                    combined_script += "\n    })();\n  } catch(e) {\n"
                    combined_script += f"    console.error('[-] Error in {file_path}:', e);\n"
                    combined_script += "  }\n"
            except FileNotFoundError:
                print(f"[-] ERROR: Could not find {file_path}")
                device.kill(pid)
                sys.exit(1)

    combined_script += "\n}\n"

    # PHASE 3: Polling mechanism to wait for the JVM
    combined_script += """
    function waitForJava() {
        // Check if the JVM has been successfully loaded into memory
        if (typeof Java !== 'undefined' && Java.available) {
            console.log("[*] JVM is ready! Executing Android hooks...");
            startAndroidHooks();
        } else {
            // Not ready yet, check again in 50ms
            setTimeout(waitForJava, 50);
        }
    }
    
    console.log("[*] Native hooks injected. Waiting for JVM to boot...");
    waitForJava();
    """

    print("[*] Injecting scripts...")
    script = session.create_script(combined_script)
    
    # Handle console messages from your JavaScript files
    def on_message(message, data):
        if message['type'] == 'send':
            print(f"[*] {message['payload']}")
        elif message['type'] == 'error':
            print(f"[-] ERROR: {message['stack']}")

    script.on('message', on_message)
    script.load()

    # App is finally resumed here; the 50ms polling loop will catch the JVM booting up
    print("[*] Resuming app execution...")
    device.resume(pid)

    print("[*] Interception active. Press Ctrl+C to exit.")
    try:
        sys.stdin.read()
    except KeyboardInterrupt:
        print("\n[*] Exiting...")
        session.detach()

if __name__ == "__main__":
    main()