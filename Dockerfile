# SPDX-FileCopyrightText: 2023 Luca Tagliavini <luca.tagliavini5@studio.unibo.it>
# SPDX-FileCopyrightText: 2025 Eyad Issa <eyadlorenzo@gmail.com>
#
# SPDX-License-Identifier: AGPL-3.0-only
# SPDX-License-Identifier: AGPL-3.0-or-later

ARG GO_VERSION=1.24
ARG ALPINE_VERSION=3.22

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download -x

COPY . /build
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags "-s -w" -o /build/filenameslinter ./cmd/

FROM alpine:${ALPINE_VERSION}
COPY --from=builder /build/filenameslinter /usr/local/bin/filenameslinter

ENTRYPOINT ["filenameslinter"]
