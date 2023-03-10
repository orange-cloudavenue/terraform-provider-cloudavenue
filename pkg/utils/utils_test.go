package utils

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_generateUUID(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "generateUUID",
			args: args{
				str: "test",
			},
			want: "e8b764da-5fe5-51ed-8af8-c5c6eca28d7a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateUUID(tt.args.str); got != tt.want {
				t.Errorf("generateUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUUID(t *testing.T) {
	type args[T tfValuesForUUID] struct {
		value T
	}

	type testTF[T tfValuesForUUID] struct {
		name string
		args args[T]
		want types.String
	}

	testString := []testTF[string]{
		{
			name: "GenerateUUIDFromString",
			args: args[string]{
				value: "test",
			},
			want: types.StringValue("e8b764da-5fe5-51ed-8af8-c5c6eca28d7a"),
		},
		{
			name: "GenerateUUIDFromSliceString",
			args: args[string]{
				value: "test2",
			},
			want: types.StringValue("a4065824-05e5-5a82-9841-cd5efc76b8c1"),
		},
	}

	testSliceString := []testTF[[]string]{
		{
			name: "GenerateUUIDFromSliceString",
			args: args[[]string]{
				value: []string{"test", "test2"},
			},
			want: types.StringValue("016fab6f-3c2d-5b38-b6fc-421aff431b61"),
		},
	}

	for _, tt := range testString {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateUUID(tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateUUID() from testString = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range testSliceString {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateUUID(tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateUUID() from testSliceString = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeBoolPointer(t *testing.T) {
	type args struct {
		value bool
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "TakeBoolPointer",
			args: args{
				value: true,
			},
			want: &[]bool{true}[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TakeBoolPointer(tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeBoolPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeIntPointer(t *testing.T) {
	type args struct {
		value int
	}
	tests := []struct {
		name string
		args args
		want *int
	}{
		{
			name: "TakeIntPointer",
			args: args{
				value: 666,
			},
			want: &[]int{666}[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TakeIntPointer(tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeIntPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeInt64Pointer(t *testing.T) {
	type args struct {
		x int64
	}
	tests := []struct {
		name string
		args args
		want *int64
	}{
		{
			name: "TakeInt64Pointer",
			args: args{
				x: 666,
			},
			want: &[]int64{666}[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TakeInt64Pointer(tt.args.x)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeInt64Pointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
