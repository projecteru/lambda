Lambda
=============
[![CircleCI](https://circleci.com/gh/projecteru2/lambda/tree/master.svg?style=shield)](https://circleci.com/gh/projecteru2/lambda/tree/master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a91da2853c4c4dc3b155f3397778f47e)](https://www.codacy.com/app/CMGS/lambda?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=projecteru2/lambda&amp;utm_campaign=Badge_Grade)

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
