package utils

const (
	RoleGuest    = "guest"
	RoleClient   = "client"
	RoleArtisant = "artisant"
	RoleAdmin    = "admin"
)

func GetValidRoles(role string) int {
	switch role {
	case RoleGuest:
		return 1
	case RoleClient:
		return 2
	case RoleArtisant:
		return 3
	case RoleAdmin:
		return 4
	default:
		return 0
	}
}
