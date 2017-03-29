#!/bin/bash

ROOT="`pwd`/build"
BIN="$ROOT/usr/bin"
CONF="$ROOT/etc/"

mkdir -p $BIN
mkdir -p $CONF

mv lambda $BIN
mv lambda.yaml.example $CONF

VERSION=$(cat VERSION)
echo $VERSION rpm build begin

fpm -f -s dir -t rpm -n eru-lambda --epoch 0 -v $VERSION --iteration 1.el7 -C $ROOT -p $PWD --verbose --rpm-auto-add-directories --category 'Development/App' --description 'docker eru lambda executor' --url 'http://gitlab.ricebook.net/platform/lambda-run/' --license 'BSD'  --no-rpm-sign usr etc

rm -rf $ROOT
