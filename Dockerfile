FROM node:15 AS ui-builder

WORKDIR /gopds
COPY . .
RUN cd ui && npm install
RUN cd ui && npm run build -- --prod --output-path=../public

FROM golang:1.15 AS app-builder

WORKDIR /gopds
COPY . .
RUN go mod download
RUN go get -v github.com/golang/mock/mockgen@v1.4.4
RUN go generate ./...
RUN CGO_ENABLED=1 go build -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension,osusergo,netgo -o gopds -v

FROM scratch

COPY --from=app-builder /gopds/gopds /gopds
COPY --from=ui-builder /gopds/public /public
COPY --from=app-builder /gopds/migrations /migrations
ENTRYPOINT ["/gopds"]
