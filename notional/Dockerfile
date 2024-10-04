#
# ATTENTION: this image must be built from the root of the repository.
# e.g. docker build -f notional/Dockerfile -t notional:latest .
#
ARG BUILDPLATFORM="linux/amd64"
ARG DEBIAN_FRONTEND="noninteractive"

#
# This stage builds the notional binary.
#
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.21-bookworm AS build-notional
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install build dependencies.
RUN apt update; apt install -y build-essential

# Copy the source code to the build container.
ADD ./notional /source/notional
ADD ./common /source/common 

# Build the notional binary.
WORKDIR /source/notional
RUN make build

#
# This stage updates the CA certificates.
#
FROM --platform=${BUILDPLATFORM} debian:bookworm AS ca-certs
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install ca-certificates.
RUN apt update; apt upgrade -y ca-certificates

#
# Build the final (minimalistic) image.
#
FROM --platform=${BUILDPLATFORM} gcr.io/distroless/static-debian12:nonroot

# Copy the CA certificates from the ca-certs image.
COPY --from=ca-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Add the notional binary.
COPY --from=build-notional "/source/notional/notional" "/notional"

# Run the binary.
ENTRYPOINT ["/notional"]
