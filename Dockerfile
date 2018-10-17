FROM golang:1.10.3-alpine as builder
WORKDIR /go/src/github.com/maorfr/orca/
COPY . .
RUN apk --no-cache add git make \
    && make

FROM alpine:3.8
ARG HELM_VERSION=v2.11.0
ARG HELM_OS_ARCH=linux-amd64
RUN apk --no-cache add ca-certificates git bash curl \
    && wget -q https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-${HELM_OS_ARCH}.tar.gz \
    && tar -zxvf helm-${HELM_VERSION}-${HELM_OS_ARCH}.tar.gz ${HELM_OS_ARCH}/helm \
    && mv ${HELM_OS_ARCH}/helm /usr/local/bin/helm \
    && rm -rf ${HELM_OS_ARCH} helm-${HELM_VERSION}-${HELM_OS_ARCH}.tar.gz
ARG LINKERD_VERSION=stable-2.0.0
ARG LINKERD_OS=linux
RUN wget -q https://github.com/linkerd/linkerd2/releases/download/${LINKERD_VERSION}/linkerd2-cli-${LINKERD_VERSION}-${LINKERD_OS} -O linkerd \
    && chmod +x linkerd \
    && mv linkerd /usr/local/bin/linkerd
COPY --from=builder /go/src/github.com/maorfr/orca/bin/orca /usr/local/bin/orca
RUN addgroup -g 1001 -S orca \
    && adduser -u 1001 -D -S -G orca orca
USER orca
WORKDIR /home/orca
RUN helm init -c \
    && helm plugin install https://github.com/chartmuseum/helm-push \
    && helm plugin install https://github.com/maorfr/helm-inject
CMD ["orca"]
