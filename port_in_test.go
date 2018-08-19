package churn

import "testing"

func TestInPortCore_AddReceiver(t *testing.T) {

	core := new(InPortCore)
	err := core.AddReceiver(nil)
	if err == nil {
		t.Error("should not allow receivers on in portss")
	}

}
