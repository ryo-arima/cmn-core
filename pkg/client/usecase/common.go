package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"

	clientauth "github.com/ryo-arima/cmn-core/pkg/client/share"
	"github.com/ryo-arima/cmn-core/pkg/client/repository"
	"github.com/ryo-arima/cmn-core/pkg/config"
	"github.com/ryo-arima/cmn-core/pkg/entity/response"
)

// Common is the business-logic interface for token-related operations.
type Common interface {
	ValidateToken() response.RrValidateToken
	GetUserInfo() response.RrCommons
}

type common struct {
	repo repository.Common
}

// NewCommon creates a Common usecase.
// manager is used to obtain and inject auth tokens automatically.
func NewCommon(conf config.BaseConfig, manager *clientauth.Manager) Common {
	return &common{repo: repository.NewCommon(conf, manager)}
}

func (u *common) ValidateToken() response.RrValidateToken { return u.repo.ValidateToken() }
func (u *common) GetUserInfo() response.RrCommons         { return u.repo.GetUserInfo() }

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
	case response.RrCommons:
		return commonTableString(data)
	case *response.RrCommons:
		return commonTableString(*data)
	case response.RrIdPUsers:
		return idpUsersTableString(data)
	case response.RrSingleIdPUser:
		return singleIdPUserTableString(data)
	case response.RrIdPGroups:
		return idpGroupsTableString(data)
	case response.RrSingleIdPGroup:
		return singleIdPGroupTableString(data)
	case response.RrResources:
		return resourcesTableString(data)
	case response.RrSingleResource:
		return singleResourceTableString(data)
	case response.RrResourceGroupRoles:
		return resourceGroupRolesTableString(data)
	default:
		b, _ := json.Marshal(data)
		return string(b) + "\n"
	}
}

func idpUsersTableString(res response.RrIdPUsers) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Users) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "USERNAME", "EMAIL", "FIRST_NAME", "LAST_NAME", "ENABLED", "ROLE"}, "\t"))
		for _, u := range res.Users {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%v\t%s\n", u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.Enabled, u.Role)
		}
	}
	w.Flush()
	return buf.String()
}

func singleIdPUserTableString(res response.RrSingleIdPUser) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if res.User != nil {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "USERNAME", "EMAIL", "FIRST_NAME", "LAST_NAME", "ENABLED", "ROLE"}, "\t"))
		u := res.User
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%v\t%s\n", u.ID, u.Username, u.Email, u.FirstName, u.LastName, u.Enabled, u.Role)
	}
	w.Flush()
	return buf.String()
}

func idpGroupsTableString(res response.RrIdPGroups) string {
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

func singleIdPGroupTableString(res response.RrSingleIdPGroup) string {
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

func resourcesTableString(res response.RrResources) string {
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

func singleResourceTableString(res response.RrSingleResource) string {
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

func resourceGroupRolesTableString(res response.RrResourceGroupRoles) string {
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

func commonTableString(res response.RrCommons) string {
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

