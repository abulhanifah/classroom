package helpers

import "fmt"

func CheckAcl(ctx Context, key string) bool {
	fmt.Println(ctx.Path())
	return true
}

func GetAclKeyFromRequest(ctx Context) string {
	return ""
}
