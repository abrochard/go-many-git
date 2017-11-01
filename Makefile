BINARY=gmg

build:
	if [ ! -d ${GOPATH}/src/github.com/fatih/color ] ; then go get -x github.com/fatih/color ; fi
	go build -o ${BINARY}

install:
	cp gmg /usr/bin/gmg

clean:
	if [ -f *~ ] ; then rm *~ ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${GOPATH}/bin/${BINARY} ] ; then rm ${GOPATH}/bin/${BINARY} ; fi
	if [ -d ${GOPATH}/src/github.com/fatih/color ] ; then rm -rf ${GOPATH}/src/github.com/fatih ; fi

.PHONY: build clean install
