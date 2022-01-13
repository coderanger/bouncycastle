FROM golang:1 as build

WORKDIR /src
ADD . /src

RUN go build -o /bouncycastle

FROM gcr.io/distroless/base
LABEL org.opencontainers.image.source https://github.com/coderanger/bouncycastle
COPY --from=build /bouncycastle /bouncycastle
USER 65534
ENTRYPOINT ["/bouncycastle"]
