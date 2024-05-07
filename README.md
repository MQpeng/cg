# cg

a text tool that saves you time and helps your team build new files with consistency.

## Install

1. Download from [release](https://github.com/MQpeng/cg/releases)
2. Install by `npm`

```shell
npm i -g @tonyer/cg
```

```shell
# other package manager
yarn add -g @tonyer/cg
pnpm add -g @tonyer/cg
```

## Quick Start

- Add `template`

```shell
cd examples
# add a template to user path
cg add base
# as new name
cg add base test
```

- Use `Template name`

```shell
# use template to generate
cg g base --name Tom --age 12
cg g test --name Job --age 44
```

- Use `Template path`

```shell
# use template to generate
cg g --template=examples/base --data="{'name': 'Tom', 'age': 12}"
cg g --template=examples/test --data="{'name': 'Tom', 'age': 12}"
```

## Use `git`

- Add `template` repo

```shell
# clone template repo
cg clone https://github.com/MQpeng/cg-templates
```

- fetch `template` repo

```shell
# pull template repo
cg pull
```
