# who dat?

Tiny Go web server that prints OS information and HTTP request to output.  I'll be evolving it to use gRPC for HTTP/2 and gRPC-gateway for HTTP 1.1.

```sh
$ docker run -d -P --name blah lukesiler/whodat
$ docker inspect --format '{{ .NetworkSettings.Ports }}' blah
map[80/tcp:[{0.0.0.0 32769}]]
$ curl "http://0.0.0.0:32769"
Hostname: 3cce2c28be76
IP: 127.0.0.1
IP: 172.17.0.2
GET / HTTP/1.1
Host: 0.0.0.0:32769
User-Agent: curl/7.54.0
Accept: */*

```
