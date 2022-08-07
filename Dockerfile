FROM  golang:alpine3.16 as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN  CGO_ENABLED=0 GOOS=linux go build -a -o ./gogoapps

FROM scratch
COPY --from=builder /app/gogoapps ./gogoapps
ENTRYPOINT [ "./gogoapps" ]