# go-packetflagon
![alt text](https://pbs.twimg.com/media/CegH8M9W4AEXeHz.png "Example image of the http interface")

An application that serves customised Proxy Auto Configuration files for your browser to help bypass Internet censorship.

## How it Works
URLs added to the PAC file will be sent to a SOCKS proxy listening on ```localhost:9050``` / ```localhost:9051``` whilst any other URLs will use your normal Internet connection.

We recommend you use the [Tor client or Tor Browser Bundle.](https://www.torproject.org/index.html.en) as the local SOCKS proxy *(it will do this automatically)* or follow our [guide on creating an SSH based SOCKS5 proxy](https://immunicity.org/how-to/create-socks5-proxy/).


## Getting Started
go-packetflagon currently requires a config file passed with the -conf argument defining the location of the URL database, a listen port etc in order to start;

```javascript
{
        "dbpath":"/tmp/pacs.db",
        "listenport": 8080,
        "debug":true,
        "tls_enabled":false
}```
