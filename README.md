# ghrel

`ghrel` concurrently downloads and verifies assets of the latest release from a GitHub repository.

```sh
❯ ghrel jreisinger/ghrel
downloaded 6 file(s)
verified 5 file(s)

❯ ghrel -p '*linux*amd64*' brave/brave-browser
downloaded 2 file(s)
verified 1 file(s)
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture. Or `go install github.com/jreisinger/ghrel@latest`.

