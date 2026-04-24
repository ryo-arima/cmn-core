package controller_test

import (
	"testing"

	"github.com/ryo-arima/cmn-core/pkg/server/controller"
	mock "github.com/ryo-arima/cmn-core/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewCommonControllerForPublic(t *testing.T) {
	// CommonPublic controller does not exist as a standalone; skip.
	t.Skip("CommonPublic controller is not a standalone type")
}

func TestNewCommonControllerForInternal(t *testing.T) {
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewCommonInternal(commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForPrivate(t *testing.T) {
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewCommonPrivate(commonRepo)

	assert.NotNil(t, ctrl)
}
