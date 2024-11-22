# bitflip

Flip bits in files.

## tldr

If you want to flip bits in files, this utility has got you covered.
Flip specific bits, randomly flip a single bit, or spray hell all around..!
Whatever rocks your boat, as long as you need some bits flipped.

## Usage

### Example usage

```bash
$ echo hello world > /tmp/hello.world
$ cat < /tmp/hello.world
hello world
$ bitflip random  /tmp/hello.world
2020/04/09 17:40:55 flipping 2th bit of byte 6 in file "/tmp/hello.world"
$ cat < /tmp/hello.world
hello sorld
```

From here, sky is the limit!

Why not flip 0.1% of the bits in your MySQL DB file?

```sh
bitflip spray percent:0.1 /var/lib/mysql/mydb.idb
```

## installation

### Debian/Ubuntu

```bash
wget https://github.com/aybabtme/bitflip/releases/download/v0.2.2/bitflip_0.2.2_linux_amd64.deb
dpkg -i bitflip_0.2.2_linux_amd64.deb
```

### darwin

```bash
brew install aybabtme/homebrew-tap/bitflip
```

## Goreleaser

Builds are done using goreleaser. A GitHub workflow will take care of this,
that being said, you can run goreleaser as follows:

> NB. Builds will land in the `dist/` directory.

### Update or install goreleaser

```sh
go install github.com/goreleaser/goreleaser/v2@latest
```

### Create a snapshot

E.g. if you haven't a tag of your current repository state.

```sh
goreleaser build --snapshot --single-target --clean -f .goreleaser.yml
```

### Create a release

```sh
goreleaser release --skip=publish --clean -f .goreleaser.yml
```
