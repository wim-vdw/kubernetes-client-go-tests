FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /k8s-cluster-info ./...

FROM alpine:3.21

RUN apk add --no-cache tzdata

WORKDIR /

COPY --from=build-stage /k8s-cluster-info /k8s-cluster-info

CMD ["/k8s-cluster-info"]
