FROM golang:1.10.3-alpine as builder
WORKDIR /go/src/github.com/maorfr/orca/
COPY . .
RUN apk --no-cache add git glide \
    && glide up \
    && for f in $(find test -type f -name "*.go"); do go test -v $f; done \
    && CGO_ENABLED=0 GOOS=linux go build -o orca cmd/orca.go

FROM alpine:3.8
ARG HELM_VERSION=v2.11.0
RUN apk --no-cache add ca-certificates curl git bash \
    && curl -O https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz \
    && mv linux-amd64/helm /usr/local/bin/helm \
    && rm -f /helm-${HELM_VERSION}-linux-amd64.tar.gz
COPY --from=builder /go/src/github.com/maorfr/orca/orca /usr/local/bin/orca
RUN addgroup -g 1001 -S orca \
    && adduser -u 1001 -D -S -G orca orca
USER orca
WORKDIR /home/orca
RUN helm init -c \
    && helm plugin install https://github.com/chartmuseum/helm-push \
    && helm plugin install https://github.com/maorfr/helm-inject
CMD ["orca"]
