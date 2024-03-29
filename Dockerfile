FROM golang:1.16.1-alpine AS build
ARG DEST_DIR="/go/src/bitbucket_pipelines_runner"
COPY . $DEST_DIR
WORKDIR $DEST_DIR
RUN apk update && apk add git openssh
RUN go build -o /usr/bin/bpr ./cmd/bpr

FROM alpine:3
COPY --from=build /usr/bin/bpr /bin/bpr
ENTRYPOINT [ "/bin/bpr" ]