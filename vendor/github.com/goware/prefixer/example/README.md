# Prefixer example

See [prefix.go](./prefix.go).

## Build "prefix" binary

```bash
$ git clone https://github.com/goware/prefixer.git
$ cd prefixer/example
$ go build -o prefix
```

## Usage

Create an email reply (`"> "` prefix) from any text easily:

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
