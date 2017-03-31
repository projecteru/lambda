Lambda runner
=============

Run anything in container.

Params
======

Must

* config
* name
* command

Optional

* env

Get from config (can override)

* pod
* network
* image
* cpu
* mem
* timeout

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

