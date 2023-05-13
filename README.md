# ghrel

```
$ ghrel -h
Download and verify assets (files) of the latest release from a GitHub repository.

ghrel [flags] <owner>/<repo>
  -l	just list assets, don't download them
  -p pattern
    	assets matching shell pattern (doesn't apply to checksum files)
  -v	be verbose
```

To use ghrel download a [binary](https://github.com/jreisinger/ghrel/releases) for your system and architecture. Or, if you have [Go installed](https://go.dev/doc/install):

```
$ go install github.com/jreisinger/ghrel@latest
```
