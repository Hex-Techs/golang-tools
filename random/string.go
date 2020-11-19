/*
Copyright The Michaelson.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package random

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
)

// idLen is a length of captcha id string.
// (20 bytes of 62-letter alphabet give ~119 bits.)
const idLen = 20

// idChars are characters allowed in captcha id.
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// rngKey is a secret key used to deterministically derive seeds for
// PRNGs used in image and audio. Generated once during initialization.
var rngKey [32]byte

func init() {
	if _, err := io.ReadFull(rand.Reader, rngKey[:]); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
}

// RandomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-9. The slice can be used as a captcha
// solution.
func RandomDigits(length int) []byte {
	return randomBytesMod(length, 10)
}

func RandomDigitsToString(length int) string {
	return randomdigitsToString(randomBytesMod(length, 10))
}

func RandomDigitsToInt(length int) int {
	return randomdigitsToInt(randomBytesMod(length, 10))
}

func randomdigitsToString(b []byte) string {
	var digitsString string
	for _, tmp := range b {
		digitsString = fmt.Sprintf("%s%d", digitsString, int(tmp))
	}

	return digitsString
}

func randomdigitsToInt(b []byte) int {
	digitsInt, _ := strconv.Atoi(randomdigitsToString(b))
	return digitsInt
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}

}

// randomId returns a new random id string.
func RandomId(l int) string {
	var length int
	if l > 0 {
		length = l
	} else {
		length = idLen
	}
	b := randomBytesMod(length, byte(len(idChars)))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b)
}
