package spikeutils_test

import (
	"testing"

	aero "github.com/aerospike/aerospike-client-go"
	"github.com/brunoyin/spike-cli/spikeutils"
)

// TestGetClient test. Start a local Docker Aerospike container before testing.
func TestGetClient(t *testing.T) {
	clientPolicy := aero.NewClientPolicy()
	client, err := spikeutils.GetClient(clientPolicy, "127.0.0.1", 3000)
	if err != nil {
		t.Errorf("Expect no error. But got %v", err)
	} else {
		defer client.Close()
		name := client.GetNodeNames()[0]
		if name == "" {
			t.Errorf("Node name should not be blank!")
		}
	}
}
