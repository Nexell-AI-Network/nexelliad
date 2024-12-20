ARG KASPAD_IMAGE
ARG KASPAMINER_IMAGE

FROM ${KASPAD_IMAGE} as nexelliad
FROM ${KASPAMINER_IMAGE} as nexelliaminer

FROM golang:1.19-alpine

RUN mkdir -p /go/src/github.com/Nexell-AI-Network/nexelliad/v2/

WORKDIR /go/src/github.com/Nexell-AI-Network/nexelliad/v2/

RUN apk add bash build-base git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

COPY --from=nexelliad /app/ /app/
COPY --from=nexelliaminer /app/ /app/
ENV PATH="/app:${PATH}"

WORKDIR /go/src/github.com/Nexell-AI-Network/nexelliad/v2/stability-tests

RUN go install ./...

ENTRYPOINT ["./run/run.sh"]