package main

import (
	"fmt"
)

func getFullNameFromOwnerAndRepo(owner string, repository string) string {
	return fmt.Sprintf("%s/%s", owner, repository)
}
