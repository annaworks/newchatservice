FROM golang:1.15.0-alpine3.12 as BUILDER

ARG PORT=3000

WORKDIR /surubot

ADD . .

RUN go build -o surubot ./cmd/surubot.go

FROM alpine:3.12

COPY --from=BUILDER /surubot/surubot /bin/.

EXPOSE $PORT

CMD ["surubot"]
