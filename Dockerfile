# Build the process-rest binary
ARG GOVERSION=1.21
FROM golang:$GOVERSION as builder
ARG GOVERSION=1.20
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
COPY . .

# Build
RUN  go build    \
    -ldflags="-X 'github.com/w6d-io/process-rest/internal/config.Version=${VERSION}' -X 'github.com/w6d-io/process-rest/internal/config.Revision=${VCS_REF}' -X 'github.com/w6d-io/process-rest/internal/config.Built=${BUILD_DATE}'" \
    -a -o process-rest main.go


FROM w6dio/kubectl:latest
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ARG PROJECT_URL
ARG USER_EMAIL="david.alexandre@w6d.io"
ARG USER_NAME="David ALEXANDRE"
LABEL maintainer="${USER_NAME} <${USER_EMAIL}>" \
    io.w6d.vcs-ref=$VCS_REF       \
    io.w6d.vcs-url=$PROJECT_URL   \
    io.w6d.build-date=$BUILD_DATE \
    io.w6d.version=$VERSION
RUN curl -sSL https://git.io/get-mo -o mo && \
    chmod +x mo && \
    mv mo /usr/local/bin/ && \
    apt update && apt install -y postgresql-client && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /
COPY --from=builder /github.com/w6d-io/process-rest/process-rest /usr/local/bin/process-rest

ENTRYPOINT ["/usr/local/bin/process-rest"]
CMD ["serve"]

