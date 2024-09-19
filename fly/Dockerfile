#
# ATTENTION: this image must be built from the root of the repository.
# e.g. docker build -f fly/Dockerfile -t fly:latest .
#
ARG BUILDPLATFORM="linux/amd64"

#
# This stage builds the fly binary. We use an Alpine Linux image in order
# to respect the fact this distro is musl-based, so we shouldn't build 
# binaries in a glibc-based image and expect them to run in Alpine Linux.
#
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.21-alpine AS build-fly

# Install build dependencies.
RUN apk update; apk add build-base

# Copy the source code to the build container.
ADD ./fly /source/fly
ADD ./common /source/common 

# Build the fly binary.
WORKDIR /source/fly
RUN CGO_ENABLED=1 make build

#
# This stage builds the libwasvm.x86_64.so library. We use an Alpine Linux
# image in order to respect the fact this distro is musl-based, so we shouldn't
# build libraries in a glibc-based image and expect them to run in Alpine Linux.
#
FROM --platform=${BUILDPLATFORM} cosmwasm/libwasmvm-builder:0100-alpine AS build-libwasvm

# Install build dependencies.
RUN apk update; apk add git

# Download libwasvm code.
RUN git clone https://github.com/CosmWasm/wasmvm.git --depth 1 --branch v1.5.4 /wasmvm

# Build libwasvm shared-library.
WORKDIR /wasmvm/libwasmvm
RUN cargo build --release

#
# This stage updates the CA certificates present in Alpine Linux, which
# will be used by the final image.
#
FROM --platform=${BUILDPLATFORM} alpine AS ca-certs

# Install ca-certificates.
RUN apk update; apk add ca-certificates
#; update ca-certificates

# Build the final (minimalistic) Alpine Linux image.
FROM --platform=${BUILDPLATFORM} alpine

# Copy the CA certificates from the ca-certs image.
COPY --from=ca-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Add the libwasvm.x86_64.so library.
COPY --from=build-libwasvm /wasmvm/libwasmvm/target/release/libwasmvm..so /usr/lib/libwasmvm.x86_64.so

# Add the fly binary.
COPY --from=build-fly "/source/fly/fly" "/fly"

# Run the binary.
ENTRYPOINT ["/fly"]
