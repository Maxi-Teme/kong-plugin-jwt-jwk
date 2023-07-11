# kong-plugin-jwt-jwk

## Lua

## Go

from go/ directory

```sh
docker build -t kong-jwt-jwk .
```

```sh
docker run --rm --name kong-example \
  --mount type=bind,source="$(pwd)"/kong.yaml,target=/kong/kong.yaml \
  -e "KONG_DATABASE=off" \
  -e "KONG_DECLARATIVE_CONFIG=/kong/kong.yaml" \
  -e "KONG_PLUGINS=bundled,jwt-jwk" \
  -e "KONG_PLUGINSERVER_NAMES=jwt-jwk" \
  -e "KONG_PLUGINSERVER_JWT_JWK_START_CMD=/kong/go-plugins/jwt-jwk" \
  -e "KONG_PLUGINSERVER_JWT_JWK_QUERY_CMD=/kong/go-plugins/jwt-jwk -dump" \
  -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
  -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
  -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
  -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
  -e "KONG_PROXY_LISTEN=0.0.0.0:8000" \
  -e "KONG_ADMIN_LISTEN=0.0.0.0:8001" \
  -e "KONG_LOG_LEVEL=info" \
  -p 8000:8000 \
  -p 127.0.0.1:8001:8001 \
  kong-jwt-jwk
```

sadly no obvious caching available
