#!/bin/bash

set -e

# 出力ディレクトリ
OUT_DIR="./out"
mkdir -p "$OUT_DIR"

# 引数チェック
if [ $# -lt 3 ]; then
    echo "Usage: $0 <target_dir> <target_file_path> <target_func>"
    exit 1
fi

TARGET_DIR="$1"
TARGET_FILE_PATH="$2"
TARGET_FUNC="$3"

echo "Target directory: $TARGET_DIR"
echo "Searching for function $TARGET_FUNC in $TARGET_FILE_PATH"

# 正規表現
REGEX_STATIC_METHOD="[A-Za-z0-9_\\]+::[A-Za-z0-9_]+"
REGEX_METHOD="[A-Za-z0-9_\\]+->[A-Za-z0-9_]+"

# 関数抽出
TARGET=$(awk "/$TARGET_FUNC *\(/, /^\}/" "$TARGET_FILE_PATH")
if [ -z "$TARGET" ]; then
    echo "Function $TARGET_FUNC not found in $TARGET_FILE_PATH"
    exit 1
fi
echo "Function $TARGET_FUNC found in $TARGET_FILE_PATH"

echo "$TARGET" | grep -oE -e "$REGEX_STATIC_METHOD" -e "$REGEX_METHOD" | while read -r method; do
    echo $method
done > "$OUT_DIR/method.txt"