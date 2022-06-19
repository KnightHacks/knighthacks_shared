package models

import (
	"fmt"
	"io"
	"strconv"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	// for now keep this the same
	RoleSponsor Role = "SPONSOR"
	RoleNormal  Role = "NORMAL"
	RoleOwns    Role = "OWNS"
)

var AllRole = []Role{
	RoleAdmin,
	RoleSponsor,
	RoleNormal,
	RoleOwns,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleSponsor, RoleNormal, RoleOwns:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Provider string

const (
	ProviderGithub Provider = "GITHUB"
	ProviderGmail  Provider = "GMAIL"
)

var AllProvider = []Provider{
	ProviderGithub,
	ProviderGmail,
}

func (e Provider) IsValid() bool {
	switch e {
	case ProviderGithub, ProviderGmail:
		return true
	}
	return false
}

func (e Provider) String() string {
	return string(e)
}

func (e *Provider) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Provider(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Provider", str)
	}
	return nil
}

func (e Provider) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
