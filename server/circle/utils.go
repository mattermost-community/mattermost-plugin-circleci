package circle

import (
	"fmt"

	"github.com/TomTucka/go-circleci/circleci"
)

func GetFullNameFromOwnerAndRepo(owner string, repository string) string {
	return fmt.Sprintf("%s/%s", owner, repository)
}

// Return in the format "FullName (username)"
func CircleciUserToString(user *circleci.User) string {
	if user.Name != "" {
		return user.Name + " (" + user.Login + ")"
	}

	return user.Login
}
