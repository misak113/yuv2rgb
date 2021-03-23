
CGO_CFLAGS := -I/opt/intel/oneapi/ipp/latest/include
CGO_LDFLAGS := -L/opt/intel/oneapi/ipp/latest/lib/intel64
# Tests has to be staticly linked because otherwise they cannot find libs
CGO_LDFLAGS_TEST := -static ${CGO_LDFLAGS}

test:
	export CGO_CFLAGS="$(CGO_CFLAGS)" ; \
	export CGO_LDFLAGS="$(CGO_LDFLAGS_TEST)" ; \
	go test ./*.go

bench:
	export CGO_CFLAGS="$(CGO_CFLAGS)" ; \
	export CGO_LDFLAGS="$(CGO_LDFLAGS_TEST)" ; \
	go test -bench=. ./*.go

clean:
	go clean
