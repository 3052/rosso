def request(flow):
    if flow.request.method == "POST" and "requestVamAPI=true" in flow.request.url:
        flow.kill()
