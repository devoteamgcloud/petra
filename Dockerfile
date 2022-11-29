FROM golang:1.19 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/petra

FROM gcr.io/distroless/base:debug
COPY --from=build /go/bin/petra /
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT [ "/petra" ]
