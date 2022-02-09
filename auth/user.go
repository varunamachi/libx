package auth

type Role string

const (
	None   Role = "None"
	Normal Role = "Normal"
	Admin  Role = "Admin"
	Super  Role = "Super"
)

func (r Role) EqualOrAbove(another Role) bool {
	// Following logic only checks above condition
	switch r {
	case None:
		return true
	case Normal:
		return another == None
	case Admin:
		return another == None || another == Normal
	case Super:
		return false
	}

	// Following checks equal
	return r == another
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
	return role.EqualOrAbove(u.Role())
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
	SeqId() int
	Id() string
	Email() string
	FullName() string
	Role() Role
	GroupIds() []string
	Permissions() PermissionSet
}

type Group struct {
	SeqId int    `json:"seqId" db:"seqId" bson:"seqId"`
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

type UserRetrieverFunc func(userId string) (User, error)
