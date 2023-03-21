# ghrel

`ghrel` lists or downloads assets (files) of the latest release from a GitHub repository.

```sh
$ ghrel -l -p '*linux*' jreisinger/ghrel
Asset                           Updated     Size     Download count
-----                           -------     ----     --------------
ghrel_0.6.0_linux_386.tar.gz    2023-03-21  2031247  6
ghrel_0.6.0_linux_amd64.tar.gz  2023-03-21  2144185  7
ghrel_0.6.0_linux_arm64.tar.gz  2023-03-21  1974932  5
ghrel_0.6.0_linux_armv6.tar.gz  2023-03-21  2009336  7

$ ghrel -p '*linux*' jreisinger/ghrel
downloaded	4 (+ 1 checksums file)
verified	4
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture or `go install github.com/jreisinger/ghrel@latest`.
