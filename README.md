# genstruct

![build](https://github.com/fifsky/genstruct/workflows/build/badge.svg)

Golang struct generator from mysql schema

[![asciicast](https://asciinema.org/a/X5sk7TqrTTjF8AhN764K0Fc6m.svg)](https://asciinema.org/a/X5sk7TqrTTjF8AhN764K0Fc6m)

## Install

```
go get github.com/suncle1993/genstruct
```

## Usage

```
genstruct -h 127.0.0.1 -u root -P 123456 -p 3306
```

* `-h` default `localhost`
* `-u` default `root`
* `-p` default `3306`
