package auth

type PermissionTree struct {
	Permissions []*PermissionNode `json:"permissions"`
}

type PermissionNode struct {
	IdNum      int               `json:"idNum" db:"idNum"`
	PermId     string            `json:"permId" db:"permId"`
	Name       string            `json:"name" db:"name"`
	Predefined bool              `json:"predefined" db:"predefined"`
	Base       string            `json:"base" db:"base"`
	Children   []*PermissionNode `json:"children"`
}

func (pn *PermissionNode) AddChild(child *PermissionNode) {
	pn.Children = append(pn.Children, child)
}

type PermissionSet map[string]struct{}

func (pm PermissionSet) HasPerm(permId string) bool {
	if permId == "" {
		return true
	}
	_, found := pm[permId]
	return found
}

func MergePerms(permSets []PermissionSet) PermissionSet {
	result := PermissionSet{}

	for _, set := range permSets {
		for perm := range set {
			if _, found := result[perm]; !found {
				result[perm] = struct{}{}
			}
		}
	}
	return result
}
