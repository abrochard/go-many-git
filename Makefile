BINARY=gmg

build:
	go get github.com/fatih/color
	go build -o ${BINARY}

install: build
	go install -x ${BINARY}.go

clean:
	if [ -f *~ ] ; then rm *~ ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${GOPATH}/bin/${BINARY} ] ; then rm ${GOPATH}/bin/${BINARY} ; fi

.PHONY: clean install
