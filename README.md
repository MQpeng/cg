# cg

A code generator

## Quick Start

1. Add `template`
```shell
# add a template to user path
cg add example
# as new name
cg add example test
```
2. Use `Template`
```shell
# use template to generate
cg g example --name Tom --age 12
cg g test --name Job --age 44
```

## Use `git`

1. Add `template` repo
```shell
# clone template repo
cg clone https://github.com/MQpeng/cg-templates
```
2. fetch `template` repo
```shell
# pull template repo
cg pull
```