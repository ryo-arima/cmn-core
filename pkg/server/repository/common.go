package repository

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/model"
	"github.com/ryo-arima/cmn-core/pkg/server/share"
	"golang.org/x/crypto/bcrypt"
)

// Common interface for repository layer
type Common interface {
	share.Common // Embed share.Common for middleware compatibility
	GetBaseConfig() config.BaseConfig
	GenerateJWTToken(claims model.JWTClaims) (string, error)
	InvalidateToken(ctx context.Context, tokenString string) error
	GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error)
	GenerateJWTSecret() (string, error)
	ValidateJWTSecretStrength(secret string) error
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	ValidatePasswordStrength(password string) error
	DeleteTokenCache(token string)
	SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error
	// SSO state management (OIDC / SAML polling flow)
	StoreAuthState(ctx context.Context, state string) error
	ConsumeAuthState(ctx context.Context, state string) error
	StoreTokenForState(ctx context.Context, state string, tokenPair *model.TokenPair) error
	GetTokenForState(ctx context.Context, state string) (*model.TokenPair, bool, error)
}

// common implements Common interface
type common struct {
	BaseConfig  config.BaseConfig
	RedisClient *redis.Client
	MailConfig  *config.Mail
}

func (rcvr *common) GetBaseConfig() config.BaseConfig {
	return rcvr.BaseConfig
}

// GetJWTSecret returns the JWT secret key from config, environment variable, or a default value
func (rcvr *common) getJWTSecret() string {
	// First try environment variable
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		return envSecret
	}

	// Then try config file
	if rcvr.BaseConfig.YamlConfig.Application.Server.JWTSecret != "" {
		return rcvr.BaseConfig.YamlConfig.Application.Server.JWTSecret
	}

	// Finally, use default (should be changed in production)
	return "your-256-bit-secret-key-change-this-in-production"
}

// GenerateJWTToken creates a JWT token with the given claims
func (rcvr *common) GenerateJWTToken(claims model.JWTClaims) (string, error) {
	// Create header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(headerJSON)

	// Create payload
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(payloadJSON)

	// Create signature
	message := headerEncoded + "." + payloadEncoded
	signature := rcvr.createSignature(message)

	// Combine all parts
	token := message + "." + signature
	return token, nil
}

// ValidateJWTToken validates and parses a JWT token
func (rcvr *common) ValidateJWTToken(tokenString string) (*model.JWTClaims, error) {
	// 1. Try cache first if enabled
	if rcvr.RedisClient != nil && rcvr.BaseConfig.YamlConfig.Application.Server.Redis.JWTCache {
		if cached, err := rcvr.getCachedTokenClaims(tokenString); err == nil && cached != nil {
			// Ensure not expired
			if cached.ExpiresAt >= time.Now().Unix() {
				return cached, nil
			}
		}
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Verify signature
	message := parts[0] + "." + parts[1]
	expectedSignature := rcvr.createSignature(message)
	if parts[2] != expectedSignature {
		return nil, errors.New("invalid token signature")
	}

	// Decode payload
	payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims model.JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	now := time.Now().Unix()
	if claims.ExpiresAt < now {
		return nil, errors.New("token expired")
	}

	// Cache claims if enabled
	if rcvr.RedisClient != nil && rcvr.BaseConfig.YamlConfig.Application.Server.Redis.JWTCache {
		_ = rcvr.cacheTokenClaims(tokenString, &claims)
	}

	return &claims, nil
}

// GenerateTokenPair creates both access and refresh tokens
func (rcvr *common) GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error) {
	now := time.Now()
	accessTokenExpiry := now.Add(24 * time.Hour).Unix()      // 24 hours
	refreshTokenExpiry := now.Add(7 * 24 * time.Hour).Unix() // 7 days

	// Generate a unique ID for the token
	jti := uuid.New().String()

	// Create access token claims
	accessClaims := model.JWTClaims{
		Jti:       jti,
		UserID:    userID,
		UUID:      userUUID,
		Email:     email,
		Name:      name,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: accessTokenExpiry,
	}

	// Generate access token
	accessToken, err := rcvr.GenerateJWTToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create refresh token claims (longer expiry, different JTI)
	refreshJti := uuid.New().String()
	refreshClaims := model.JWTClaims{
		Jti:       refreshJti,
		UserID:    userID,
		UUID:      userUUID,
		Email:     email,
		Name:      name,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: refreshTokenExpiry,
	}

	// Generate refresh token
	refreshToken, err := rcvr.GenerateJWTToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessTokenExpiry - now.Unix(),
	}, nil
}

// createSignature creates HMAC-SHA256 signature for JWT
func (rcvr *common) createSignature(message string) string {
	secret := rcvr.getJWTSecret()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
	return signature
}

