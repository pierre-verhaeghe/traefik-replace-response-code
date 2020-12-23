# traefik-replace-response-code
Traefik plugin used to replace a http response code. Can be useful for example to mask an internal error code with another status code.


### Warning : still under test. Not yet ready for production.


### Configuration



```toml
# Static configuration
[pilot]
    token = "xxxx"

[experimental.plugins.replace]
  modulename = "github.com/pierre-verhaeghe/traefik-replace-response-code"
  version = "v0.2.0"
#
```
removeBody optional flag can be use to remove response body. By default removeBody is set to `false`.
```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    my-plugin:
      plugin:
        replace:
          inputCode: 429
          outputCode: 200
          removeBody: true
```