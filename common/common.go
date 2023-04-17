package common

var binPath string

func GetBinPath() string {
	return binPath
}

func SetBinPath(path string) {
	binPath = path
}
