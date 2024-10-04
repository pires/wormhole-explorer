#
# ATTENTION: this image must be built from the root of the repository.
# e.g. docker build -f spy/Dockerfile -t spy:latest .
#
ARG BUILDPLATFORM="linux/amd64"
ARG DEBIAN_FRONTEND="noninteractive"

#
# This stage builds the spy binary.
#
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.21-bookworm AS build
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install build dependencies.
RUN apt update; apt install -y build-essential

# Copy the source code to the build container.
ADD ./spy /source/spy
ADD ./common /source/common 

# Build the spy binary.
WORKDIR /source/spy
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

# Add the spy binary.
COPY --from=build "/source/spy/spy" "/spy"

# Run the binary.
ENTRYPOINT ["/spy"]
