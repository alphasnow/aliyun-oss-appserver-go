package appserver

import (
	"github.com/jarcoal/httpmock"
	"testing"
)

func TestCallbackGetPublicKey(t *testing.T) {
	pk := "-----BEGIN PUBLIC KEY-----\nMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAKs/JBGzwUB2aVht4crBx3oIPBLNsjGs\nC0fTXv+nvlmklvkcolvpvXLTjaxUHR3W9LXxQ2EHXAJfCB+6H2YF1k8CAwEAAQ==\n-----END PUBLIC KEY-----"

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	httpmock.RegisterResponder("GET", "https://gosspublic.alicdn.com/callback_pub_key_v1.pem",
		httpmock.NewStringResponder(200, pk))

	resp, err := GetPublicKey("aHR0cHM6Ly9nb3NzcHVibGljLmFsaWNkbi5jb20vY2FsbGJhY2tfcHViX2tleV92MS5wZW0=")
	if err != nil {
		t.Error(err)
	}
	if string(resp) != pk {
		t.Errorf("expect %s, got %s", pk, string(resp))
	}
}
