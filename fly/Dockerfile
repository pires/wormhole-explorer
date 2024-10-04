#
# ATTENTION: this image must be built from the root of the repository.
# e.g. docker build -f fly/Dockerfile -t fly:latest .
#
ARG BUILDPLATFORM="linux/amd64"
ARG DEBIAN_FRONTEND="noninteractive"

#
# This stage builds the fly binary.
#
FROM --platform=${BUILDPLATFORM} docker.io/golang:1.21-bookworm AS build-fly
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install build dependencies.
RUN apt update; apt install -y build-essential

# Copy the source code to the build container.
ADD ./fly /source/fly
ADD ./common /source/common 

# Build the fly binary.
WORKDIR /source/fly
RUN CGO_ENABLED=1 make build

#
# This stage builds the libwasvm.x86_64.so shared-library.
#
FROM --platform=${BUILDPLATFORM} cosmwasm/libwasmvm-builder:0100-debian AS build-libwasmvm
ENV DEBIAN_FRONTEND=${DEBIAN_FRONTEND}

# Install build dependencies.
RUN apt update; apt install -y git

# Download libwasvm code.
RUN git clone https://github.com/CosmWasm/wasmvm.git --depth 1 --branch v1.5.4 /wasmvm

# Build libwasvm shared-library.
WORKDIR /wasmvm/libwasmvm
RUN cargo build --release

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

# Add the libwasvm.x86_64.so library.
COPY --from=build-libwasmvm /wasmvm/libwasmvm/target/release/libwasmvm.so /usr/lib/libwasmvm.x86_64.so

# Add the fly binary.
COPY --from=build-fly "/source/fly/fly" "/fly"

# Run the binary.
ENTRYPOINT ["/fly"]
