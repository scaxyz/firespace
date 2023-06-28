# Firespace
(Work)Spaces/Profiles for firejail

## Usage
### Example contig
`config.yaml`
```yaml
spaces:
  working:
    home: "/some/path/working"

  private:
    home: ""
    allow_empty_home: true
    env:
      PROXY: "socks5://..."

  nospace:
    home: ""
    allow_empty_home: true
    no_private: true

programms:
  firefox:
    flags:
      - "--no-remote"
      - "https://duckduckgo.com/?q=firejail"
    spaces:
      private:
        overwrites:
          flags: true
        flags:
          - "--no-remote"
          - "--private-window"
          - "https://duckduckgo.com/?q=firejail"
```
### Cli usage

`firespace working firefox` will be executed as => `/usr/bin/firejail --private=/some/path/working firefox --no-remote https://duckduckgo.com/?q=firejail`

`firespace private firefox` will be executed as => `/usr/bin/firejail --private firefox --no-remote --private-window https://duckduckgo.com/?q=firejail`

`firespace nospace firefox` will be executed as => `/usr/bin/firejail firefox --no-remote https://duckduckgo.com/?q=firejail`