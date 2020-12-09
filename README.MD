# Traefik 2 check request headers middleware plugin

This plugin checks the incoming request for specific headers and their values to be present and matching the configuration. If the request does not validate against the configured headers, the middleware will return a 403 Forbidden status code.

## Dev `traefik.yml` configuration file for traefik

```yaml
# Static configuration
api:
  dashboard: true
  insecure: true

pilot:
  token: <your-token-here>

experimental:
  devPlugin:
    goPath: /home/user/go
    moduleName: github.com/dkijkuit/checkheadersplugin

entryPoints:
  http:
    address: ":4000"
    forwardedHeaders:
      insecure: true

providers:
  file:
    filename: dynamic-dev-config.yaml

```

## Launch Traefik using dev config (config of plugin can be found in dynamic-dev-config.yaml)
```bash
$ docker run --rm -d -p 5000:80 containous/whoami
```

## Test using cURL
```bash
curl --location --insecure --request GET "http://localhost:4000/whoami" --header "HEADER_1: VALUE_99" --header "HEADER_2: VALUE_2" --header "HEADER_3: VALUE_3" --header "HEADER_4: VALUE_X_and_VALUE_4_and_VALUE_5_AND_6" 
```

Should return a 200 showing details about the request.

#
## Configuration documentation

Supported configurations per header

| Setting | Allowed values | Description |
|---|---|---|
| name | string | Name of the request header |
| matchtype | one, all | Match on all values or one of the values specified. The value 'all' is only allowed in combination with the 'contains' setting.|
| values | []string | A list of allowed values which are matched against the request header value|
| contains | boolean | If set to true (default false), the request is allowed if the rtequest header value contains the value specified in the configuration |
| required | boolean | If set to false (default true), the request is allowed if the header is absent or the value is empty|

#
## Example config
```yaml
middlewares:
  my-checkheadersplugin:
    plugin:
      checkheadersplugin:
        headers:
          - header:
            name: "HEADER_1"
            matchtype: one
            values: 
              - "VALUE_1"
              - "VALUE_99"
          - header:
            name: "HEADER_2"
            matchtype: one
            values: 
              - "VALUE_2"
          - header:
            name: "HEADER_3"
            matchtype: one
            values: 
              - "VALUE_3"
            required: false
          - header:
            name: "HEADER_4"
            matchtype: all
            values: 
              - "LUE_4"
              - "VALUE_5"
            contains: true
            required: true
```