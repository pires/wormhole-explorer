#
# ATTENTION: this image must be built from the root of the repository.
# e.g. docker build -f jobs/Dockerfile -t jobs:latest .
#
ARG BUILDPLATFORM="linux/amd64"
ARG DEBIAN_FRONTEND="noninteractive"

#
# This stage builds the jobs binary.
#
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.21-bookworm AS build-jobs
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install build dependencies.
RUN apt update; apt install -y build-essential

# Copy the source code to the build container.
ADD ./jobs /source/jobs
ADD ./common /source/common 

# Build the jobs-app binary.
WORKDIR /source/jobs
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

# Add the jobs binary.
COPY --from=build-jobs "/source/jobs/jobs-app" "/jobs-app"

# Run the binary.
ENTRYPOINT ["/jobs-app"]
