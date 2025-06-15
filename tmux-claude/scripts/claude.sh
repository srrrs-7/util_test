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
tmux send-keys -t %52 "claude" && sleep 0.1 && tmux send-keys -t %52 Enter & \
    tmux send-keys -t %53 "claude" && sleep 0.1 && tmux send-keys -t %53 Enter & \
    tmux send-keys -t %54 "claude" && sleep 0.1 && tmux send-keys -t %54 2Enter & \
    tmux send-keys -t %55 "claude" && sleep 0.1 && tmux send-keys -t %55 Enter & \
    tmux send-keys -t %56 "claude" && sleep 0.1 && tmux send-keys -t %56 Enter & \
    wait

# tasks
tmux send-keys -t %53 "GCP app engineについて詳細にまとめて、obsideanに書き込んで" && sleep 0.1 && tmux send-keys -t %53 Enter & \
    tmux send-keys -t %55 "AWS ECSとGCP app engineとの違いについて調査して、obsideanに書き込んで" && sleep 0.1 && tmux send-keys -t %55 Enter & \
    tmux send-keys -t %56 "AWSとGCPの価格帯系についてまとめて、安い方を教えて" && sleep 0.1 && tmux send-keys -t %56 Enter & \
    wait
