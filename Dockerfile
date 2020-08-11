FROM golang:1 as build

WORKDIR /src
ADD . /src

RUN go build -o /bouncycastle

FROM gcr.io/distroless/base
COPY --from=build /bouncycastle /bouncycastle
USER 65534
ENTRYPOINT ["/bouncycastle"]
