package resolver

import (
	structs "github.com/callistom/go-graphql-auth/structs"
	graphql "github.com/neelance/graphql-go"
)

// UserResolver struct
type UserResolver struct {
	User structs.User
}

type UsersResolver struct {
	User structs.User
}

func (v *UsersResolver) ID() graphql.ID {
	return graphql.ID(v.User.ID)
}

// Name return name
func (v *UsersResolver) Name() string {
	return v.User.Name
}

// Mail return mail
func (v *UsersResolver) Mail() string {
	return v.User.Mail
}

// ID return ID
func (v *UserResolver) ID() graphql.ID {
	return graphql.ID(v.User.ID)
}

// Name return name
func (v *UserResolver) Name() string {
	return v.User.Name
}

// Mail return mail
func (v *UserResolver) Mail() string {
	return v.User.Mail
}
