FROM 175781750719.dkr.ecr.us-east-1.amazonaws.com/base/image/golang:1.24-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3 && \
	swag init -g ./cmd/main.go -o ./api

RUN go build -v -tags appsec -o ./cmd/main ./cmd/main.go

FROM 175781750719.dkr.ecr.us-east-1.amazonaws.com/base/image/golang:1.24-alpine AS final

ENV CGO_ENABLED=1

RUN addgroup -S service && \
	adduser app -S -G service

COPY --from=build /app/cmd/ /home/app

RUN chmod -R 775 /home/app
RUN chown -R app:service /home/app

USER app
WORKDIR /home/app

EXPOSE 8080

CMD ./main