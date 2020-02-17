# Copyright 2018 Adrian Todorov All Rights Reserved.
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

# Pull base golang image to compile the drone-firebase plugin
FROM golang:1.13-alpine

# add project files
RUN mkdir /drone
ADD . /drone/
WORKDIR /drone

# Install dependencies for go build
RUN apk add --no-cache git build-base

# Build drone-firebase
RUN go build -o /bin/drone-firebase . 

# Pull base alpine image for rest
FROM alpine:3.11

# Upgrade and install dependencies
RUN apk add --no-cache nodejs npm 

# Install required firebase-tools NPM package.
RUN npm install -g firebase-tools

# Copy the drone-firebase binary compiled in the first stage
COPY --from=0 /bin/drone-firebase /bin/drone-firebase

ENTRYPOINT ["/bin/drone-firebase"]
