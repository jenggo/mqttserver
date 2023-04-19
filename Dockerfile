ARG BASE_IMAGE=golang:1.20-alpine
ARG CONTAINER=scratch

FROM ${BASE_IMAGE}

ARG APP_NAME=mqttserver

RUN apk add --update upx gcc musl-dev git tzdata && wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s

WORKDIR /src

COPY go.* ./
RUN go mod download -x

COPY . .
RUN go vet . && golangci-lint run

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ${APP_NAME} && upx -q --best --lzma ${APP_NAME}


FROM ${CONTAINER}

WORKDIR /app

COPY --from=0 /src/mqttserver .
COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo
COPY certs/ certs/
COPY config.yaml .

EXPOSE 8883

ENTRYPOINT ["./mqttserver"]

