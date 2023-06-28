# Firespace
(Work)Spaces/Profiles for firejail

## Usage
### Example config
`config.yaml`  (use `firespace --help` to see the default location)
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

## Installation
`go install github.com/scaxyz/firespace/cli/firespace@latest`

## Templating
> TODO: improve readme
- go templating
- see `Firespacecontext.go`#TemplateContext
- e. g. 
  ```yaml
  global:
    env:
      _proxy_host: "some-server.mullvad.net"
      _proxy_port: "1080"
      PROXY: "socks5://{{.Space.Env._proxy_host}}:{{.Space.Env._proxy_port}}"
      HTTP_PROXY: "socks5://{{.Space.Env._proxy_host}}:{{.Space.Env._proxy_port}}"
      HTTP_PROXY: "socks5://{{.Space.Env._proxy_host}}:{{.Space.Env._proxy_port}}"
      
  ```

## Development
```bash
git clone ...
cd ...
make install_hook # installed hook will run go tests before commiting #TODO maybe replace with github workflow
```