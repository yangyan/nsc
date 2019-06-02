package nettest

import "testing"

func TestTestConnection(t *testing.T) {
	ok := TestConnection("127.0.0.1", 445)
	t.Log(ok)
}
