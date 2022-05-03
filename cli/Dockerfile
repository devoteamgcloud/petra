FROM gcr.io/distroless/base:nonroot
WORKDIR /
COPY petra-cli /
USER nonroot:nonroot
EXPOSE 3001
ENTRYPOINT [ "/petra-cli" ]