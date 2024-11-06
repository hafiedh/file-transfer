FROM golang:1.23 as build

WORKDIR /app
COPY . .

RUN make compile

FROM alpine:3.13

COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Asia/Jakarta
ENV ZONEINFO=/zoneinfo.zip

COPY --from=build /app/bin/app /app

EXPOSE 3000
ENTRYPOINT ["/app"]
