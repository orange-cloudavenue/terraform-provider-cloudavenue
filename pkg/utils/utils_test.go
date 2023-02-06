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
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    types.String
		wantErr bool
	}{
		{
			name: "GenerateUUIDFromString",
			args: args{
				value: "test",
			},
			want:    types.StringValue("e8b764da-5fe5-51ed-8af8-c5c6eca28d7a"),
			wantErr: false,
		},
		{
			name: "GenerateUUIDFromSliceString",
			args: args{
				value: []string{"test", "test2"},
			},
			want:    types.StringValue("56a2aec8-8045-5727-91fc-7194fd4e339f"),
			wantErr: false,
		},
		{
			name: "GenerateUUIDFromIntError",
			args: args{
				value: 1,
			},
			want:    types.StringNull(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateUUID(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
