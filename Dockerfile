FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

EXPOSE 8080

CMD ["./main"]
```

**Dan ganti baris pertama `go.mod` menjadi:**
```
go 1.23
```

---

## Langkah-langkahnya di VSCode:

**1. Edit `go.mod`** — ubah baris:
```
go 1.25.5  ← hapus ini
```
jadi:
```
go 1.23
```

**2. Edit `Dockerfile`** — ubah:
```
FROM golang:1.22-alpine  ← hapus ini
```
jadi:
```
FROM golang:1.23-alpine