FROM gcr.io/distroless/base:nonroot
WORKDIR /
COPY petra /
USER nonroot:nonroot
EXPOSE 3000
ENTRYPOINT [ "/petra" ]