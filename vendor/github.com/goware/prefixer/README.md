# Prefixer
[Golang](http://golang.org/)'s [io.Reader](http://golang.org/pkg/io/#Reader) wrapper prepending every line with a given string.

[![GoDoc](https://godoc.org/github.com/goware/prefixer?status.png)](https://godoc.org/github.com/goware/prefixer)
[![Travis](https://travis-ci.org/goware/prefixer.svg?branch=master)](https://travis-ci.org/goware/prefixer)


## Use cases
1. Logger that prefixes every line with a timestamp etc.
    ```bash
    16:54:49 My awesome server | Creating etcd client pointing to http://localhost:4001
    16:54:49 My awesome server | Listening on http://localhost:8080
    16:54:49 My awesome server | [restful/swagger] listing is available at 127.0.0.1:8080/swaggerapi
    ```

2. SSH multiplexer prepending output from multiple servers with a hostname
    ```bash
    host1.example.com | SUCCESS
    host2.example.com | SUCCESS
    host3.example.com | -bash: cd: workdir: No such file or directory
    host4.example.com | SUCCESS
    ```

3. Create an email reply (`"> "` prefix) from any text easily.
    ```bash
    $ ./prefix
    Dear John,               
    did you know that https://github.com/goware/prefixer is a golang pkg
    that prefixes every line with a given string and accepts any io.Reader?

    Cheers,
    - Jane
    ^D     
    > Dear John,               
    > did you know that https://github.com/goware/prefixer is a golang pkg
    > that prefixes every line with a given string and accepts any io.Reader?
    >
    > Cheers,
    > - Jane
    ```

## Example

See the ["Prefix Line Reader" example](./example).

```go
package main

import (
    "io/ioutil"
    "os"

    "github.com/goware/prefixer"
)

func main() {
    // Prefixer accepts anything that implements io.Reader interface
    prefixReader := prefixer.New(os.Stdin, "> ")

    // Read all prefixed lines from STDIN into a buffer
    buffer, _ := ioutil.ReadAll(prefixReader)

    // Write buffer to STDOUT
    os.Stdout.Write(buffer)
}
```

## License
Prefixer is licensed under the [MIT License](./LICENSE).
