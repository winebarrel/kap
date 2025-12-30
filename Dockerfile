FROM golang:1.25 AS build

WORKDIR /src
COPY go.* /src/
RUN go mod download

COPY *.go /src/
COPY cmd/kap/*.go /src/cmd/kap/
ARG KAP_VERSION
RUN CGO_ENABLED=0 go build -o kap -ldflags "-X main.version=${KAP_VERSION#v}" ./cmd/kap

FROM gcr.io/distroless/static

COPY --from=build /src/kap /

ENTRYPOINT ["/kap"]
