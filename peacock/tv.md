# tv

~~~
$env:path = 'C:\windows\system32'
.\rootAVD.bat system-images\android-29\android-tv\x86\ramdisk.img

adb shell mkdir -p /data/local/tmp/cacerts
adb push C:/Users/Steven/.mitmproxy/mitmproxy-ca-cert.pem /data/local/tmp/cacerts/c8750f0d.0
adb shell cp /system/etc/security/cacerts/* /data/local/tmp/cacerts

adb shell su -c 'mount -t tmpfs tmpfs /system/etc/security/cacerts'
adb shell su -c 'cp /data/local/tmp/cacerts/* /system/etc/security/cacerts'
adb shell su -c 'chcon u:object_r:system_file:s0 /system/etc/security/cacerts/*'

adb push frida-server-16.2.1-android-x86 /data/local/tmp/frida-server
adb shell chmod +x /data/local/tmp/frida-server
adb shell su -c /data/local/tmp/frida-server
~~~

https://issuetracker.google.com/issues/331256113
