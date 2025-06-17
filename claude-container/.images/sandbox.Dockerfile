FROM oven/bun

RUN apt update && \
    apt install -y \
    vim \
    curl \
    iputils-ping \
    tmux

WORKDIR /src

RUN bun install -g @anthropic-ai/claude-code

CMD [ "sleep", "infinity" ]
