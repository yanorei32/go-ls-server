FROM golang

MAINTAINER yanorei32

RUN go get -u github.com/labstack/echo/ && \
	go get -u github.com/labstack/echo/middleware && \
	go get -u gopkg.in/go-playground/validator.v9 && \
	go get -u gopkg.in/yaml.v3

