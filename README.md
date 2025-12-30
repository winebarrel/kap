# kap

[![CI](https://github.com/winebarrel/kap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/kap/actions/workflows/ci.yml)

Simple HTTP proxy with Auth key.

## Usage

```
Usage: kap --backend=BACKEND --port=UINT --key=STRING --secret=SECRET,... [flags]

Flags:
  -h, --help                 Show help.
  -b, --backend=BACKEND      Backend URL ($KAP_BACKEND).
  -p, --port=UINT            Listening port ($KAP_PORT).
  -k, --key=STRING           Auth key name ($KAP_KEY).
  -s, --secret=SECRET,...    Auth secret value ($KAP_SECRET).
      --version
```

```sh
$ export KAP_SECRET=my-secret
$ go run ./cmd/kap -p 8080 -b https://example.com -k my-key
```

```sh
$ curl -H 'my-key: xxx' localhost:8080
forbidden
$ curl -H 'my-key: my-secret' localhost:8080
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
...
```

## Use hashed secret

```
$ export KAP_SECRET=$(htpasswd -nbB '' my-secret | cut -b 2-)
$ go run ./cmd/kap -p 8080 -b https://example.com -k my-key
```

### Run with Docker

```
% cat kap.env
KAP_PORT=8080
KAP_KEY=my-key
KAP_SECRET=my-secret
KAP_BACKEND=https://example.com
$ docker run --rm -p 8080:8080 --env-file kap.env ghcr.io/winebarrel/kap
```
