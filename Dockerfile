ARG GO_VERSION=1.26.2

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/tcp-echo .

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

USER app

EXPOSE 8090

ENTRYPOINT ["/bin/tcp-echo"]

COPY --from=build /bin/tcp-echo /bin/tcp-echo
