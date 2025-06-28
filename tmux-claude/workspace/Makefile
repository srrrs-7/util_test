.PHONY: claude claude-tmux exec kill

claude:
	claude --dangerously-skip-permissions

claude-tmux:
	/scripts/tmux-claude.sh

exec:
	tmux send-keys -t %1 Enter & \
	tmux send-keys -t %2 Enter & \
	tmux send-keys -t %3 Enter & \
	tmux send-keys -t %4 Enter & \
	tmux send-keys -t %5 Enter &

kill:
	tmux kill-session