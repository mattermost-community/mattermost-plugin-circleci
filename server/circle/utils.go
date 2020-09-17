package circle

import "fmt"

func GetFullNameFromOwnerAndRepo(owner string, repository string) string {
	return fmt.Sprintf("%s/%s", owner, repository)
}
