############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk add git
RUN apk add openldap-clients
WORKDIR $GOPATH/src/yarn
COPY /yarn $GOPATH/src/yarn
ENV GOBIN=$GOPATH/bin

RUN go get -d -v
#RUN go build -o /go/bin/yarn-prometheus-exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/yarn-prometheus-exporter

############################
# STEP 2 build a small image
############################
#FROM scratch
## Copy the  executable.
#COPY --from=builder /go/bin/yarn-prometheus-exporter /go/bin/yarn-prometheus-exporter
# Run the binary.
ENTRYPOINT ["/go/bin/yarn-prometheus-exporter"]
