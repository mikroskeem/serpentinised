FROM golang:1.22.3 AS builder

WORKDIR $GOPATH/src/github.com/astei/serpentinised
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -ldflags="-s -w" -o /serpentinised .

FROM scratch
COPY --from=builder serpentinised /usr/bin/serpentinised

ENV PATH="/usr/bin"
EXPOSE 6379

CMD ["serpentinised", "-bind=0.0.0.0:6379"]
