# Build the process-rest binary
ARG GOVERSION=1.19
FROM golang:$GOVERSION as builder
ARG GOVERSION=1.19
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ENV GO111MODULE="on" \
    GOOS=linux       \
    GOARCH=amd64

WORKDIR /github.com/w6d-io/process-rest
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Copy the go source
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/

# Build
RUN  go build    \
    -ldflags="-X 'github.com/w6d-io/process-rest/internal/config.Version=${VERSION}' -X 'github.com/w6d-io/process-rest/internal/config.Revision=${VCS_REF}' -X 'github.com/w6d-io/process-rest/internal/config.GoVersion=go${GOVERSION}' -X 'github.com/w6d-io/process-rest/internal/config.Built=${BUILD_DATE}' -X 'github.com/w6d-io/process-rest/internal/config.OsArch=${GOOS}/${GOARCH}'" \
    -a -o process-rest cmd/process-rest/main.go


FROM w6dio/kubectl:v1.4.0
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ARG PROJECT_URL
ARG USER_EMAIL="david.alexandre@w6d.io"
ARG USER_NAME="David ALEXANDRE"
LABEL maintainer="${USER_NAME} <${USER_EMAIL}>" \
    io.w6d.ci.vcs-ref=$VCS_REF       \
    io.w6d.ci.vcs-url=$PROJECT_URL   \
    io.w6d.ci.build-date=$BUILD_DATE \
    io.w6d.ci.version=$VERSION
RUN curl -sSL https://git.io/get-mo -o mo && \
    chmod +x mo && \
    mv mo /usr/local/bin/ && \
    apt update && apt install -y postgresql-client && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=builder /github.com/w6d-io/process-rest/process-rest /usr/local/bin/process-rest

ENTRYPOINT ["/usr/local/bin/process-rest"]
CMD ["serve"]

