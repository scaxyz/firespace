# Firespace
Workspace/Profiles for firejail

e.g.  

`config.yaml`
```yaml
spaces:
  working:
    home: "/some/path/working"


programms:
  firefox:
    flags:
      - "--no-remote"
      - "https://duckduckgo.com/?q=firejail"
```

`firespace working firefox` => `/usr/bin/firejail --private=/some/path/working firefox --no-remote https://duckduckgo.com/?q=firejail`