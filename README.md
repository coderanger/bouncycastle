# BouncyCastle

A small Go app to redirect (bounce) incoming HTTP requests.

## Options

* `--host` - Hostname to put in the URL while redirecting. *(required)*
* `--scheme` - Scheme to put in the URL while redirecting. *(default: https)*
* `--path` - Path to put in the URL while redirecting. *(default: preserve the requested path)*
* `--status` - HTTP status to use for redirecting. *(default: 302)*
* `--port` - Port to listen on. *(default: 8080/8443)*
* `--tls` - Enable HTTPS support.
* `--cert` - Path to the TLS certificate to use.
* `--key` - Path to the TLS key to use.


