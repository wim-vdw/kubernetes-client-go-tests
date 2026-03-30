FROM golang:1.26-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /k8s-cluster-info ./...

FROM gcr.io/distroless/static-debian13

WORKDIR /

COPY --from=build-stage /k8s-cluster-info /k8s-cluster-info

ENTRYPOINT ["/k8s-cluster-info"]
