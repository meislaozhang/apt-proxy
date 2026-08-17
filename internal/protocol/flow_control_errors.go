package protocol

import "errors"

// ErrFlowWindowOverflow indicates that a WINDOW_UPDATE would exceed the
// configured maximum receive/send credit. It is a sentinel so callers and
// tests can distinguish protocol overflow from other errors.
var ErrFlowWindowOverflow = errors.New("flow-window update exceeds maximum")
