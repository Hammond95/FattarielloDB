# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /

# add safe user
ARG UNAME=myuser
ARG UID=1000
RUN useradd -o -r -u $UID $UNAME

# compilator parameters
ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .

RUN go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o /myapp

RUN mkdir /temp && \
    chown $UID /temp && \
    chmod 555 /temp


###################
# The final stage #
###################
FROM scratch

ENV PATH="/"

COPY --from=builder /myapp /FattarielloDB
COPY --from=builder /temp /tmp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER $UID
CMD ["main"]
