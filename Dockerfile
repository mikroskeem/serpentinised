FROM --platform=${BUILDPLATFORM} golang:1.22.3 AS builder

WORKDIR $GOPATH/src/github.com/astei/serpentinised

COPY go.mod go.sum .
RUN go mod download

COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /serpentinised .

FROM --platform=${TARGETPLATFORM} scratch
COPY --from=builder serpentinised /usr/bin/serpentinised

ENV SERPENTINISED_BIND="0.0.0.0:6379"
ENV PATH="/usr/bin"
EXPOSE 6379

CMD ["serpentinised"]
