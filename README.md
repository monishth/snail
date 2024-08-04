# snail

`snail` is a simple forward proxy written in go. It is __not__ as slow as a snail.

Features:

- HTTP server
- Minimal configuration
- `Basic` proxy authentication
  - Supports `.htpasswd`
- Supports HTTPS requests (however is not a HTTPS server)

## Installation

#### From `source`

> 1. Clone the snail repository:
>
> ```console
> $ git clone https://github.com/monishth/snail
> ```
>
> 2. Change to the project directory:
>
> ```console
> $ cd snail
> ```
>
> 3. Install the dependencies:
>
> ```console
> $ go mod download
> ```
>
> 4. Build:
>
> ```console
> $ go build -o snail
> ```

## Usage

```console
./snail
```

### Help

```console

A quick, simple, faster-than-snail forward proxy server

Usage:
  snail [flags]

Flags:
  -a, --auth AuthProvider   auth provider - must be "htpasswd", "simple" or "none" (default none)
  -f, --filename string     filename for htpasswd auth
  -h, --help                help for snail
  -p, --port int            server port (default 8080)
  -u, --userpass string     user:pass pair for simple
```
