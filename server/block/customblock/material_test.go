package customblock

import "testing"

// The client reads ambient_occlusion as an intensity. Encoding it as a bool is silently
// rejected, so the wire type matters as much as the value.
func TestMaterial_EncodeAmbientOcclusionIsAnIntensity(t *testing.T) {
	for _, tc := range []struct {
		name     string
		material Material
		want     float32
	}{
		{"opaque render method defaults to vanilla strength", NewMaterial("t", OpaqueRenderMethod()), 1},
		{"alpha test render method defaults to off", NewMaterial("t", AlphaTestRenderMethod()), 0},
		{"enabled explicitly", NewMaterial("t", AlphaTestRenderMethod()).WithAmbientOcclusion(), 1},
		{"disabled explicitly", NewMaterial("t", OpaqueRenderMethod()).WithoutAmbientOcclusion(), 0},
		{"custom intensity", NewMaterial("t", OpaqueRenderMethod()).WithAmbientOcclusionIntensity(2.5), 2.5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.material.Encode()["ambient_occlusion"]
			got, ok := encoded.(float32)
			if !ok {
				t.Fatalf("ambient_occlusion encoded as %T, want float32", encoded)
			}
			if got != tc.want {
				t.Fatalf("ambient_occlusion = %v, want %v", got, tc.want)
			}
		})
	}
}
