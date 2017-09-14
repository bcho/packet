package packet

import "os"

func inTravis() bool {
	return os.Getenv("TRAVIS_GO_VERSION") != ""
}
