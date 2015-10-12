### Introduction

Simple blog engine by golang ,with tango and tidb , cross platform, easy use.

### Quick Start

- Download and extract [binary release](https://github.com/go-xiaohei/pugo/releases) , run `pugo[.exe] server`
- Open browser to visit `http://localhost:9899` to preview default website

### Advanced

- Configuration file is `config.json`
- Database is based on `goleveldb` and stores in `data` directory. **Be careful**
- Theme directory contains user themes and administrator theme(cant delete)
- More documentation in [wiki](https://github.com/go-xiaohei/pugo/wiki)

### Thanks

- [tango](https://github.com/lunny/tango) - middleware-design and pluginable web framework
- [tidb](https://github.com/pingcap/tidb) - sql database based on kv store engine
- [xorm](https://github.com/go-xorm/xorm) - awesome golang orm library
- [codegangsta/cli](https://github.com/codegangsta/cli) - powerful command-line application framework
- [log15](https://gopkg.in/inconshreveable/log15.v2) - fantastic logging library
- [russross/blackfriday](https://github.com/russross/blackfriday) - markdown render engine

Meanwhile , thanks a lot to [lunny](https://github.com/lunny) and [Unknwon](https://github.com/Unknwon) for new features and beta testing