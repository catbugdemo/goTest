package init

import "testing"

func TestConfig(t *testing.T) {
	t.Run("GetConfig", func(t *testing.T) {
		getConfig := GetConfig()
		getConfig.Get("")
	})
}
