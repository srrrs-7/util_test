.PHONY: run setup install update freeze dev test clean

run:
	uv run ./main.py

# 基本的な環境セットアップ
setup:
	pip install uv

# 依存関係のインストール
install:
	uv pip install -r requirements.txt

# 開発用依存関係のインストール
dev-install:
	uv pip install -r requirements-dev.txt

# すべてのパッケージを最新バージョンに更新
update:
	uv pip install --upgrade -r requirements.txt

# 現在の環境を requirements.txt に出力
freeze:
	uv pip freeze > requirements.txt

# 開発サーバーの起動
dev:
	uv run python -m your_app

# テストの実行
test:
	uv run pytest

# キャッシュやビルドファイルの削除
clean:
	rm -rf .uv
	rm -rf __pycache__
	rm -rf *.egg-info
	rm -rf .pytest_cache
	rm -rf dist
	rm -rf build

# 仮想環境の作成
venv:
	uv venv

# 特定のパッケージのインストール (make pkg PKG=package_name)
pkg:
	uv pip install $(PKG)

# 依存関係の監査
audit:
	uv pip audit