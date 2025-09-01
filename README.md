# kap

[![CI](https://github.com/winebarrel/kap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/kap/actions/workflows/ci.yml)

Simple HTTP proxy with Auth key.

## Usage

```
Usage: kap --backend=BACKEND --port=UINT --key=STRING --secret=STRING [flags]

Flags:
  -h, --help               Show help.
  -b, --backend=BACKEND    Backend URL ($GAP_BACKEND).
  -p, --port=UINT          Listening port ($GAP_PORT).
  -k, --key=STRING         Auth key name ($GAP_KEY).
  -s, --secret=STRING      Auth secret value ($GAP_SECRET).
      --version
```

```sh
$ export GAP_SECRET=my-secret
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
