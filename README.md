# bitflip

Flip bits in files.


# example

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
```
$ bitflip spray percent:0.1 /var/lib/mysql/mydb.idb
```


## tldr

If you want to flip bits in files, this utility has got you covered. Flip specific bits, randomly flip a single bit, or spray hell all around..! Whatever rocks your boat, as long as you need some bits flipped.

## installation

### Debian/Ubuntu

```bash
wget https://github.com/aybabtme/bitflip/releases/download/v0.2.0/bitflip_0.2.0_linux_amd64.deb
dpkg -i bitflip_0.2.0_linux_amd64.deb
```

### darwin

```bash
brew install aybabtme/homebrew-tap/bitflip
```
