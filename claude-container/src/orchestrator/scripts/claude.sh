#!/bin/bash

tmux split-window -h && \
    tmux split-window -v && \
    tmux select-pane -t 0 && \
    tmux split-window -v && \
    tmux select-pane -t 2 && \
    tmux split-window -v

tmux list-panes -F "#{pane_id}"

tmux list-panes -F "#{pane_id}" | xargs -I {} sh -c 'tmux send-keys -t {} "claude"; sleep 1; tmux send-keys -t {} Enter'
