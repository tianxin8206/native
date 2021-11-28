FROM golang:1.17-alpine AS base
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /src
COPY . .
RUN go build ./src/main.go

FROM scratch
EXPOSE 8080
WORKDIR /app
COPY --from=base /src/main .
COPY --from=base /src/src/config.json .
CMD [ "./main" ]