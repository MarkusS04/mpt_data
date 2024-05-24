package vars

import (
	"mpt_data/models/apimodel"
	"os"
)

func GetDbPAth() string {
	return os.ExpandEnv("$HOME/apps/mpt_testing")
}

var (
	UserAPI = apimodel.UserLogin{Username: "Max", Password: "admin"}
)
