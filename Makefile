deps: 
	glide install

build: deps
	go build

install: deps
	go install