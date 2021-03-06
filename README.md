# What is this?

This server implements [lego's httpreq _RAW_
API](https://go-acme.github.io/lego/dns/httpreq/), and uses lego itself to then
create or remove ACME challenges from one of [its supported DNS
providers](https://go-acme.github.io/lego/dns/).

This was built to work around the fact that Traefik 1.7 was built against a now
old version of lego, which means you can't use many of the DNS providers that
have been added to lego in the meantime.

# Security

**DO NOT EXPOSE THIS SERVER TO THE INTERNET**

This server provides no security to speak of, and as such, should not be
exposed to the internet. If you need authentication, put a reverse proxy in
front of it that handles that for you.

# How to use?

The server is configured entirely through environment variables.

## `LISTEN_ADDRESS`

The address the server will listen on. See [go's net.Dial
documentation](https://golang.org/pkg/net/#Dial).

Defaults to: `:8080`.

Example values:

- `:8080`: listen on all IP's, on port 8080
- `127.0.0.1:3000`: listen on localhost, on port 3000
- `[::1]:3000`: listen on localhost's IPv6 address, on port 3000

## `DNS_PROVIDER`

The name of the DNS provider to use when handling requests. See lego's
documentation for the list of supported providers.

## Other environment variables

Just in case you didn't know: lego's providers are also configured through
environment variables. Odds are you'll need to add credentials somehow. Refer
to [lego's documentation](https://go-acme.github.io/lego/dns/) to find out what
variables apply to your provider.

# Traefik 1.7 integration

## Traefik configuration file

```toml
[acme.dnschallenge]
  provider = "httpreq"
```

## Traefik environment variables

- `HTTPREQ_ENDPOINT=<URL to lego-httpreq-server>`
- `HTTPREQ_MODE=RAW`: _lego-httpreq-server_ does **not** support the _default_ mode
