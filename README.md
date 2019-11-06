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

Command line flag can append additional rsync flags. Default rsync flags are `-avzP --stats --delete --delete-excluded`.

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
    - name: fuse
      source: /home/mitz/Documents
      destination: LOCAL_PATH
      # Some cases should disable sudo with fuse mounts
      disable_sudo: true
```

### `backup` command

`rsync` wrapper command for backup.

- `z backup`: Run backup to the remote server in parallel.

Command line flag can append additional rsync flags. Default rsync flags are `-avzP --stats --delete`.

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

### `repos-update` command

Run `git clean`, `git checkout master` and `git pull` for checked out git repositories.

Configuration example `$HOME/.z.yaml`:

```yaml
repos:
  # Root directory of the git repositories
  root: /your/root/git/repos
```

### `dotfiles` command

dotfiles manager inspired by <https://github.com/dotphiles/dotphiles>

- `z dotfiles init GITURL`: Initialize dotfiles in local.
- `z dotfiles pull`: Download latest dotfiles and make symbolic links.

```yaml
dotfiles:
  # default is $HOME/.zdotfiles
  dir: /home/mitz/.zdotfiles
  hosts:
    - YOUR_HOSTNAME
  files:
    - bashrc
    - config/sway/config
    - spacemacs
    - ssh
```

### `proxy` command

***NOT IMPLEMENTED YET***

Make multiple ssh port forward at once.

- `z proxy`: Setup ssh port forward.

```yaml
proxy:
  - name: testforward1
    user: ubuntu
    address: remotehost1
    privateKey: ~/.ssh/id_rsa
    forward:
      # Local forwarding
      - type: local
        # default bindAddress is *
        bindAddress: localhost
        bindPort: 13128
        remoteAddress: localhost
        remotePort: 3128
      # Dynamic forwarding for SOCK4, 5
      - type: dynamic
        bindAddress: localhost
        bindPort: 1080
  - name: testforward2
    user: admin
    address: remotehost2
    privateKey: ~/.ssh/id_ecdsa
    port: 10000
    forward:
      # Remote forwarding
      - type: remote
        bindAddress: localhost
        bindPort: 9000
        remoteAddress: localhost
        remotePort: 3000
```

### Misc

To load bash completion scripts, run:

```console
. <(z completion)
```

## License

MIT
