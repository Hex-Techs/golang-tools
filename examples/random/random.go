package main

import (
	"strings"

	"github.com/ymichaelson/golang-tools/random"
	"github.com/ymichaelson/klog"
)

func main() {
	// random 11 bytes string
	randomString := random.RandomId(11)
	klog.Infof("random 11 bytes string: %s", randomString)

	randomStringToLower := strings.ToLower(randomString)
	klog.Infof("random 11 bytes string lower: %s", randomStringToLower)

	randomStringToUpper := strings.ToUpper(randomString)
	klog.Infof("random 11 bytes string upper: %s", randomStringToUpper)

	// random int, it can be used to generate digital captcha
	randomIntBytes := random.RandomDigits(11)
	klog.Infof("randomIntBytes: %v", randomIntBytes)

	randomInt := random.RandomDigitsToInt(6)
	klog.Infof("randomInt: %d", randomInt)

	randomStrings := random.RandomDigitsToString(6)
	klog.Infof("randomString: %s", randomStrings)
}
