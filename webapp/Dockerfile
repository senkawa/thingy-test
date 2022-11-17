FROM golang:1.19 AS builder
RUN apt update && apt install -y git

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /server

FROM gcr.io/distroless/static AS final
COPY --from=builder --chown=nonroot:nonroot /server /server

EXPOSE 80
CMD ["/server"]
