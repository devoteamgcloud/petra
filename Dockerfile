FROM golang:1.17-alpine AS build

ENV BASEDIR /go/src/github.com/arthur-laurentdka/petra

WORKDIR ${BASEDIR}

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/petra

FROM gcr.io/distroless/base:nonroot

COPY --from=build /go/bin/petra /

EXPOSE 3000
ENTRYPOINT [ "/petra" ]