FROM alpine:3.9

ARG HELM_VERSION=v2.14.0
ARG HELM_OS_ARCH=linux-amd64

RUN apk --no-cache add ca-certificates git bash curl jq \
  && wget -q https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-${HELM_OS_ARCH}.tar.gz \
  && tar -zxvf helm-${HELM_VERSION}-${HELM_OS_ARCH}.tar.gz ${HELM_OS_ARCH}/helm \
  && mv ${HELM_OS_ARCH}/helm /usr/local/bin/helm \
  && rm -rf ${HELM_OS_ARCH} helm-${HELM_VERSION}-${HELM_OS_ARCH}.tar.gz

COPY orca /usr/local/bin/orca

RUN addgroup -g 1001 -S orca \
  && adduser -u 1001 -D -S -G orca orca

USER orca

WORKDIR /home/orca

ENV HELM_HOME /home/orca/.helm

RUN helm init -c \
  && helm plugin install https://github.com/chartmuseum/helm-push

CMD ["orca"]
