# zcmd: mitZ's CoMmanD line collections

[![CircleCI](https://circleci.com/gh/mitsutaka/zcmd.svg?style=svg)](https://circleci.com/gh/mitsutaka/zcmd)
[![GoDoc](https://godoc.org/github.com/mitsutaka/zcmd?status.svg)](https://godoc.org/github.com/mitsutaka/zcmd)
[![Go Report Card](https://goreportcard.com/badge/github.com/mitsutaka/zcmd)](https://goreportcard.com/report/github.com/mitsutaka/zcmd)

## Installation

```console
go get github.com/mitsutaka/zcmd/pkg/z
```

## Usage

```console
z
```

### `sync` command

`rsync` wrapper command for servers. It plans replace it with Go native sync library.

- `z sync pull`: Pull files from the remote server in parallel.
- `z sync push`: Push files to the remote server in parallel.

Configuration example `$HOME/.z.yaml`:

```yaml
sync:
  # Config for pulling
  pull:
    # It uses when command runs with particular path
    - name: movie
      # Source URL for pulling
      source: RSYNC_URL
      # Destination directory
      destination: /mnt/nas/movie
    - name: picture
      source: RSYNC_URL
      destination: /mnt/nas/picture
      # Exclude pattern
      excludes:
        - xxxx
        - yyyy
  # Config for pushing
  push:
    - name: music
      source: /mnt/nas/music
      destination: RSYNC_URL
      excludes:
        - zzzz
```

### `backup` command

`rsync` wrapper command for backup.

- `z backup`: Run backup to the remote server in parallel.

Configuration example `$HOME/.z.yaml`:

```yaml
backup:
  # Backup URLs
  destinations:
    - rsync://BACKUP_URL
  # Include backup paths
  includes:
    - /
    - /boot
    - /home
 # Exclude paths and pattern
  excludes:
    - .cache/
    - /dev
    - /media
    - /misc
    - /mnt
    - /proc
    - /run
    - /sys
    - /var/cache
```

### `repos` command

Operation command for checked out git repositories.

- `z repos update`: Fetch and Checkout git repositories in parallel.

Configuration example `$HOME/.z.yaml`:

```yaml
repos:
  # Root directory of the git repositories
  root: /home/mitz/git/repos
```

### Misc

To load bash completion scripts, run:

```console
. <(z completion)
```

## License

MIT
