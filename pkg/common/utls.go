package common

import (
	utls "github.com/refraction-networking/utls"
)

// List of actually supported ciphers(not a list of offered ciphers!)
// Essentially all working AES_GCM_128 ciphers
var tapDanceSupportedCiphers = []uint16{
	utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
}

func ForceSupportedCiphersFirst(suites []uint16) []uint16 {
	swapSuites := func(i, j int) {
		if i == j {
			return
		}
		tmp := suites[j]
		suites[j] = suites[i]
		suites[i] = tmp
	}
	lastSupportedCipherIdx := 0
	for i := range suites {
		for _, supportedS := range tapDanceSupportedCiphers {
			if suites[i] == supportedS {
				swapSuites(i, lastSupportedCipherIdx)
				lastSupportedCipherIdx += 1
			}
		}
	}
	alwaysSuggestedSuite := utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	for i := range suites {
		if suites[i] == alwaysSuggestedSuite {
			return suites
		}
	}
	return append([]uint16{alwaysSuggestedSuite}, suites[lastSupportedCipherIdx:]...)
}
