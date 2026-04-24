# Scripts

このディレクトリには、cmn-coreプロジェクトで使用する各種スクリプトが整理されています。

## 統合CLIツール (main.sh)

全スクリプト機能を統合したCLIツールを提供しています。

### 使用方法

```bash
# ヘルプを表示
./scripts/main.sh help

# 環境管理
./scripts/main.sh env recreate    # 環境を完全再構築（ephemeral）
./scripts/main.sh env up           # サービス起動
./scripts/main.sh env down         # サービス停止
./scripts/main.sh env status       # ステータス確認
./scripts/main.sh env logs [service]  # ログ表示

# メール管理
./scripts/main.sh mail create_accounts     # メールアカウント作成
./scripts/main.sh mail generate_accounts   # postfix-accounts.cf生成

# CI/テスト
./scripts/main.sh ci run [auth|mailflow]   # テストスイート実行
./scripts/main.sh ci matrix                # 全テスト実行

# ドキュメント
./scripts/main.sh docs build       # GitHub Pagesドキュメント生成
./scripts/main.sh docs serve       # ローカルサーバー起動

# DNS検証
./scripts/main.sh dns check        # DNS設定確認
```

## ディレクトリ構造

```
scripts/
├── main.sh           # 統合CLIツール（推奨）
├── lib/              # 機能別ライブラリ
│   ├── common.sh     # 共通ユーティリティ（色付きログ関数）
│   ├── env.sh        # 環境管理機能
│   ├── mail.sh       # メール管理機能
│   ├── ci.sh         # CI/テスト機能
│   ├── docs.sh       # ドキュメント機能
│   ├── dns.sh        # DNS検証機能
│   └── help.sh       # ヘルプ表示
├── data/             # データファイル
│   ├── mysql/        # MySQL初期化SQL
│   └── mailserver/   # メールサーバー設定
├── ci/               # CI/テストスクリプト
│   └── parallel.sh
├── mail/             # メール関連スクリプト
│   └── accounts.sh
├── setup/            # セットアップスクリプト
│   ├── container.sh          # 環境完全再構築（スタンドアロン版）
│   ├── docs.sh               # ドキュメントビルド
│   ├── dns.sh                # DNS検証
│   ├── architecture.sh       # アーキテクチャ図生成
│   └── secrets.sh            # LocalStack secrets初期化
└── README.md         # このファイル
```

## lib/

main.shから呼び出される機能別ライブラリ

**共通:**
- `common.sh` - ユーティリティ関数（info, success, warn, err）

**機能モジュール:**
- `env.sh` - 環境管理（recreate, up, down, status, logs）
- `mail.sh` - メール管理（create_accounts, generate_accounts）
- `ci.sh` - CI/テスト実行（run, matrix）
- `docs.sh` - ドキュメント生成（build, serve）
- `dns.sh` - DNS検証（check）
- `help.sh` - ヘルプメッセージ表示

各ライブラリは独立して動作可能ですが、main.shからの呼び出しを推奨します。

## data/

コンテナで使用する設定・データファイル

- `mysql/roundcube-init.sql` - Roundcube DB初期化
- `mailserver/postfix-accounts.cf` - メールアカウント定義
- `mailserver/postfix-main.cf` - Postfix設定オーバーライド
- `mailserver/postfix-master.cf` - Postfix master設定
- `mailserver/dovecot-quotas.cf` - Dovecot quota設定

## ci/

CI/テスト用スクリプト

- `parallel.sh` - docker-compose並列起動＆テスト実行

## mail/

メール関連スクリプト

- `accounts.sh` - メールアカウント作成（コンテナ実行中）
  - `./scripts/mail/accounts.sh` で直接実行可能
  - main.sh の `mail create_accounts` コマンドからも呼び出し

## setup/

セットアップ・初期化スクリプト

**主要スクリプト:**
- `container.sh` - 環境完全再構築（スタンドアロン版、main.sh経由推奨）
- `docs.sh` - GitHub Pagesドキュメント生成
- `dns.sh` - DNS設定検証

**補助スクリプト:**
- `architecture.sh` - アーキテクチャ図生成（mmdc必須）
- `secrets.sh` - LocalStack Secrets Manager初期化

すべてスタンドアロン実行可能ですが、main.shに統合されているものはCLI経由を推奨します。
