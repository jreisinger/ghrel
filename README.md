# ghrel

`ghrel` concurrently downloads and verifies assets (files) of the latest release from a GitHub repository.

```sh
# donwload all
❯ ghrel jreisinger/ghrel
downloaded 11 file(s)
verified 10 file(s)

# donwload those matching shell pattern
❯ ghrel -p '*linux*amd64*' brave/brave-browser
downloaded 2 file(s)
verified 1 file(s)

# list those matching shell pattern
❯ ghrel -l -p '*linux*amd64*' brave/brave-browser
brave-browser-1.45.116-linux-amd64.zip
brave-browser-1.45.116-linux-amd64.zip.sha256
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture. Or `go install github.com/jreisinger/ghrel@latest`.
