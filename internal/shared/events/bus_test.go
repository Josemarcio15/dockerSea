package events

import "testing"

func TestBusPublish(t *testing.T) {
	bus := NewBus()
	called := false
	unsubscribe := bus.Subscribe("test", func(payload any) { called = payload == "ok" })
	bus.Publish("test", "ok")
	if !called {
		t.Fatal("subscriber was not called")
	}
	called = false
	unsubscribe()
	bus.Publish("test", "ok")
	if called {
		t.Fatal("unsubscribed handler was called")
	}
}
