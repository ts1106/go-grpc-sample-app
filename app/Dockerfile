ARG VERSION=1.18
ARG USERNAME=golang
ARG USER_UID=1000
ARG USER_GID=1000

FROM golang:${VERSION}-bullseye as base
ARG USERNAME
ARG USER_UID
ARG USER_GID
WORKDIR /workspace
RUN groupadd --gid ${USER_GID} ${USERNAME} \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} \
    && apt-get update \
    && apt-get install -y sudo \
    && echo ${USERNAME} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}
RUN mkdir -p /go/pkg
RUN chown -R ${USERNAME} /go/pkg

FROM base as development
ARG USERNAME
WORKDIR /workspace/app
RUN apt-get update &&\
    apt-get install -y git
USER ${USERNAME}
COPY --chown=${USERNAME}:${USERNAME} . .
RUN go install golang.org/x/tools/gopls@latest
RUN go mod tidy
