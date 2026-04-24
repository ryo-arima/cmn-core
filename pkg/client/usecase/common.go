package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/auth"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// Common is the business-logic interface for token-related operations.
// Token management is handled transparently by the auth.Manager.
type Common interface {
	Logout() response.Commons
	ValidateToken() response.ValidateToken
	GetUserInfo() response.Commons
}

type common struct {
	repo repository.Common
}

// NewCommon creates a Common usecase.
// manager is used to obtain and inject auth tokens automatically.
func NewCommon(conf config.BaseConfig, manager *clientauth.Manager) Common {
	return &common{repo: repository.NewCommon(conf, manager)}
}

func (u *common) Logout() response.Commons            { return u.repo.Logout() }
func (u *common) ValidateToken() response.ValidateToken { return u.repo.ValidateToken() }
func (u *common) GetUserInfo() response.Commons        { return u.repo.GetUserInfo() }

// Format formats the given value into table, json, or yaml and returns it as string.
func Format(format string, v interface{}) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		b, _ := json.MarshalIndent(v, "", "  ")
		return string(b) + "\n"
	case "yaml":
		b, _ := yaml.Marshal(v)
		return string(b)
	default:
		return tableString(v)
	}
}

func tableString(v interface{}) string {
	switch data := v.(type) {
	case response.RefreshToken:
		return refreshTableString(data)
	case *response.RefreshToken:
		return refreshTableString(*data)
	case response.Commons:
		return commonTableString(data)
	case *response.Commons:
		return commonTableString(*data)
	default:
		b, _ := json.Marshal(data)
		return string(b) + "\n"
	}
}

func newTabWriterBuf() (*tabwriter.Writer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 2, 4, 2, ' ', 0)
	return w, buf
}

func refreshTableString(res response.RefreshToken) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"FIELD", "VALUE"}, "\t"))
	fmt.Fprintf(w, "Code\t%s\n", res.Code)
	fmt.Fprintf(w, "Message\t%s\n", res.Message)
	if res.TokenPair != nil {
		fmt.Fprintf(w, "AccessToken\t%s\n", res.TokenPair.AccessToken)
		fmt.Fprintf(w, "RefreshToken\t%s\n", res.TokenPair.RefreshToken)
		fmt.Fprintf(w, "TokenType\t%s\n", res.TokenPair.TokenType)
		fmt.Fprintf(w, "ExpiresIn\t%d\n", res.TokenPair.ExpiresIn)
	}
	w.Flush()
	return buf.String()
}

func commonTableString(res response.Commons) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Commons) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "CREATED_AT"}, "\t"))
		for _, c := range res.Commons {
			created := ""
			if c.CreatedAt != nil {
				created = c.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", c.ID, c.UUID, created)
		}
	}
	w.Flush()
	return buf.String()
}

