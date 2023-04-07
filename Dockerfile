FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /
COPY petra /
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT [ "/petra" ]