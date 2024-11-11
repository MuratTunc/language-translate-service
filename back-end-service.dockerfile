# base go image
FROM golang:1.23-alpine AS builder

RUN mkdir /app

COPY backendServiceApp /app

CMD [ "/app/backendServiceApp" ]