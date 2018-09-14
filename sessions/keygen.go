package sessions

import "ztaylor.me/keygen"

// Keygen is replaceable session key generator
var Keygen = func() string {
	return keygen.NewVal()
}
