FROM golang:latest as builder

ENV CGO_ENABLED=0

COPY . /go/src/github.com/akhmettolegen/proxy
WORKDIR /go/src/github.com/akhmettolegen/proxy
RUN \
    version=git describe --abbrev=6 --always --tag; \
    echo "version=$version" && \
    cd cmd/app && \
    go build -a -tags proxy -installsuffix proxy -ldflags "-X main.version=${version} -s -w" -o /go/bin/proxy -mod vendor

#
# Контейнер рантайма
#
FROM scratch
COPY --from=builder /go/bin/proxy /bin/proxy

ENTRYPOINT ["/bin/proxy"]