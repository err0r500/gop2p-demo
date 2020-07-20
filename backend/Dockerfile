FROM golang:1.13-alpine3.11 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /build/app .
ENTRYPOINT [ "./app" ]