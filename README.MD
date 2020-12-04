# Traefik 2 check request headers middleware plugin

This plugin checks the incoming request for specific headers and their values to be present and matching the configuration. If the request does not validate against the configured headers, the middleware will return a 403 Forbidden status code.