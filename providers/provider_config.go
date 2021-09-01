package providers

import "os"

type ProviderConfig struct {
	keydir     string
	btcAddr    string
	accountNum int
	pwdFile    *os.File
}
