FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

COPY *.go /src/
COPY cmd /src/cmd
RUN CGO_ENABLED=0 go build -o kap ./cmd/kap

FROM gcr.io/distroless/static

COPY --from=build /src/kap /

ENTRYPOINT ["/kap"]
