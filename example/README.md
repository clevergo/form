# Form Decoder Example

```shell
$ cd example
$ go run main.go
```

```shell
$ curl -XPOST -d "username=foo&password=bar"  http://localhost:12345/login
username: foo, password: bar

$ curl -XPOST -H "Content-Type: application/json" -d '{"username":"foo", "password": "bar"}' http://localhost:12345/login
username: foo, password: bar

$ curl -XPOST -H "Content-Type: application/xml" -d '<xml><username>foo</username><password>bar</password></xml>' http://localhost:12345/login
username: foo, password: bar

$ curl -XPOST -F "username=foo" -F "password=bar" http://localhost:12345/login
username: foo, password: bar
```