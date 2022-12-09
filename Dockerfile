FROM gcr.io/distroless/base:nonroot
WORKDIR /
COPY petra /
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT [ "/petra" ]