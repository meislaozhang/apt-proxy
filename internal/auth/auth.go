package auth

import "crypto/subtle"

func EqualToken(a, b []byte) bool { return subtle.ConstantTimeCompare(a, b) == 1 }
