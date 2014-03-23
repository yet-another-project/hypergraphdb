go test -c
./hypergraphdb.test -test.bench=. -test.cpuprofile=cpu.out
go tool pprof hypergraphdb.test cpu.out
