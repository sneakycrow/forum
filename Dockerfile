FROM golang

WORKDIR /app

COPY . .

RUN make build

ENTRYPOINT ["/app/bin/forum"]