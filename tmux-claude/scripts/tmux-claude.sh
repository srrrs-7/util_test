#!/bin/sh

# 5つのpaneに分割
tmux split-window -h && \
    tmux split-window -v && \
    tmux select-pane -t 0 && \
    tmux split-window -v && \
    tmux select-pane -t 2 && \
    tmux split-window -v && \
    tmux select-pane -t 4 && \
    tmux split-window -v

# pane info
tmux list-panes -F "#{pane_index}: #{pane_id} #{pane_current_command} #{pane_active}"
tmux display-panes

# claude code session
tmux send-keys -t %0 "claude --dangerously-skip-permissions" && sleep 0.1 && tmux send-keys -t %0 Enter &  \
    tmux send-keys -t %1 "claude --dangerously-skip-permissions" && sleep 0.1 && tmux send-keys -t %1 Enter & \
    tmux send-keys -t %2 "claude --dangerously-skip-permissions" && sleep 0.1 && tmux send-keys -t %2 Enter & \
    tmux send-keys -t %3 "claude --dangerously-skip-permissions" && sleep 0.1 && tmux send-keys -t %3 Enter & \
    tmux send-keys -t %4 "claude --dangerously-skip-permissions" && sleep 0.1 && tmux send-keys -t %4 Enter & \
    tmux send-keys -t %5 "claude --dangerously-skip-permissions" && sleep 0.1 && tmux send-keys -t %5 Enter & \
    wait
