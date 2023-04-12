# ghrel

`ghrel` downloads and verifies, or just lists, assets (files) of the latest release from a GitHub repository.

```sh
> ghrel -p '*linux*' jreisinger/ghrel
downloaded	4 (+ 1 checksums files)
verified	4

> ghrel -l -p '*linux*' jreisinger/ghrel
Asset                           Checksums file  Updated     Size     Download count
-----                           --------------  -------     ----     --------------
checksums.txt                   true            2023-04-07  976      15
ghrel_0.6.2_linux_386.tar.gz    false           2023-04-07  2031455  13
ghrel_0.6.2_linux_amd64.tar.gz  false           2023-04-07  2144218  28
ghrel_0.6.2_linux_arm64.tar.gz  false           2023-04-07  1974958  14
ghrel_0.6.2_linux_armv6.tar.gz  false           2023-04-07  2009334  12
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture or `go install github.com/jreisinger/ghrel@latest`.
