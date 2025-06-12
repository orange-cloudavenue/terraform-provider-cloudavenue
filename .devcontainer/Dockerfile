ARG VARIANT="1-bullseye"

FROM --platform=linux/amd64 mcr.microsoft.com/vscode/devcontainers/go:${VARIANT}

ARG USERNAME=cav
ARG USER_UID=1001
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME

RUN /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

COPY .terraformrc /home/cav/.terraformrc
ENV PATH="/home/linuxbrew/.linuxbrew/bin:${PATH}"

RUN brew install "go-task/tap/go-task"