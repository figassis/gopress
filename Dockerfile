FROM golang:alpine as builder
ENV GO111MODULE on
RUN mkdir -p /build

WORKDIR /build
COPY go.mod go.sum ./
RUN apk add git && go mod download
ADD . .
RUN mkdir -p /goinagbe/uploads && mv assets /goinagbe/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /goinagbe/app .

# FROM alpine
FROM gcr.io/distroless/base
COPY --from=builder /goinagbe /goinagbe

WORKDIR /goinagbe

ENV DATADIR /goinagbe

VOLUME [ "/goinagbe/uploads" ]

EXPOSE 80

CMD ["./app"]