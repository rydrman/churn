package churn

import "testing"

func TestChannelBufferSize(t *testing.T) {

	g := NewGraph(ChannelBufferSize(10))
	if g.channelBufferSize != 10 {
		t.Error("expected channel buffer size to be set")
	}

}
