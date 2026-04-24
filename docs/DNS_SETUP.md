# ローカルDNS設定ガイド

## 概要

cmn-coreプロジェクトでは、メールサーバーのドメイン名解決のために、ローカルDNSサーバー（dnsmasq）を使用しています。

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────┐
│                  cmn-network (172.20.0.0/16)          │
│                                                         │
│  ┌──────────────┐                                       │
│  │ DNS Server   │  172.20.0.2                           │
│  │ (dnsmasq)    │                                       │
│  └──────┬───────┘                                       │
│         │                                               │
│         ├────────> 他の全サービスがこのDNSを参照             │
│         │                                               │
│  ┌──────▼───────┐        ┌──────────────┐              │
│  │ Mail Server  │        │ Roundcube    │              │
│  │ 172.20.0.10  │◄───────┤ 172.20.0.11  │              │
│  └──────────────┘        └──────────────┘              │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

## DNS レコード

以下のドメインが自動的に解決されます：

| ドメイン | IPアドレス | 用途 |
|---------|-----------|------|
| `cmn.local` | 172.20.0.10 | メインドメイン |
| `mail.cmn.local` | 172.20.0.10 | メールサーバー |
| `roundcube.cmn.local` | 172.20.0.11 | Webメールクライアント |

## DNSサーバー管理

### Web UI アクセス

DNSサーバーの設定はWeb UIで確認できます：

- URL: http://localhost:8053
- ユーザー名: `admin`
- パスワード: `admin`

### DNS動作確認

コンテナ内からDNS解決をテストする：

```bash
# mailserverコンテナ内でテスト
docker exec -it cmn-mailserver nslookup mail.cmn.local 172.20.0.2

# ホストマシンからテスト（dockerの内部DNSを使用）
docker run --rm --network cmn_cmn-network alpine nslookup mail.cmn.local 172.20.0.2
```

### ログ確認

DNSクエリのログを確認：

```bash
docker logs cmn-dns -f
```

## ホストマシンでのDNS設定（オプション）

ホストマシンからも `*.cmn.local` にアクセスしたい場合は、以下の設定を行います。

### macOS

`/etc/hosts` に追加：

```bash
sudo nano /etc/hosts
```

以下を追加：

```
127.0.0.1 mail.cmn.local
127.0.0.1 roundcube.cmn.local
```

または、macOSのリゾルバー設定：

```bash
sudo mkdir -p /etc/resolver
sudo bash -c 'echo "nameserver 127.0.0.1" > /etc/resolver/cmn.local'
```

### Linux

`/etc/hosts` に追加：

```bash
sudo nano /etc/hosts
```

以下を追加：

```
127.0.0.1 mail.cmn.local
127.0.0.1 roundcube.cmn.local
```

### Windows

`C:\Windows\System32\drivers\etc\hosts` を管理者権限で編集し、以下を追加：

```
127.0.0.1 mail.cmn.local
127.0.0.1 roundcube.cmn.local
```

## トラブルシューティング

### DNSが解決されない場合

1. DNSサーバーが起動しているか確認：
   ```bash
   docker ps | grep cmn-dns
   ```

2. DNSサーバーのログを確認：
   ```bash
   docker logs cmn-dns
   ```

3. コンテナのDNS設定を確認：
   ```bash
   docker exec -it cmn-mailserver cat /etc/resolv.conf
   ```

   `nameserver 172.20.0.2` が含まれていることを確認

### メールサーバーがドメインを解決できない場合

mailserverコンテナ内でDNS解決をテスト：

```bash
docker exec -it cmn-mailserver bash
# コンテナ内で
nslookup mail.cmn.local
nslookup cmn.local
dig mail.cmn.local
```

### ネットワークの問題

ネットワークの再作成：

```bash
docker-compose down
docker network prune
docker-compose up -d
```

## 技術詳細

### dnsmasq設定

DNSサーバーは以下のコマンドで起動されます：

```
--address=/cmn.local/172.20.0.10
--address=/mail.cmn.local/172.20.0.10
--address=/roundcube.cmn.local/172.20.0.11
--server=8.8.8.8
--server=8.8.4.4
--log-queries
--no-resolv
```

- `--address`: ドメインとIPアドレスのマッピング
- `--server`: 上位DNSサーバー（外部ドメイン解決用）
- `--log-queries`: クエリログを有効化
- `--no-resolv`: `/etc/resolv.conf` を読まない

### ネットワーク設定

- ネットワーク名: `cmn-network`
- サブネット: `172.20.0.0/16`
- ゲートウェイ: `172.20.0.1`

### 静的IPアドレス

| サービス | IPアドレス |
|---------|-----------|
| DNS | 172.20.0.2 |
| Mail Server | 172.20.0.10 |
| Roundcube | 172.20.0.11 |
| その他 | DHCP（172.20.x.x） |

## セキュリティ考慮事項

- DNSサーバーはローカル開発環境専用です
- 本番環境では適切なDNSサーバー（Route53、CloudDNS等）を使用してください
- 自己署名証明書を使用しているため、本番環境では正式な証明書を取得してください

## 参考リンク

- [dnsmasq公式ドキュメント](http://www.thekelleys.org.uk/dnsmasq/doc.html)
- [docker-mailserver ドキュメント](https://docker-mailserver.github.io/docker-mailserver/latest/)