// HashPassword hashes a password using bcrypt
func (rcvr *common) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against its hash
func (rcvr *common) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength validates password strength
func (rcvr *common) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// GenerateJWTSecret generates a secure random JWT secret
func (rcvr *common) GenerateJWTSecret() (string, error) {
	// Generate 32 bytes (256 bits) of random data
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 for use as a string
	secret := base64.URLEncoding.EncodeToString(bytes)
	return secret, nil
}

// ValidateJWTSecretStrength checks if JWT secret meets minimum security requirements
func (rcvr *common) ValidateJWTSecretStrength(secret string) error {
	if len(secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	// Check for common weak secrets
	weakSecrets := []string{
		"secret",
		"jwt-secret",
		"your-secret-key",
		"change-me",
		"your-256-bit-secret-key-change-this-in-production",
	}

	for _, weak := range weakSecrets {
		if secret == weak {
			return fmt.Errorf("JWT secret is too weak: %s", secret)
		}
	}

	return nil
}

// InvalidateToken adds a token's JTI to the Redis denylist
func (rcvr *common) InvalidateToken(ctx context.Context, tokenString string) error {
	claims, err := rcvr.ValidateJWTToken(tokenString)
	if err != nil {
		// If token is already expired or invalid, we don't need to do anything.
		// We can consider it "successfully" invalidated.
		if strings.Contains(err.Error(), "token expired") {
			return nil
		}
		return fmt.Errorf("error validating token before invalidation: %w", err)
	}

	// Calculate remaining time until expiration
	now := time.Now()
	expiresAt := time.Unix(claims.ExpiresAt, 0)
	if now.After(expiresAt) {
		// Already expired
		return nil
	}
	ttl := expiresAt.Sub(now)

	// Add to denylist in Redis
	err = rcvr.RedisClient.Set(ctx, claims.Jti, "invalidated", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to denylist: %w", err)
	}

	return nil
}

// ParseTokenUnverified decodes the claims from a token without verifying its signature.
// This is used to get the JTI for denylist checking before full validation.
func (rcvr *common) ParseTokenUnverified(tokenString string) (*model.JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims model.JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	return &claims, nil
}

// IsTokenInvalidated checks if a token's JTI exists in the Redis denylist.
func (rcvr *common) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	if rcvr.RedisClient == nil {
		return false, nil
	}
	result, err := rcvr.RedisClient.Exists(ctx, jti).Result()
	if err != nil {
		return true, fmt.Errorf("error checking token in redis: %w", err)
	}
	return result == 1, nil
}

// helper: cache key builder
func (rcvr *common) tokenCacheKey(token string) string {
	return "auth:token:" + token
}

// helper: store token claims in redis with configurable TTL
func (rcvr *common) cacheTokenClaims(token string, claims *model.JWTClaims) error {
	if rcvr.RedisClient == nil || claims == nil {
		return nil
	}
	data, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	remaining := time.Until(time.Unix(claims.ExpiresAt, 0))
	if remaining <= 0 {
		return nil
	}
	// Use configured TTL or default to 30 minutes
	maxTTL := 30 * time.Minute
	if rcvr.BaseConfig.YamlConfig.Application.Server.Redis.CacheTTL > 0 {
		maxTTL = time.Duration(rcvr.BaseConfig.YamlConfig.Application.Server.Redis.CacheTTL) * time.Second
	}
	// Don't cache longer than token expiry
	if remaining < maxTTL {
		maxTTL = remaining
	}
	return rcvr.RedisClient.Set(context.Background(), rcvr.tokenCacheKey(token), string(data), maxTTL).Err()
}

