Eru Lambda
=============

Run anything in container on Eru

Params
======

Must

* config
* name
* command

Optional

* env
* volumes

Get from config (can override)

* pod
* network
* image
* cpu
* mem
* timeout
* stdin

Default

* count
* admin
* debug
* help
* version

Admin
======

if admin is `True`

pod and volumes will be rewrited.

their values get from config file.

DEV
======

#### Test

```make test```

#### Build

```make build```

#### RPM

```./make-rpm```

To make rpm, you should install [fpm](https://github.com/jordansissel/fpm) first.