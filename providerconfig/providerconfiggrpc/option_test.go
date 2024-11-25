package providerconfiggrpc

import "testing"

func TestValidateEndpoint(t *testing.T) {
	err := validateEndpoint("localhost:8888")
	if err != nil {
		t.Errorf("Unexpected error while validating endpoint: %v", err)
	}
	err = validateEndpoint("http://127.0.0.1:8888")
	if err != nil {
		t.Errorf("Unexpected error while validating endpoint: %v", err)
	}
	err = validateEndpoint("https://127.0.0.1:8888")
	if err != nil {
		t.Errorf("Unexpected error while validating endpoint: %v", err)
	}

	err = validateEndpoint("httpx://127.0.0.1:8888")
	if err == nil {
		t.Error("Expected error while validating endpoint")
	}

	err = validateEndpoint("0:0:0:0:0:0:0:1:8888")
	if err != nil {
		t.Errorf("Unexpected error while validating endpoint: %v", err)
	}
	err = validateEndpoint("::1:8888")
	if err != nil {
		t.Errorf("Unexpected error while validating endpoint: %v", err)
	}
}
