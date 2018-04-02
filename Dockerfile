FROM golang:1.10

WORKDIR /go/src/app

RUN curl -fsSL -o \
	/usr/local/bin/dep \
	https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
	chmod +x /usr/local/bin/dep

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only

COPY . .

