PROJECT_NAME=rms-mirror
BINARY_NAME=${PROJECT_NAME}.out
SOURCE_MAIN=${PROJECT_NAME}.go
LDFLAGS="-X main.Version=`git tag --sort=-version:refname | head -n 1`"

all: build test

build:
	go build -ldflags ${LDFLAGS} -o ${BINARY_NAME} ${SOURCE_MAIN}

test:
	go test -v ${SOURCE_MAIN}

run:
	go build -ldflags ${LDFLAGS} -o ${BINARY_NAME} ${SOURCE_MAIN}
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}