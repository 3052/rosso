# MITM Proxy

First download:

https://mitmproxy.org/downloads

then start `mitmproxy.exe`. Address and port should be in the bottom right
corner. Default should be `*:8080`. Assuming the above, go to Android Emulator
and set proxy:

~~~
127.0.0.1:8080
~~~

Then open Google Chrome on Virtual Device, and browse to http://example.com. To
exit, press `q`, then `y`. To capture HTTPS, open Google Chrome on Virtual
Device, and browse to <http://mitm.it>. Click on the Android certificate. Under
"Certificate name" enter "MITM", then click "OK". Then browse to
https://example.com. To disable compression:

~~~
set anticomp true
~~~

## Blocklist

Press `O` to enter options. Move to `block_list` and press Enter. Then press
`a` to add a new entry. Press Esc when finished, then `q`.

~~~
/~u tag.js/444
~~~

https://docs.mitmproxy.org/stable/overview-features#blocklist

## Modify Headers

Press `O` to enter options. Move to `modify_headers` and press Enter. Then press
`a` to add a new entry. Press Esc when finished, then `q`.

~~~
/~q/X-Forwarded-For/99.224.0.0
/~u hello.world/X-Forwarded-For/99.224.0.0
~~~

https://docs.mitmproxy.org/stable/overview-features#modify-headers
