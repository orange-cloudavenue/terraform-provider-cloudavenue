package uuid

import (
	"testing"
)

const (
	validUUIDv4 = "12345678-1234-1234-1234-123456789012"
)

func TestVcloudUUID_ContainsPrefix(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "ContainsPrefix",
			uuid: VcloudUUID(VM.String() + validUUIDv4),
			want: true,
		},
		{
			name: "DoesNotContainPrefix",
			uuid: VcloudUUID("urn:vm:" + validUUIDv4),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
				uuid: validUUIDv4,
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
		{ // Empty string
			name: "EmptyString",
			args: args{
				uuid: "",
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
			uuid: VcloudUUID(VM.String() + validUUIDv4),
			args: args{
				prefix: VM,
			},
			want: true,
		},
		{
			name: "IsNotType",
			uuid: VcloudUUID(VM.String() + validUUIDv4),
			args: args{
				prefix: User,
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
			args: args{
				prefix: VM,
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
				uuid:   "urn:vcloud:vm:" + validUUIDv4,
				prefix: VM,
			},
			want: validUUIDv4,
		},
		{
			name: "EmptyString",
			args: args{
				uuid:   "",
				prefix: VM,
			},
			want: "",
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
				uuid: "urn:vcloud:vm:" + validUUIDv4,
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
		{ // Empty string
			name: "EmptyString",
			args: args{
				uuid: "",
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
				uuid:   validUUIDv4,
			},
			want: VcloudUUID("urn:vcloud:vm:" + validUUIDv4),
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
			uuid: VcloudUUID(VM.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVM",
			uuid: VcloudUUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(User.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotUser",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(Group.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGroup",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(Gateway.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGateway",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(VDC.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDC",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(Network.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotNetwork",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(LoadBalancerPool.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotLoadBalancerPool",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(VDCStorageProfile.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDCStorageProfile",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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
			uuid: VcloudUUID(VAPP.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVAPP",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
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

// IsDisk.
func TestVcloudUUID_IsDisk(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsDisk",
			uuid: VcloudUUID(Disk.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotDisk",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsDisk(); got != tt.want {
				t.Errorf("VcloudUUID.IsDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsSecurityGroup.
func TestVcloudUUID_IsSecurityGroup(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsSecurityGroup",
			uuid: VcloudUUID(SecurityGroup.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotSecurityGroup",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsSecurityGroup(); got != tt.want {
				t.Errorf("VcloudUUID.IsSecurityGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsAppPortProfile.
func TestVcloudUUID_IsAppPortProfile(t *testing.T) {
	tests := []struct {
		name string
		uuid VcloudUUID
		want bool
	}{
		{
			name: "IsAppPortProfile",
			uuid: VcloudUUID(AppPortProfile.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotAppPortProfile",
			uuid: VcloudUUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: VcloudUUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsAppPortProfile(); got != tt.want {
				t.Errorf("VcloudUUID.IsAppPortProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsType tests the TestIsType function.
func TestTestIsType(t *testing.T) {
	testCases := []struct {
		name     string
		uuidType VcloudUUID
		uuid     VcloudUUID
		want     bool
	}{
		{
			name:     "valid uuid",
			uuidType: VM,
			uuid:     VcloudUUID(VM.String() + validUUIDv4),
			want:     true,
		},
		{
			name:     "invalid uuid",
			uuidType: VM,
			uuid:     "invalid-uuid",
			want:     false,
		},
		{
			name:     "empty value",
			uuidType: VM,
			uuid:     "",
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := TestIsType(tc.uuidType)(tc.uuid.String())
			if tc.want && err != nil {
				t.Errorf("TestIsType() = %v, want %v", err, tc.want)
			}
		})
	}
}
