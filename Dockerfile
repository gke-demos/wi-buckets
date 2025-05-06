FROM golang AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o demo

FROM gcr.io/distroless/static
COPY --from=builder /app/demo /demo


