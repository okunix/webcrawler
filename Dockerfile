FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app
RUN apk add --no-cache make
COPY . .
ARG TARGETARCH
RUN make TARGETARCH=${TARGETARCH}

FROM alpine:3.22

COPY --from=builder /app/bin/webcrawler /bin/webcrawler
ENTRYPOINT [ "/bin/webcrawler" ]
