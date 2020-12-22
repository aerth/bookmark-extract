bin/bookmark-extract: *.go
	GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOBUILDFLAGS} -o ${OUTDIR}$@
clean:
	rm -rf bin
