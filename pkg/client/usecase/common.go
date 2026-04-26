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
type Common interface {
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

func (u *common) ValidateToken() response.ValidateToken { return u.repo.ValidateToken() }
func (u *common) GetUserInfo() response.Commons         { return u.repo.GetUserInfo() }

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
	case response.Commons:
		return commonTableString(data)
	case *response.Commons:
		return commonTableString(*data)
	case response.IdPUsers:
		return idpUsersTableString(data)
	case response.SingleIdPUser:
		return singleIdPUserTableString(data)
	case response.IdPGroups:
		return idpGroupsTableString(data)
	case response.SingleIdPGroup:
		return singleIdPGroupTableString(data)
	case response.Resources:
		return resourcesTableString(data)
	case response.SingleResource:
		return singleResourceTableString(data)
	case response.ResourceGroupRoles:
		return resourceGroupRolesTableString(data)
	default:
		b, _ := json.Marshal(data)
		return string(b) + "\n"
	}
}

func idpUsersTableString(res response.IdPUsers) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Users) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "USERNAME", "EMAIL", "FIRST_NAME", "LAST_NAME", "ENABLED"}, "\t"))
		for _, u := range res.Users {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%v\n", u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.Enabled)
		}
	}
	w.Flush()
	return buf.String()
}

func singleIdPUserTableString(res response.SingleIdPUser) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if res.User != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "USERNAME", "EMAIL", "FIRST_NAME", "LAST_NAME", "ENABLED"}, "\t"))
		u := res.User
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%v\n", u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.Enabled)
	}
	w.Flush()
	return buf.String()
}

func idpGroupsTableString(res response.IdPGroups) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Groups) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "NAME", "PATH"}, "\t"))
		for _, g := range res.Groups {
			fmt.Fprintf(w, "%s\t%s\t%s\n", g.ID, g.Name, g.Path)
		}
	}
	w.Flush()
	return buf.String()
}

func singleIdPGroupTableString(res response.SingleIdPGroup) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if res.Group != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "NAME", "PATH"}, "\t"))
		g := res.Group
		fmt.Fprintf(w, "%s\t%s\t%s\n", g.ID, g.Name, g.Path)
	}
	w.Flush()
	return buf.String()
}

func resourcesTableString(res response.Resources) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Resources) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"UUID", "NAME", "DESCRIPTION", "CREATED_BY"}, "\t"))
		for _, r := range res.Resources {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.UUID, r.Name, r.Description, r.CreatedBy)
		}
	}
	w.Flush()
	return buf.String()
}

func singleResourceTableString(res response.SingleResource) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if res.Resource != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"UUID", "NAME", "DESCRIPTION", "CREATED_BY"}, "\t"))
		r := res.Resource
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.UUID, r.Name, r.Description, r.CreatedBy)
	}
	w.Flush()
	return buf.String()
}

func resourceGroupRolesTableString(res response.ResourceGroupRoles) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Groups) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"RESOURCE_UUID", "GROUP_UUID", "ROLE"}, "\t"))
		for _, g := range res.Groups {
			fmt.Fprintf(w, "%s\t%s\t%s\n", g.ResourceUUID, g.GroupUUID, g.Role)
		}
	}
	w.Flush()
	return buf.String()
}

func newTabWriterBuf() (*tabwriter.Writer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 2, 4, 2, ' ', 0)
	return w, buf
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

