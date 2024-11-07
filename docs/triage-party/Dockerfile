# syntax=docker/dockerfile:1.1-experimental

# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM golang:1.22.6  as builder

RUN go get github.com/google/triage-party/cmd/server
RUN go install github.com/google/triage-party/cmd/server@latest
RUN mkdir -p /.cache
WORKDIR /git
RUN git clone https://github.com/google/triage-party.git
ENTRYPOINT ["/go/bin/server"]

FROM gcr.io/distroless/base:latest
WORKDIR /app
COPY --chown=nobody:nobody --from=builder /go/bin/server /app/triage-party
COPY --chown=nobody:nobody --from=builder /git/triage-party/site /app/site
COPY --chown=nobody:nobody --from=builder /git/triage-party/third_party /app/third_party
COPY --chown=nobody:nobody config/config.yaml /app/config/config.yaml
COPY --chown=nobody:nobody --from=builder /.cache /app/.cache
ENV HOME /app
USER nobody
ENTRYPOINT ["/app/triage-party"]
