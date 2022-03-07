# Build the manager binary
FROM golang:1.15 as builder

# add safe user
ARG UNAME=myuser
ARG UID=1000
RUN useradd -o -r -u $UID $UNAME

# compilator parameters
ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOARCH=amd64 \
	GOOS=linux

RUN mkdir /app
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY proto/ proto/
COPY server/ server/


RUN ls -la -R /app
RUN go build -o /fattarielloServer -a -installsuffix cgo -ldflags '-w -extldflags "-static"' /app/server

RUN mkdir /temp && \
    chown $UID /temp && \
    chmod 555 /temp


###################
# The final stage #
###################
FROM scratch

ENV PATH="/"

COPY --from=builder /fattarielloServer /fattarielloServer
COPY --from=builder /temp /tmp
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER $UID
CMD ["/fattarielloServer"]
