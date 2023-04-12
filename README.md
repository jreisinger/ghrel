# ghrel

`ghrel` downloads and verifies, or just lists, assets (files) of the latest release from a GitHub repository.

```sh
> ghrel -p '*linux*' jreisinger/ghrel
downloaded	4 + 1 checksum file(s)
verified	4

> ghrel -l -p '*linux*' jreisinger/ghrel
Asset                           Checksum file  Updated     Size     Download count
-----                           -------------  -------     ----     --------------
checksums.txt                   true           2023-04-12  976      1
ghrel_0.7.0_linux_386.tar.gz    false          2023-04-12  2028791  1
ghrel_0.7.0_linux_amd64.tar.gz  false          2023-04-12  2142561  1
ghrel_0.7.0_linux_arm64.tar.gz  false          2023-04-12  1972445  1
ghrel_0.7.0_linux_armv6.tar.gz  false          2023-04-12  2006885  1
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture or `go install github.com/jreisinger/ghrel@latest`.
