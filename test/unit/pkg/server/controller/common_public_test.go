package controller_test

import (
	"testing"
)

func TestNewCommonControllerForPublic(t *testing.T) {
	// CommonPublic controller does not exist as a standalone; skip.
	t.Skip("CommonPublic controller is not a standalone type")
}

func TestNewCommonControllerForInternal(t *testing.T) {
	t.Skip("CommonInternal merged into CommonShare; constructor removed")
}

func TestNewCommonControllerForPrivate(t *testing.T) {
	t.Skip("CommonPrivate merged into CommonShare; constructor removed")
}
