# Traefik 2 check request headers middleware plugin

This plugin checks the incoming request for specific headers and their values to be present and matching the configuration. If the request does not validate against the configured headers, the middleware will return a 403 Forbidden status code. This is can also be used to check client certificate information in combination with the [PassTLSClientCert](https://doc.traefik.io/traefik/middlewares/passtlsclientcert/) Traefik middleware, details can be found at the end of this document.

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

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: Path(`/whoami`)
      service: service-whoami
      entryPoints:
        - http
      middlewares:
        - checkheaders

  services:
    service-whoami:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    checkheaders:
      plugin:
        dev:
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
            - header:
              name: "HEADER_4"
              matchtype: one
              values:
                - "VALUE_\d"
              required: true
              regex: true
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

| Setting   | Allowed values | Description |
| :--       | :--            | :--         |
| name      | string         | Name of the request header |
| matchtype | one, all, none | Match on all values, one of the values specified or none of the values. The value 'all' is only allowed in combination with the 'contains' and 'regex' setting. |
| values    | []string       | A list of allowed values which are matched against the request header value |
| contains  | boolean        | If set to true (default false), the request is allowed if the request header value contains the value specified in the configuration |
| regex     | boolean        | If set to true (default false), the match is done using a regular expression. The value of the request header is matched against the value specified in the configuration. via the [regexp](https://pkg.go.dev/regexp) package |
| required  | boolean        | If set to false (default true), the request is allowed if the header is absent or the value is empty |
| urldecode | boolean        | If set to true (default false), the value of the request header will be URL decoded before further processing with the plugin. This is useful when using this plugin with the [PassTLSClientCert](https://doc.traefik.io/traefik/middlewares/passtlsclientcert/) middleware that Traefik offers. |
| debug     | boolean        | If set to true (default false), the request headers, values and validation will be printed to the console |

#

## Example 1 config

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
          - header:
              name: "HEADER_4"
              matchtype: one
              values:
                - "VALUE_\d"
              regex: true
              required: true
```

## Example 2 config

You can also use this plugin to check on client certificate fields when using mTLS configuration. The [PassTLSClientCert](https://doc.traefik.io/traefik/middlewares/passtlsclientcert/) Traefik middleware adds the client certificate information to the request header `X-Forwarded-Tls-Client-Cert-Info` in a URL encoded format. Using this plugin as second middleware for route, you can verify the client certificate fields.

Example client certificate request header:

```http
X-Forwarded-Tls-Client-Cert-Info: Subject="C=US,ST=Ohio,L=Akron,O=Google,CN=server0.google.com";Issuer="DC=us,DC=google.com,DC=com,CN=GoogleRootCA";NB="1687386830";NA="1750458830";SAN="server0.google.com"
```

You could configure the plugin to check for the `CN` and the `DC` fields:

```yaml
middlewares:
  my-checkheadersplugin:
    plugin:
      checkheadersplugin:
        headers:
          - header:
            name: "X-Forwarded-Tls-Client-Cert-Info"
            matchtype: all
            values:
              - "CN=server0.google.com"
              - "DC=google.com"
            contains: true
            required: true
            urldecode: true
```

## Example 3 config

This plugin give you also the possibility to validate header via a regular expression. This can be useful when you want to validate a header value against a pattern. For example, you want to validate a JWT token in the `Authorization` header. The JWT token has a specific format and you can validate this with a regular expression.

```yaml
middlewares:
  my-checkheadersplugin:
    plugin:
      checkheadersplugin:
        headers:
          - header:
            name: "Authorization"
            matchtype: one
            values:
              - "^Bearer .*"
            regex: true
```

## Example 4 config

You can also use this plugin to check if header has a certain value that is not allowed. This way you can allow every value except a the provide ones, acting as blacklist. For example, you want to block requests that have `Content-Language` header that are set to `de-DE` or `de-AT`. You can use 'none' in combination with matchtype 'regex' or 'contains'.

```yaml
middlewares:
  my-checkheadersplugin:
    plugin:
      checkheadersplugin:
        headers:
          - header:
            name: "Content-Language"
            matchtype: none
            values:
              - "de-DE"
              - "de-AT"
            required: true
````
