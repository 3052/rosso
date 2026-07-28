Java.perform(function () {
  var RealInterceptorChain = Java.use("okhttp3.internal.http.RealInterceptorChain");
  var BodyClass = Java.use("okhttp3.RequestBody$Companion$toRequestBody$2");
  var ArrayCls = Java.use("java.lang.reflect.Array");
  var FileOutputStream = Java.use("java.io.FileOutputStream");
  var File = Java.use("java.io.File");
  var StringCls = Java.use("java.lang.String");

  var bodyField = BodyClass.class.getDeclaredField("$this_toRequestBody");
  bodyField.setAccessible(true);
  var offsetField = BodyClass.class.getDeclaredField("$offset");
  offsetField.setAccessible(true);
  var countField = BodyClass.class.getDeclaredField("$byteCount");
  countField.setAccessible(true);

  var TARGET_URL = "https://play.clients.peacocktv.com/drm/widevine/acquirelicense";
  var done = false;

  RealInterceptorChain.proceed.overloads.forEach(function (overload) {
    overload.implementation = function () {
      var response = overload.apply(this, arguments);

      if (done) return response;

      try {
        var request = this.request();
        var url = request.url().toString();

        if (!url.startsWith(TARGET_URL)) return response;
        done = true;

        var dir = Java.use("android.app.ActivityThread")
          .currentApplication()
          .getApplicationContext()
          .getExternalFilesDir(null)
          .getAbsolutePath();

        var filePath = dir + "/license_request.txt";

        var text = request.method() + " " + url + " HTTP/1.1\r\n";
        var headers = request.headers();
        for (var i = 0; i < headers.size(); i++) {
          text += headers.name(i) + ": " + headers.value(i) + "\r\n";
        }
        text += "\r\n";

        var fos = FileOutputStream.$new(File.$new(filePath));
        fos.write(StringCls.$new(text).getBytes());

        var body = request.body();
        if (body !== null) {
          var casted = Java.cast(body, BodyClass);
          var arr = bodyField.get(casted);
          var off = offsetField.getInt(casted);
          var cnt = countField.getInt(casted);
          var bytes = Java.array("byte", new Array(cnt).fill(0));
          for (var i = 0; i < cnt; i++) {
            bytes[i] = ArrayCls.getByte(arr, off + i);
          }
          fos.write(bytes);
        }

        fos.close();
        console.log("[*] Captured: " + url);
        console.log("[*] Saved to: " + filePath);
      } catch (e) {
        console.log("[!] " + e);
      }

      return response;
    };
  });

  console.log("[*] Ready — start playback");
});
