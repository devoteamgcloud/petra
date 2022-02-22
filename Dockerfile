FROM golang:1.17-alpine AS build

ENV BASEDIR /go/src/github.com/arthur-laurentdka/petra

WORKDIR ${BASEDIR}

ADD . ${BASEDIR}

RUN go install github.com/arthur-laurentdka/petra

FROM gcr.io/distroless/base:nonroot

COPY --from=build /go/bin/petra /

ENTRYPOINT [ "/petra" ]