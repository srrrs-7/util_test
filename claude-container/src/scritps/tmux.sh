#!/bin/sh

tmux split-window -h && \
    tmux split-window -v && \
    tmux select-pane -t 0 && \
    tmux split-window -v && \
    tmux select-pane -t 2 && \
    tmux split-window -v && \