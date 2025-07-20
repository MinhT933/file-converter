FROM golang:1.24-bullseye

RUN apt-get update && \
    apt-get install -y --fix-missing wget fontconfig fonts-dejavu xvfb libxrender1 libxext6 libjpeg62-turbo && \
    wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb && \
    apt-get install -y ./wkhtmltox_0.12.6-1.buster_amd64.deb && \
    rm -rf /var/lib/apt/lists/* wkhtmltox_0.12.6-1.buster_amd64.deb

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o myapp ./cmd/server

EXPOSE 8080

CMD ["./myapp"]