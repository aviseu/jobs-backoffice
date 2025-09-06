FROM golang:1.25-alpine AS builder

ARG APP

WORKDIR /app

COPY ./ ./

RUN apk add --update --no-cache git make
RUN make build-${APP}

# ---

FROM gcr.io/distroless/static:nonroot AS api

ARG APP

WORKDIR /

COPY --from=builder /app/dist/app /app
COPY --from=builder /app/config/migrations/ ./config/migrations/

EXPOSE 8080 9090 4000

ENTRYPOINT ["/app"]


FROM gcr.io/distroless/static:nonroot AS job

ARG APP

WORKDIR /

COPY --from=builder /app/dist/app /app
COPY --from=builder /app/config/migrations/ ./config/migrations/

ENTRYPOINT ["/app"]
