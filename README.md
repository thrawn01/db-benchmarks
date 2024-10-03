### Embedded Key/Value Stores
Benchmarking some embedded key value stores.

### Results
```
goos: darwin
goarch: amd64
pkg: github.com/thrawn01/db-benchmarks
cpu: VirtualApple @ 2.50GHz
BenchmarkDatabases
BenchmarkDatabases/BuntDB
BenchmarkDatabases/BuntDB/BuntDB-Set
BenchmarkDatabases/BuntDB/BuntDB-Set-10               387495          3571 ns/op
BenchmarkDatabases/BuntDB/BuntDB-Get
BenchmarkDatabases/BuntDB/BuntDB-Get-10              2280698           529.8 ns/op
BenchmarkDatabases/BadgerDB
BenchmarkDatabases/BadgerDB/BadgerDB-Set
BenchmarkDatabases/BadgerDB/BadgerDB-Set-10           136003          8731 ns/op
BenchmarkDatabases/BadgerDB/BadgerDB-Get
BenchmarkDatabases/BadgerDB/BadgerDB-Get-10           575146          1889 ns/op
BenchmarkDatabases/BBolt
BenchmarkDatabases/BBolt/BBolt-Set
BenchmarkDatabases/BBolt/BBolt-Set-10                    127      10030328 ns/op
BenchmarkDatabases/BBolt/BBolt-Get
BenchmarkDatabases/BBolt/BBolt-Get-10                2023285           581.7 ns/op
PASS
```
