# ghrel

`ghrel` downloads and verifies, or just lists, assets (files) of the latest release from a GitHub repository.

```sh
> ghrel -p '*linux*' jreisinger/ghrel
downloaded	4 + 1 checksum file(s)
verified	4

> ghrel -l -p '*linux*' jreisinger/ghrel
Asset                           Updated     Size     Download count
-----                           -------     ----     --------------
ghrel_0.7.1_linux_386.tar.gz    2023-04-12  2029844  2
ghrel_0.7.1_linux_amd64.tar.gz  2023-04-12  2143213  2
ghrel_0.7.1_linux_arm64.tar.gz  2023-04-12  1973067  2
ghrel_0.7.1_linux_armv6.tar.gz  2023-04-12  2007574  2
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture or `go install github.com/jreisinger/ghrel@latest`.
