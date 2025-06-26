FROM oven/bun

RUN apt update && \
    apt install -y \
    vim \
    curl \
    iputils-ping \
    tmux \
    nodejs \
    make \
    fish \
    docker

COPY ./src/orchestrator/tmux/tmux.conf /root/.tmux.conf

RUN bun install -g @anthropic-ai/claude-code @google/gemini-cli

WORKDIR /src

CMD [ "sleep", "infinity" ]
