package uuid

import (
	"testing"
)

func TestVcloudUUID_ContainsPrefix(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "ContainsPrefix",
			uuid: VcloudUUID("urn:vcloud:vm:"),
			want: true,
		},
		{
			name: "DoesNotContainPrefix",
			uuid: VcloudUUID("urn:vm:"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.ContainsPrefix(); got != tt.want {
				t.Errorf("VcloudUUID.ContainsPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUUIDV4(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidUUID",
			args: args{
				uuid: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			},
			want: true,
		},
		{
			name: "InvalidUUID",
			args: args{
				uuid: "f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUUIDV4(tt.args.uuid); got != tt.want {
				t.Errorf("isUUIDV4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVcloudUUID_IsType(t *testing.T) {
	type args struct {
		prefix VcloudUUID
	}
	tests := []struct {
		name string
		uuid VcloudUUID
		args args
		want bool
	}{
		{
			name: "IsType",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			args: args{
				prefix: VM,
			},
			want: true,
		},
		{
			name: "IsNotType",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			args: args{
				prefix: User,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsType(tt.args.prefix); got != tt.want {
				t.Errorf("VcloudUUID.IsType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractUUIDv4(t *testing.T) {
	type args struct {
		uuid   string
		prefix VcloudUUID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ExtractUUID",
			args: args{
				uuid:   "urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479",
				prefix: VM,
			},
			want: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUUIDv4(tt.args.uuid, tt.args.prefix); got != tt.want {
				t.Errorf("extractUUIDv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidUUID",
			args: args{
				uuid: "urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479",
			},
			want: true,
		},
		{
			name: "InvalidUUID",
			args: args{
				uuid: "f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{
			name: "InvalidPrefix",
			args: args{
				uuid: "urn:vm:f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.uuid); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		prefix VcloudUUID
		uuid   string
	}
	tests := []struct {
		name string
		args args
		want VcloudUUID
	}{
		{
			name: "Normalize",
			args: args{
				prefix: VM,
				uuid:   "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			},
			want: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.prefix, tt.args.uuid); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVcloudUUID_IsVM(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsVM",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotVM",
			uuid: VcloudUUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVM(); got != tt.want {
				t.Errorf("VcloudUUID.IsVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVcloudUUID_IsUser(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsUser",
			uuid: VcloudUUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotUser",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsUser(); got != tt.want {
				t.Errorf("VcloudUUID.IsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGroup.
func TestVcloudUUID_IsGroup(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsGroup",
			uuid: VcloudUUID("urn:vcloud:group:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotGroup",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsGroup(); got != tt.want {
				t.Errorf("VcloudUUID.IsGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGateway.
func TestVcloudUUID_IsGateway(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsGateway",
			uuid: VcloudUUID("urn:vcloud:gateway:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotGateway",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsGateway(); got != tt.want {
				t.Errorf("VcloudUUID.IsGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDC.
func TestVcloudUUID_IsVDC(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsVDC",
			uuid: VcloudUUID("urn:vcloud:vdc:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotVDC",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVDC(); got != tt.want {
				t.Errorf("VcloudUUID.IsVDC() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsNetwork.
func TestVcloudUUID_IsNetwork(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsNetwork",
			uuid: VcloudUUID("urn:vcloud:network:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotNetwork",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsNetwork(); got != tt.want {
				t.Errorf("VcloudUUID.IsNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsLoadBalancerPool.
func TestVcloudUUID_IsLoadBalancerPool(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsLoadBalancerPool",
			uuid: VcloudUUID("urn:vcloud:loadBalancerPool:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotLoadBalancerPool",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsLoadBalancerPool(); got != tt.want {
				t.Errorf("VcloudUUID.IsLoadBalancerPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCStorageProfile.
func TestVcloudUUID_IsVDCStorageProfile(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsVDCStorageProfile",
			uuid: VcloudUUID("urn:vcloud:vdcstorageProfile:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotVDCStorageProfile",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVDCStorageProfile(); got != tt.want {
				t.Errorf("VcloudUUID.IsVDCStorageProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPP.
func TestVcloudUUID_IsVAPP(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsVAPP",
			uuid: VcloudUUID("urn:vcloud:vapp:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: true,
		},
		{
			name: "IsNotVAPP",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVAPP(); got != tt.want {
				t.Errorf("VcloudUUID.IsVAPP() = %v, want %v", got, tt.want)
			}
		})
	}
}
