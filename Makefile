
test:
	go test ./*.go

bench:
	go test -bench=. ./*.go
