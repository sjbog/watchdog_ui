package models

import (
	"fmt"
	"github.com/revel/revel"
	"regexp"
)

type User struct {
	Id		uint
	Name	string
	Username, Password	string
	HashedPassword		[]byte
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

var users	= map [ string ] * User {
	"root"	: & User {
			Id			: 0,
			Name		: "Root",
			Username	: "root",
			HashedPassword	: []byte ( "$2a$10$ypJobnj6NvOnEzfMGdAC8eKM/Q0bAkyfk6b0zxt1KqWzruL9KMxbW" ),
		},
}

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{31},
		revel.MinSize{4},
		revel.Match{userRegex},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

	v.Check(user.Name,
		revel.Required{},
		revel.MaxSize{31},
		revel.MinSize{4},
	)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{63},
		revel.MinSize{8},
	)
}

func	GetUser ( username   string )	( user  * User )	{
	user	= users [ username ]
	return
}


func	( self  * User )	RememberAuth ( request  * revel.Request )	{
}

