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

Configuration example `$HOME/.z.yaml`:

```yaml
nas:
  # Config for pulling
  pull:
    # Source URL for pulling
    source: rsync://URL
    # Destination paths
    destinations:
      # This name uses when command runs with particular path
      - name: tv
      # Destination directory
        path: /mnt/nas/tv
      - name: movie
        path: /mnt/nas/movie
        # Exclude pattern
        excludes:
          - xxxx
          - yyyy
  # Config for pushing
  push:
    # Source directories
    sources:
      - name: dvd
        path: /mnt/nas/dvd
        excludes:
          - zzzz
    # Destination URL for pushing
    destination: rsync://URL
```

### `backup` command

Operation command for backup

- `z backup`: Run backup to the remote server.

Configuration example `$HOME/.z.yaml`:

```yaml
backup:
  # Backup URL
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
