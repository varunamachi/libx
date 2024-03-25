package auth

import "context"

type Role string

const (
	None   Role = ""
	Normal Role = "Normal"
	Admin  Role = "Admin"
	Super  Role = "Super"
)

func (r Role) IsOneOf(others ...Role) bool {
	for _, oth := range others {
		if r == oth {
			return true
		}
	}
	return false
}

func (r Role) EqualOrAbove(another Role) bool {
	switch r {
	case None:
		return another == None
	case Normal:
		return another.IsOneOf(None, Normal)
	case Admin:
		return another.IsOneOf(None, Normal, Admin)
	case Super:
		return true
	}
	return false
}

// type User struct {
// 	SeqId       int           `json:"seqId" db:"seqId" bson:"seqId"`
// 	Id          string        `json:"id" db:"id" bson:"id"`
// 	EMail       string        `json:"email" db:"email" bson:"email"`
// 	FirstName   string        `json:"firstName" db:"firstName" bson:"firstName"`
// 	LastName    string        `json:"lastName" db:"lastName" bson:"lastName"`
// 	Role        Role          `json:"role" db:"role" bson:"role"`
// 	GroupsIDs   []string      `json:"groups" bson:"groups"`
// 	Permissions PermissionSet `json:"permissions" bson:"permissions"`
// }

func HasRole(u User, role Role) bool {
	// return role.EqualOrAbove(u.Role())
	return u.Role().EqualOrAbove(role)
}

func HasPerms(u User, permIds ...string) bool {
	for _, perm := range permIds {
		if !u.Permissions().HasPerm(perm) {
			return false
		}
	}
	return true
}

type User interface {
	SeqId() int64
	Id() string
	Email() string
	FullName() string
	Role() Role
	GroupIds() []string
	Permissions() PermissionSet
}

type Group struct {
	SeqId int    `json:"seqId" db:"seq_id" bson:"seqId"`
	Id    string `json:"id" db:"id" bson:"id"`
	Name  string `json:"name" db:"name" bson:"name"`
}

func ToRole(roleStr string) Role {
	switch roleStr {
	case "None":
		return None
	case "Normal":
		return Normal
	case "Admin":
		return Admin
	case "Super":
		return Super
	}
	return None
}

// TODO - make this and interface, so that implementation has more freedom
// to access the logic that gets the user information
// type UserRetrieverFunc func(
// 	gtx context.Context, userId string) (User, error)

type UserRetriever interface {
	GetUser(gtx context.Context, userId string) (User, error)
}