// helper: get token claims from redis cache
func (rcvr *common) getCachedTokenClaims(token string) (*model.JWTClaims, error) {
	if rcvr.RedisClient == nil {
		return nil, errors.New("redis client nil")
	}
	val, err := rcvr.RedisClient.Get(context.Background(), rcvr.tokenCacheKey(token)).Result()
	if err != nil {
		return nil, err
	}
	var claims model.JWTClaims
	if err := json.Unmarshal([]byte(val), &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// DeleteTokenCache removes cached claims for the given raw token string (if present)
func (rcvr *common) DeleteTokenCache(token string) {
	if rcvr.RedisClient == nil || token == "" {
		return
	}
	_ = rcvr.RedisClient.Del(context.Background(), rcvr.tokenCacheKey(token)).Err()
}

// SendEmail sends an email using the configured mail settings
func (rcvr *common) SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error {
	if rcvr.MailConfig == nil || rcvr.MailConfig.Host == "" {
		return errors.New("mail sender not configured")
	}

	auth := smtp.PlainAuth("", rcvr.MailConfig.Username, rcvr.MailConfig.Password, rcvr.MailConfig.Host)

	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: %s; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		rcvr.MailConfig.From, to, subject, contentType, body)

	addr := fmt.Sprintf("%s:%d", rcvr.MailConfig.Host, rcvr.MailConfig.Port)
	return smtp.SendMail(addr, auth, rcvr.MailConfig.From, []string{to}, []byte(message))
}

// SendWelcomeEmail sends a welcome email to a new user
func (rcvr *common) SendWelcomeEmail(ctx context.Context, to, name string) error {
	if rcvr.MailConfig == nil || rcvr.MailConfig.Host == "" {
		return errors.New("mail sender not configured")
	}

	subject := "Welcome to cmn-core!"
	body := fmt.Sprintf("Hello %s,\n\nWelcome to cmn-core! Your account has been created successfully.\n\nBest regards,\nThe cmn-core Team", name)

	return rcvr.SendEmail(ctx, to, subject, body, false)
}

// SendPasswordResetEmail sends a password reset email to a user
func (rcvr *common) SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error {
	if rcvr.MailConfig == nil || rcvr.MailConfig.Host == "" {
		return errors.New("mail sender not configured")
	}

	subject := "Password Reset Request"
	body := fmt.Sprintf("Hello %s,\n\nYou have requested to reset your password. Please click the link below to reset your password:\n\n%s\n\nIf you did not request this, please ignore this email.\n\nBest regards,\nThe cmn-core Team", name, resetURL)

	return rcvr.SendEmail(ctx, to, subject, body, false)
}

func NewCommon(baseConfig config.BaseConfig, redisClient *redis.Client) Common {
	// Initialize mail config reference from base config
	var mailConfig *config.Mail
	if baseConfig.YamlConfig.Application.Mail.Host != "" {
		mailConfig = &baseConfig.YamlConfig.Application.Mail
	}

	return &common{
		BaseConfig:  baseConfig,
		RedisClient: redisClient,
		MailConfig:  mailConfig,
	}
}

// ssoStateKey returns the Redis key for an OIDC/SAML auth state value.
func ssoStateKey(state string) string { return "sso:state:" + state }

// ssoTokenKey returns the Redis key where the token is stored after callback.
func ssoTokenKey(state string) string { return "sso:token:" + state }

const ssoStateTTL = 10 * time.Minute
const ssoTokenTTL = 5 * time.Minute

// StoreAuthState stores a one-time state value in Redis to guard against CSRF.
func (rcvr *common) StoreAuthState(ctx context.Context, state string) error {
	if rcvr.RedisClient == nil {
		return errors.New("redis client not available")
	}
	return rcvr.RedisClient.Set(ctx, ssoStateKey(state), "1", ssoStateTTL).Err()
}

// ConsumeAuthState verifies and atomically removes the state value from Redis.
// Returns an error if the state does not exist (expired or unknown).
func (rcvr *common) ConsumeAuthState(ctx context.Context, state string) error {
	if rcvr.RedisClient == nil {
		return errors.New("redis client not available")
	}
	n, err := rcvr.RedisClient.Del(ctx, ssoStateKey(state)).Result()
	if err != nil {
		return fmt.Errorf("failed to consume auth state: %w", err)
	}
	if n == 0 {
		return errors.New("auth state not found or expired")
	}
	return nil
}

// StoreTokenForState stores the issued token pair in Redis after a successful SSO callback.
// The CLI polls for this token using the same state value.
func (rcvr *common) StoreTokenForState(ctx context.Context, state string, tokenPair *model.TokenPair) error {
	if rcvr.RedisClient == nil {
		return errors.New("redis client not available")
	}
	data, err := json.Marshal(tokenPair)
	if err != nil {
		return fmt.Errorf("failed to marshal token pair: %w", err)
	}
	return rcvr.RedisClient.Set(ctx, ssoTokenKey(state), string(data), ssoTokenTTL).Err()
}

// GetTokenForState retrieves the token pair stored after SSO callback.
// Returns (pair, true, nil) if found, (nil, false, nil) if not yet available.
func (rcvr *common) GetTokenForState(ctx context.Context, state string) (*model.TokenPair, bool, error) {
	if rcvr.RedisClient == nil {
		return nil, false, errors.New("redis client not available")
	}
	val, err := rcvr.RedisClient.GetDel(ctx, ssoTokenKey(state)).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to get token for state: %w", err)
	}
	var pair model.TokenPair
	if err := json.Unmarshal([]byte(val), &pair); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal token pair: %w", err)
	}
	return &pair, true, nil
}
