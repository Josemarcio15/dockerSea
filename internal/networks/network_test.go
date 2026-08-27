package networks

import "testing"

func TestNormalizeNetworkCreateRequest(t *testing.T) {
	req := NormalizeNetworkCreateRequest(NetworkCreateRequest{
		Name:   "  frontend  ",
		Driver: "bridge",
	})

	if req.Name != "frontend" || req.Driver != "bridge" {
		t.Fatalf("request was not normalized: %#v", req)
	}
}

func TestValidateNetworkCreateRequest(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkCreateRequest
		want string
	}{
		{name: "missing name", req: NetworkCreateRequest{}, want: "Nome da rede não informado"},
		{name: "valid", req: NetworkCreateRequest{Name: "frontend"}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateNetworkCreateRequest(tt.req); got != tt.want {
				t.Fatalf("ValidateNetworkCreateRequest() = %q, want %q", got, tt.want)
			}
		})
	}
}
