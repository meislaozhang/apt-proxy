package protocol

import "testing"

func TestCapabilitiesRoundTrip(t *testing.T) {
    c := Capabilities{}
    c.Enable(CapabilityFlowControl)
    c.Enable(CapabilityDatagram)
    c.Enable(CapabilityConnectUDP)
    got, err := DecodeCapabilities(c.Encode())
    if err != nil { t.Fatal(err) }
    if got.Bits != c.Bits { t.Fatalf("bits=%x want=%x", got.Bits, c.Bits) }
    if !got.Has(CapabilityDatagram) || got.Has(CapabilityConnectIP) { t.Fatal("unexpected capability set") }
}

func TestCapabilitiesRejectBadLength(t *testing.T) {
    if _, err := DecodeCapabilities([]byte{1, 2, 3}); err == nil { t.Fatal("expected malformed capabilities payload") }
}
