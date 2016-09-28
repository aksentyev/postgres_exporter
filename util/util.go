package util

import (
    "github.com/aksentyev/hubble/hubble"
)

func IncludesStr(source []string, target string) bool {
    for _, el := range source {
        if el == target {
            return true
        }
    }
    return false
}

func PgConnURL(s *hubble.ServiceAtomic) string {
	url := "postgres://" + s.Address + ":" + s.Port + "/"
	url += s.ExporterOptions["db"] + "?"
	url += "user=" + s.ExporterOptions["user"] + "&"
	url += "password=" + s.ExporterOptions["password"] + "&"
	url += "sslmode=disable" + "&"
	return url
}
