# Mail Server Setup

This project includes a complete mail server setup using docker-mailserver and Roundcube webmail UI.

## Components

- **docker-mailserver**: Full-featured mail server (SMTP/IMAP)
- **Roundcube**: Web-based email client UI

## Quick Start

### 1. Start the Mail Server

```bash
# Start all services including mailserver
docker-compose up -d

# Wait for services to be ready
sleep 10

# Setup mail accounts
chmod +x scripts/setup-mailserver.sh
./scripts/setup-mailserver.sh
```

### 2. Access Roundcube Webmail

Open your browser and navigate to: http://localhost:3005

### 3. Default Mail Accounts

#### Administrator
| Email | Password | Purpose |
|-------|----------|---------|
| admin@cmn.local | AdminPassword123! | Administrator account |

#### Test Users
| Email | Password | Purpose |
|-------|----------|---------|
| test1@cmn.local | TestPassword123! | Test user 1 |
| test2@cmn.local | TestPassword123! | Test user 2 |
| test3@cmn.local | TestPassword123! | Test user 3 |
| test4@cmn.local | TestPassword123! | Test user 4 |
| test5@cmn.local | TestPassword123! | Test user 5 |

#### Regular Users
| Email | Password | Purpose |
|-------|----------|---------|
| user1@cmn.local | User1Password123! | Regular user 1 |
| user2@cmn.local | User2Password123! | Regular user 2 |
| developer@cmn.local | DevPassword123! | Developer account |

#### System Accounts
| Email | Password | Purpose |
|-------|----------|---------|
| noreply@cmn.local | NoReplyPassword123! | System notifications (no-reply) |
| support@cmn.local | SupportPassword123! | Support account |

**Total Accounts**: 11

## Configuration

### SMTP Settings (Sending Mail)

- **Host**: localhost
- **Port**: 587 (TLS/STARTTLS)
- **Port**: 25 (SMTP)
- **Authentication**: Required

### IMAP Settings (Receiving Mail)

- **Host**: localhost
- **Port**: 993 (SSL/TLS)
- **Port**: 143 (IMAP)
- **Authentication**: Required

## Usage in Application

### Sending Email from Go Code

```go
import "github.com/ryo-arima/cmn-core/pkg/mail"

// Create mail sender
mailConfig := mail.Config{
    Host:     "localhost",
    Port:     587,
    Username: "noreply@cmn.local",
    Password: "NoReplyPassword123!",
    From:     "noreply@cmn.local",
    UseTLS:   true,
}
sender := mail.NewSender(mailConfig)

// Send welcome email
err := sender.SendWelcomeEmail("user@cmn.local", "John Doe")
if err != nil {
    log.Printf("Failed to send email: %v", err)
}

// Send custom email
msg := mail.Message{
    To:      []string{"recipient@cmn.local"},
    Subject: "Test Email",
    Body:    "<h1>Hello World</h1>",
    IsHTML:  true,
}
err = sender.Send(msg)
```

## Managing Mail Accounts

### Add New Account

```bash
docker exec -it cmn-mailserver setup email add user@cmn.local "password"
```

### List All Accounts

```bash
docker exec -it cmn-mailserver setup email list
```

### Delete Account

```bash
docker exec -it cmn-mailserver setup email del user@cmn.local
```

### Change Password

```bash
docker exec -it cmn-mailserver setup email update user@cmn.local "newpassword"
```

## Troubleshooting

### View Mail Server Logs

```bash
docker logs cmn-mailserver
```

### View Roundcube Logs

```bash
docker logs cmn-roundcube
```

### Check Mail Queue

```bash
docker exec -it cmn-mailserver postqueue -p
```

### Test SMTP Connection

```bash
telnet localhost 587
# or using openssl for TLS
openssl s_client -connect localhost:587 -starttls smtp
```

### Test IMAP Connection

```bash
openssl s_client -connect localhost:993
```

## Ports

| Service | Port | Protocol | Description |
|---------|------|----------|-------------|
| Mailserver | 25 | SMTP | Mail transfer |
| Mailserver | 143 | IMAP | Mail retrieval |
| Mailserver | 587 | SMTP | Mail submission (TLS) |
| Mailserver | 993 | IMAPS | Secure mail retrieval |
| Roundcube | 3005 | HTTP | Webmail UI |

## Production Considerations

### Security

1. **Change default passwords** in production
2. **Use real domain names** instead of `.local`
3. **Configure SPF, DKIM, and DMARC** records
4. **Enable SSL certificates** (Let's Encrypt)
5. **Configure firewall rules**

### DNS Configuration

For production, configure these DNS records:

```
MX     @     10 mail.yourdomain.com
A      mail  <your-server-ip>
TXT    @     "v=spf1 mx ~all"
```

### Volume Persistence

Mail data is stored in Docker volumes:

- `mailserver-data`: Mail files
- `mailserver-state`: Server state
- `mailserver-logs`: Log files
- `mailserver-config`: Configuration

Backup these volumes regularly in production.

## References

- [docker-mailserver Documentation](https://docker-mailserver.github.io/docker-mailserver/latest/)
- [Roundcube Documentation](https://github.com/roundcube/roundcubemail/wiki)
- [Zenn Article: Docker Mailserver Setup](https://zenn.dev/takaha/articles/docker-mailserver)
