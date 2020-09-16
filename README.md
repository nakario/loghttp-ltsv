loghttp-ltsv
==========

Package loghttpltsv enables a `http.Client` to do logging with LTSV format.
LTSV-formatted logs can be easily profiled with tools such as [alp](https://github.com/tkuchiki/alp).

# Usage

```sh
go get -u github.com/nakario/loghttp-ltsv
```

Use a custom `http.Client` to write logs
```go
package main

import (
	"log"
	"net/http"
	"os"

	lhl "github.com/nakario/loghttp-ltsv"
)

func main() {
	f, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cli := &http.Client{
		Transport: lhl.NewTransport(f),
	}

	resp, err := cli.Get(os.Args[1])
	...
}
```

# Acknowledgements

This package is inspired by https://github.com/motemen/go-loghttp
