BINARY=gmg

build:
	go get github.com/fatih/color
	go build -o ${BINARY}

install:
	cp gmg /usr/bin/gmg

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: build clean install
