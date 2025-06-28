# task

リッチなECサイトの画面を作成して

## context

./pjt配下で作業してください
作業する際は、%1, %2, %3, %4, %5のpaneにタスクを割り振って実行させてください
また、指示を出す際は、必ず'ultrathink'を付与して指示してください
taskを各paneに割り当て終わったら、下記コマンドを実行してください

```bash
tmux send-keys -t %1 Enter & \
tmux send-keys -t %2 Enter & \
tmux send-keys -t %3 Enter & \
tmux send-keys -t %4 Enter & \
tmux send-keys -t %5 Enter &
```

各paneはタスクの進捗や疑問点、完了を都度報告するように指示してください

## logging

./log/にlogファイルに書き出しておいてください
証跡として使用したいので詳細に記述をお願いします
