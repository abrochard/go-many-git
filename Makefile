BINARY=gmg

build:
	go build -o ${BINARY}

install:
	go install

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
