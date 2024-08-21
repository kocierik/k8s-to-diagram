FROM golang AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o k8s-to-diagram .

RUN chmod +x k8s-to-diagram

FROM busybox

COPY --from=builder /app/k8s-to-diagram .

RUN ls -la 

CMD ["./k8s-to-diagram"]

