FROM golang:1.14.4-alpine3.12 as build

RUN apk update && \
	apk add --no-cache \
	git \
	ca-certificates

ENV CGO_ENABLED 0

WORKDIR /opt
COPY ./ ./
RUN go mod download
RUN go build -o /main /opt/cmd/aws-budgets-exporter/main.go

# ---

FROM scratch

COPY --from=build /main /main
COPY --from=build /etc/ssl/certs /etc/ssl/certs

CMD ["/main"]
