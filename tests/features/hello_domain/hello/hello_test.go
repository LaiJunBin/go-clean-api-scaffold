package hello_test

import (
	"testing"

	"go-clean-api-scaffold/tests/testutils"
)

func TestHello(t *testing.T) {
	testutils.TestFeature(t, "say_hello", testutils.NewDefaultScenarioInitializer())
}
