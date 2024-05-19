FROM --platform=${BUILDPLATFORM} golang:1.22.3 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /build/serpentinised .

FROM --platform=${TARGETPLATFORM} scratch
COPY --from=builder /build/serpentinised /usr/bin/serpentinised

ENV SERPENTINISED_BIND="0.0.0.0:6379"
ENV PATH="/usr/bin"
EXPOSE 6379

CMD ["serpentinised"]
