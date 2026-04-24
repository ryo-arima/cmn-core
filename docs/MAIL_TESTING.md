# メールの送受信テストガイド

Roundcube Webメールクライアントを使用して、ブラウザからメールの送受信をテストできます。

## 前提条件

メールサーバーとRoundcubeが起動している必要があります。

```bash
# mailserver と roundcube を強制再作成して起動（設定変更を確実に反映）
docker compose up -d --force-recreate mailserver roundcube

# メールアカウントをセットアップ（未作成の場合）
python3 scripts/mail/generate_mail_accounts.py
docker compose exec mailserver setup email list || true
```

## Roundcube Webメールへのアクセス

### 1. ブラウザでアクセス

http://localhost:3005 をブラウザで開く

### 2. ログイン

以下のいずれかのアカウントでログインできます：

#### テストユーザー
- **ユーザー名**: `test1@cmn.local`
- **パスワード**: `TestPassword123!`

その他のアカウント：
- `test2@cmn.local` / `TestPassword123!`
- `test3@cmn.local` / `TestPassword123!`
- `test4@cmn.local` / `TestPassword123!`
- `test5@cmn.local` / `TestPassword123!`
- `user1@cmn.local` / `User1Password123!`
- `user2@cmn.local` / `User2Password123!`
- `admin@cmn.local` / `AdminPassword123!`

## メール送受信テスト手順

### テスト1: 内部ユーザー間のメール送信

1. **送信者でログイン**
   - ブラウザで http://localhost:3005 を開く
   - `test1@cmn.local` / `TestPassword123!` でログイン

2. **メール作成**
   - 左上の「作成」または「Compose」ボタンをクリック
   - 宛先: `test2@cmn.local`
   - 件名: `テストメール`
   - 本文: `これはテストメールです。`
   - 「送信」をクリック

3. **受信者で確認**
   - 新しいブラウザタブまたはシークレットウィンドウで http://localhost:3005 を開く
   - `test2@cmn.local` / `TestPassword123!` でログイン
   - 受信トレイにメールが届いているか確認

### テスト2: 複数の受信者へ送信

1. `test1@cmn.local` でログイン
2. メール作成
   - 宛先: `test2@cmn.local, test3@cmn.local`
   - CC: `user1@cmn.local`
   - BCC: `user2@cmn.local`
   - 件名と本文を入力して送信

3. 各受信者のアカウントでログインして確認

### テスト3: 添付ファイル付きメール

1. `user1@cmn.local` でログイン
2. メール作成
   - 宛先: `user2@cmn.local`
   - 添付ファイルを追加（クリップアイコン）
   - 小さなテキストファイルや画像をアップロード
   - 送信

3. `user2@cmn.local` で受信確認し、添付ファイルをダウンロード

### テスト4: 返信とフォワード

1. `test2@cmn.local` でログイン
2. `test1@cmn.local` から受信したメールを開く
3. 「返信」をクリックして返信メールを送信
4. または「転送」をクリックして別のユーザーに転送

### テスト5: システムメール送信のテスト（Go言語から）

アプリケーションからメールを送信：

```bash
# テスト用のGoスクリプトを実行
cd /Users/ryo/mysrc/github/ryo-arima/cmn-core
```

```go
// test/manual/send_test_email.go
package main

import (
    "log"
    "github.com/ryo-arima/cmn-core/pkg/mail"
)

func main() {
    // メール設定
    config := mail.Config{
        Host:     "localhost",
        Port:     587,
        Username: "noreply@cmn.local",
        Password: "NoReplyPassword123!",
        From:     "noreply@cmn.local",
        UseTLS:   true,
    }
    
    sender := mail.NewSender(config)
    
    // Welcomeメール送信
    err := sender.SendWelcomeEmail("test1@cmn.local", "Test User 1")
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }
    
    log.Println("Email sent successfully!")
    log.Println("Check http://localhost:3005 and login as test1@cmn.local")
}
```

実行後、`test1@cmn.local` でログインしてメールを確認

## トラブルシューティング

### メールが届かない場合

1. **コンテナの状態確認**
   ```bash
   docker ps | grep -E "mailserver|roundcube"
   ```
   両方のコンテナが起動していることを確認

2. **メールサーバーログ確認**
   ```bash
   docker logs cmn-mailserver -f
   ```

3. **DNS解決確認**
   ```bash
   docker exec -it cmn-mailserver nslookup mail.cmn.local 172.20.0.2
   ```

4. **メールアカウント確認**
   ```bash
   docker exec -it cmn-mailserver setup email list
   ```

### Roundcubeにログインできない場合

1. **Roundcubeログ確認**
   ```bash
   docker logs cmn-roundcube -f
   ```

2. **データベース接続確認**
   - Roundcubeは初回起動時にデータベースを自動セットアップ
   - MySQLが起動していることを確認

3. **ブラウザのキャッシュクリア**
   - ブラウザのキャッシュとCookieをクリアして再試行

### メール送信はできるが受信できない場合

1. **IMAP接続確認**
   ```bash
   docker exec -it cmn-mailserver telnet localhost 143
   ```

2. **メールボックス確認**
   ```bash
   docker exec -it cmn-mailserver ls -la /var/mail
   ```

## SMTP/IMAPクライアント設定（参考）

Roundcube以外のメールクライアント（Thunderbird、Apple Mail等）を使用する場合：

### IMAP設定（受信）
- **サーバー**: localhost または mail.cmn.local
- **ポート**: 993 (SSL/TLS)
- **セキュリティ**: SSL/TLS
- **認証**: 通常のパスワード認証

### SMTP設定（送信）
- **サーバー**: localhost または mail.cmn.local
- **ポート**: 587 (STARTTLS)
- **セキュリティ**: STARTTLS
- **認証**: 通常のパスワード認証

## よくあるテストシナリオ

### シナリオ1: ユーザー登録時のWelcomeメール

1. アプリケーションで新規ユーザー登録
2. 登録したメールアドレスでRoundcubeにログイン
3. Welcomeメールが届いていることを確認

### シナリオ2: パスワードリセット

1. パスワードリセットリクエストを送信
2. 該当ユーザーのメールボックスで確認
3. リセットリンクをクリックして動作確認

### シナリオ3: 通知メール

1. システムイベント発生時の通知メール送信
2. 複数ユーザーへの一斉通知
3. 管理者への警告メール

## クリーンアップ

テスト後、メールデータをクリアしたい場合：

```bash
# 推奨: 強制再作成で設定変更やクリーン状態を確実に反映
docker compose up -d --force-recreate mailserver roundcube

# さらに完全にクリーンにする場合
docker compose down
docker compose up -d --force-recreate mailserver roundcube
python3 scripts/mail/generate_mail_accounts.py
docker compose exec mailserver setup email list || true
```

## 参考リンク

- Roundcube公式: https://roundcube.net/
- docker-mailserver: https://docker-mailserver.github.io/docker-mailserver/
- 内部ドキュメント: 
  - [MAIL_SERVER.md](./MAIL_SERVER.md) - メールサーバー設定
  - [DNS_SETUP.md](./DNS_SETUP.md) - DNS設定
