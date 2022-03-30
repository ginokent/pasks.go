# pasks.go - Parallel Tasks

[![pkg](https://pkg.go.dev/badge/github.com/newtstat/pasks.go)](https://pkg.go.dev/github.com/newtstat/pasks.go)
[![goreportcard](https://goreportcard.com/badge/github.com/newtstat/pasks.go)](https://goreportcard.com/report/github.com/newtstat/pasks.go)
[![workflow](https://github.com/newtstat/pasks.go/workflows/CI/badge.svg)](https://github.com/newtstat/pasks.go/tree/main)
[![codecov](https://codecov.io/gh/newtstat/pasks.go/branch/main/graph/badge.svg?token=cTMvSwlv60)](https://codecov.io/gh/newtstat/pasks.go)
[![sourcegraph](https://sourcegraph.com/github.com/newtstat/pasks.go/-/badge.svg)](https://sourcegraph.com/github.com/newtstat/pasks.go)

## HOW TO USE

See [cmd/example/main.go](/cmd/example/main.go) for details.  

The actual behavior can be seen below:

```console
$ go run cmd/example/main.go -t 10 -p 2
2022/03/31 01:48:14 result: ok=true, task=2, current/total=1/10
2022/03/31 01:48:14 result: ok=true, task=1, current/total=2/10
2022/03/31 01:48:15 result: ok=true, task=4, current/total=3/10
2022/03/31 01:48:15 result: ok=true, task=3, current/total=4/10
2022/03/31 01:48:16 result: ok=true, task=6, current/total=5/10
2022/03/31 01:48:16 result: ok=true, task=5, current/total=6/10
2022/03/31 01:48:17 result: ok=true, task=7, current/total=7/10
2022/03/31 01:48:17 result: ok=true, task=8, current/total=8/10
2022/03/31 01:48:18 result: ok=true, task=9, current/total=9/10
2022/03/31 01:48:18 result: ok=true, task=10, current/total=10/10
2022/03/31 01:48:18 end
```
