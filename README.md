[![test](https://github.com/joelee2012/nacosctl/actions/workflows/test.yml/badge.svg)](https://github.com/joelee2012/nacosctl/actions/workflows/test.yml)
[![goreleaser](https://github.com/joelee2012/nacosctl/actions/workflows/release.yml/badge.svg)](https://github.com/joelee2012/nacosctl/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/joelee2012/nacosctl/graph/badge.svg?token=PY470EX7J6)](https://codecov.io/gh/joelee2012/nacosctl)
# nctl
Command line tools for [Nacos](https://nacos.io/) version must be compatible with [api v1](https://nacos.io/docs/v1/open-api/?spm=5238cd80.2ef5001f.0.0.3f613b7cibLcyN)
  - v2.1.x
  - v2.2.x
  - v2.3.x
  - v2.4.x
  - v2.5.x
# Usage

```sh
Command line tools for Nacos

Usage:
  nctl [command]

Available Commands:
  apply       Apply configuration file to nacos
  completion  Generate the autocompletion script for the specified shell
  config      Manage nacos instance config
  create      Create one resource
  delete      Delete one or many resources
  get         Display one or many resources
  help        Help about any command
  version     Print the version number

Flags:
  -h, --help             help for nctl
  -s, --setting string   config file (default is $HOME/.nacos.yaml)

Use "nctl [command] --help" for more information about a command.
```

# Setting

default setting file path is `$HOME/.nacos.yaml`

```yaml
servers:
  test: # server name
    url: http://127.0.0.1:8848/nacos # nacos url with context path
    user: "nacos" # username
    password: "password" # password
context: test # current context name
```
