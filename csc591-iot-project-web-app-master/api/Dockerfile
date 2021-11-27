FROM golang:1.10-stretch
ARG src_dir

RUN go get -u github.com/alecthomas/gometalinter
RUN go get -u github.com/golang/dep/cmd/dep
RUN gometalinter --install

COPY Gopkg.toml Gopkg.lock $src_dir
WORKDIR $src_dir
RUN dep ensure -v -vendor-only

COPY . .

EXPOSE 8000

CMD ["go", "run", "main.go"]
