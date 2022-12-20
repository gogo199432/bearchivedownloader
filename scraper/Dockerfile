FROM golang:1.19-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY / ./

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /scraper

FROM scratch
COPY --from=build /scraper /scraper
# Needed for colly
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./config.yaml ./
CMD [ "/scraper" ]
