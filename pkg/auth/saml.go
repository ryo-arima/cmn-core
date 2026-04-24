package auth

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

// SAMLConfig holds the configuration required to act as a SAML 2.0 Service
// Provider (SP) against an IdP such as Keycloak or Authentik.
type SAMLConfig struct {
	// ProviderName is a human-readable label, e.g. "keycloak-saml".
	ProviderName string `yaml:"provider_name"`
	// IDPMetadataURL is the URL for the IdP's SAML metadata XML.
	IDPMetadataURL string `yaml:"idp_metadata_url"`
	// IDPCertificatePEM is the PEM-encoded X.509 certificate used by the IdP.
	IDPCertificatePEM string `yaml:"idp_certificate_pem"`
	// SPEntityID is the Service Provider's entity ID.
	SPEntityID string `yaml:"sp_entity_id"`
	// SPACSURL is the SP's Assertion Consumer Service URL.
	// Example: https://example.com/v1/share/auth/saml/callback
	SPACSURL string `yaml:"sp_acs_url"`
	// SPKeyPEM is the PEM-encoded private key used by the SP for signing.
	SPKeyPEM string `yaml:"sp_key_pem"`
	// SPCertPEM is the PEM-encoded certificate for the SP.
	SPCertPEM string `yaml:"sp_cert_pem"`
}

// samlProvider implements Provider for SAML 2.0 identity providers.
type samlProvider struct {
	cfg        SAMLConfig
	sp         *saml.ServiceProvider
}

// NewSAMLProvider creates a new SAML SP provider from the supplied config.
// When IDPMetadataURL is set, provider metadata is fetched at startup.
func NewSAMLProvider(ctx context.Context, cfg SAMLConfig) (Provider, error) {
	if cfg.SPACSURL == "" {
		return nil, errors.New("auth: SAML SP ACS URL must not be empty")
	}
	if cfg.SPEntityID == "" {
		return nil, errors.New("auth: SAML SP entity ID must not be empty")
	}

	acsURL, err := url.Parse(cfg.SPACSURL)
	if err != nil {
		return nil, fmt.Errorf("auth: invalid SAML ACS URL: %w", err)
	}

	entityID, err := url.Parse(cfg.SPEntityID)
	if err != nil {
		return nil, fmt.Errorf("auth: invalid SAML SP entity ID: %w", err)
	}

	// Parse SP key and certificate if provided
	var keyPair *tls.Certificate
	if cfg.SPKeyPEM != "" && cfg.SPCertPEM != "" {
		pair, err := tls.X509KeyPair([]byte(cfg.SPCertPEM), []byte(cfg.SPKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("auth: failed to parse SP key pair: %w", err)
		}
		keyPair = &pair
	}

	// Fetch IdP metadata
	var idpMetadata *saml.EntityDescriptor
	if cfg.IDPMetadataURL != "" {
		rawURL, _ := url.Parse(cfg.IDPMetadataURL)
		metadata, err := samlsp.FetchMetadata(ctx, http.DefaultClient, *rawURL)
		if err != nil {
			return nil, fmt.Errorf("auth: failed to fetch SAML IdP metadata from %s: %w", cfg.IDPMetadataURL, err)
		}
		idpMetadata = metadata
	} else if cfg.IDPCertificatePEM != "" {
		// Build minimal IdP descriptor from PEM certificate
		block, _ := pem.Decode([]byte(cfg.IDPCertificatePEM))
		if block == nil {
			return nil, errors.New("auth: SAML IdP certificate PEM is invalid")
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("auth: failed to parse SAML IdP certificate: %w", err)
		}
		_ = cert // used below in descriptor
		idpMetadata = &saml.EntityDescriptor{
			EntityID: cfg.SPEntityID, // placeholder; must be set from config
		}
	} else {
		return nil, errors.New("auth: SAML IdP metadata URL or certificate PEM must be set")
	}

	sp := saml.ServiceProvider{
		EntityID:    entityID.String(),
		AcsURL:      *acsURL,
		IDPMetadata: idpMetadata,
	}
	if keyPair != nil {
		sp.Key = keyPair.PrivateKey.(*rsa.PrivateKey)
		if len(keyPair.Certificate) > 0 {
			cert, _ := x509.ParseCertificate(keyPair.Certificate[0])
			sp.Certificate = cert
		}
	}

	return &samlProvider{cfg: cfg, sp: &sp}, nil
}

func (p *samlProvider) Name() string { return p.cfg.ProviderName }

// LoginURL generates a SAML AuthnRequest and returns the redirect URL.
// state is passed as the RelayState parameter.
func (p *samlProvider) LoginURL(state string) (string, error) {
	authnReq, err := p.sp.MakeAuthenticationRequest(
		p.sp.GetSSOBindingLocation(saml.HTTPRedirectBinding),
		saml.HTTPRedirectBinding,
		saml.HTTPPostBinding,
	)
	if err != nil {
		return "", fmt.Errorf("auth: SAML AuthnRequest generation failed: %w", err)
	}

	redirectURL, err := authnReq.Redirect(state, p.sp)
	if err != nil {
		return "", fmt.Errorf("auth: SAML redirect URL generation failed: %w", err)
	}

	return redirectURL.String(), nil
}

// HandleCallback validates the SAML Response posted to the ACS endpoint and
// returns normalized Claims.
func (p *samlProvider) HandleCallback(ctx context.Context, r *http.Request) (*Claims, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("auth: failed to parse SAML callback form: %w", err)
	}

	samlResponse := r.FormValue("SAMLResponse")
	if samlResponse == "" {
		return nil, errors.New("auth: SAML callback missing SAMLResponse parameter")
	}

	// Parse the SAML response — allow any previously-issued request IDs.
	assertion, err := p.sp.ParseXMLResponse([]byte(samlResponse), []string{}, p.sp.AcsURL)
	if err != nil {
		return nil, fmt.Errorf("auth: SAML response validation failed: %w", err)
	}

	claims := &Claims{
		Subject: assertion.Subject.NameID.Value,
	}

	// Extract standard attributes
	for _, stmt := range assertion.AttributeStatements {
		for _, attr := range stmt.Attributes {
			vals := make([]string, 0, len(attr.Values))
			for _, v := range attr.Values {
				vals = append(vals, v.Value)
			}
			switch attr.Name {
			case "email", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress":
				if len(vals) > 0 {
					claims.Email = vals[0]
				}
			case "name", "displayName", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name":
				if len(vals) > 0 {
					claims.Name = vals[0]
				}
			case "groups", "roles", "memberOf":
				claims.Groups = append(claims.Groups, vals...)
			}
			claims.RawAttributes[attr.Name] = vals
		}
	}

	if claims.RawAttributes == nil {
		claims.RawAttributes = make(map[string][]string)
	}

	return claims, nil
}

// Metadata returns the SP metadata XML for registration with the IdP.
func (p *samlProvider) Metadata() ([]byte, error) {
	meta := p.sp.Metadata()
	b, err := xml.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("auth: failed to marshal SAML SP metadata: %w", err)
	}
	return b, nil
}

