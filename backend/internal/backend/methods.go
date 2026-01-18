package backend

import(
	"strconv"
)
func StringToInt(param string) (int64,error){
	intVersion, err := strconv.ParseInt(param, 10, 64)
	return intVersion,err
}