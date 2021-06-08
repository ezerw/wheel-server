# Build  stage
FROM golang:1.16-alpine as build
WORKDIR /go/bin/wheel
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o api ./cmd/api/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o migration ./cmd/migration/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o seeder ./cmd/seeder/

# Run stage
FROM alpine
WORKDIR /bin/wheel
COPY --from=build /go/bin/wheel/api ./api
COPY --from=build /go/bin/wheel/migration ./migration
COPY --from=build /go/bin/wheel/seeder ./seeder
COPY --from=build /go/bin/wheel/app.env ./app.env
RUN chmod +x api migration seeder

ENTRYPOINT ["./api"]
