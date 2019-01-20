# zcmd: mitZ's CoMmanD utilities for personal use

[![CircleCI](https://circleci.com/gh/mitsutaka/zcmd.svg?style=svg)](https://circleci.com/gh/mitsutaka/zcmd)
[![GoDoc](https://godoc.org/github.com/mitsutaka/zcmd?status.svg)][godoc]
[![Go Report Card](https://goreportcard.com/badge/github.com/mitsutaka/zcmd)](https://goreportcard.com/report/github.com/mitsutaka/zcmd)

## Usage

```console
z
```

### `nas` command

Operation command for home NAS.

- `z nas pull`: Pull media files from the remote server.
- `z nas push`: Push media files to the remote server.

### `backup` command

Operation command for backup

- `z backup`: Run backup to the remote server.

### `repos` command

Operation command for checked out git repositories.

- `z repos update`: Fetch and Checkout git repositories in parallel.

### Misc

To load bash completion scripts, run:

```console
. <(z completion)
```

## License

MIT
