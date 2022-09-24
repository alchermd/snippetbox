# SnippetBox

This repository contains my code-along with [Alex Edward's Let's Go book](https://lets-go.alexedwards.net/).

## Setup

```console
$ cd tls/
$ go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
$ cd ..
$ go run ./cmd/web
```