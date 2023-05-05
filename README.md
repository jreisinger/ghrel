# ghrel

`ghrel` lists, downloads and verifies assets (files) of the latest release from a GitHub repository:

```sh
$ ghrel -v -p '*linux*' -l jreisinger/ghrel
Asset                           Updated     Size     Download count
-----                           -------     ----     --------------
ghrel_0.8.0_linux_386.tar.gz    2023-04-12  2031657  7
ghrel_0.8.0_linux_amd64.tar.gz  2023-04-12  2144060  42
ghrel_0.8.0_linux_arm64.tar.gz  2023-04-12  1974168  8
ghrel_0.8.0_linux_armv6.tar.gz  2023-04-12  2008815  7

$ ghrel -v -p '*linux*' jreisinger/ghrel
downloaded	4 + 1 checksum file(s)
verified	4
removed checksum file(s)
```

To use ghrel, download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture or if you have [Go installed](https://go.dev/doc/install):

```sh
$ go install github.com/jreisinger/ghrel@latest
```
