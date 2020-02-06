# httping simple healtcheck in golang

```bash
$ httping -url http://server -code 200 -code 202 -contains "this is a string on the webpace"
```

## Download

```bash
$ go get github.com/rwese/httping
```