FROM golang:1.23 as go_app_dev

WORKDIR /app/go-clean-api-scaffold
COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
RUN go install github.com/vektra/mockery/v2@v2.45.0

COPY . .
RUN go mod tidy

ENTRYPOINT ["./entrypoint.sh"]
CMD air


FROM golang:1.23 as go_app_prod

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/bin/go_app cmd/main.go
ENTRYPOINT ["./entrypoint.sh"]
CMD /app/bin/go_app
