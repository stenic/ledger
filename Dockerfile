FROM golang:1.19 AS build-server

WORKDIR /workspace
COPY ./go.* .
RUN go mod download
COPY ./cmd ./cmd
COPY ./graph ./graph
COPY ./internal ./internal
COPY ./migrations ./migrations
RUN CGO_ENABLED=1 GOOS=linux go build -a -o ledger ./cmd/ledger


FROM node:alpine AS build-ui

WORKDIR /workspace
COPY ./ui/package*.json .
RUN npm ci
ADD ./ui/ .
RUN npm run build


FROM alpine AS downloader

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64
RUN chmod +x /usr/local/bin/dumb-init

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base:debug-nonroot
WORKDIR /app

COPY --from=downloader /usr/local/bin/dumb-init /app/dumb-init
COPY --from=build-server /workspace/ledger /app/ledger
COPY --from=build-server /workspace/migrations /app/migrations
COPY --from=build-ui /workspace/build /app/static
USER 65532:65532

ENV STATIC_ASSET_PATH=/app/static
ENV LOG_FORMAT=json

ENTRYPOINT ["/app/dumb-init", "--", "/app/ledger", "server"]
