# ghrel

`ghrel` lists or downloads assets (files) of the latest release from a GitHub repository.

```sh
$ ghrel -l -p '*linux*' jreisinger/ghrel
Asset                           Updated     Size     Download count
-----                           -------     ----     --------------
ghrel_0.6.0_linux_386.tar.gz    2023-03-21  2031247  2
ghrel_0.6.0_linux_amd64.tar.gz  2023-03-21  2144185  2
ghrel_0.6.0_linux_arm64.tar.gz  2023-03-21  1974932  2
ghrel_0.6.0_linux_armv6.tar.gz  2023-03-21  2009336  2

$ ghrel -p '*linux*' jreisinger/ghrel
downloaded 5 file(s)
verified 4 file(s)
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture. Or `go install github.com/jreisinger/ghrel@latest`.
