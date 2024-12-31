package main

import (
	goxAuthSimpleS2S "github.com/devlibx/gox-base/v2/auth/s2s/simple"
	"os"
	"time"
)

func main() {
	user := os.Getenv("USER_TO_GIVE_ACCESS")
	salt := os.Getenv("USER_TO_GIVE_ACCESS_SALT")
	token := goxAuthSimpleS2S.GenerateAdhocUserSecret(time.Duration(180)*time.Minute, user, salt)
	println(token)
}
