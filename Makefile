BINARY=gmg

build:
	go build -o ${BINARY}

install:
	go get github.com/fatih/color
	go build -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
