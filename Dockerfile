FROM golang:1.17-alpine AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /petra

FROM gcr.io/distroless/base:nonroot
WORKDIR /
COPY --from=build /petra /petra
USER nonroot:nonroot
EXPOSE 3000
ENTRYPOINT [ "/petra" ]