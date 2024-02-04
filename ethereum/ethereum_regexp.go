package ethereum

import (
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

var regexpAddress = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func IsValidAddress(addr string) (string, common.Address, bool) {
	const msg = "invalid ethereum address"

	if !regexpAddress.MatchString(addr) {
		return msg, common.Address{}, false
	}

	return "", common.HexToAddress(addr), true
}
