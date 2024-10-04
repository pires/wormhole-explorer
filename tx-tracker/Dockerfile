#
# ATTENTION: this image must be built from the root of the repository.
# e.g. docker build -f tx-tracker/Dockerfile -t tx-tracker:latest .
#
ARG BUILDPLATFORM="linux/amd64"
ARG DEBIAN_FRONTEND="noninteractive"

#
# This stage builds the tx-tracker binary.
#
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.21-bookworm AS build
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install build dependencies.
RUN apt update; apt install -y build-essential

# Copy the source code to the build container.
ADD ./tx-tracker /source/tx-tracker
ADD ./common /source/common 

# Build the tx-tracker binary.
WORKDIR /source/tx-tracker
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

# Add the tx-tracker binary.
COPY --from=build "/source/tx-tracker/tx-tracker" "/tx-tracker"

# Run the binary.
ENTRYPOINT ["/tx-tracker"]
