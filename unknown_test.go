package unknown

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	type args struct {
		obj         interface{}
		otherFields map[string]json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		wantBts []byte
		wantErr bool
	}{
		{
			"marshal same field",
			args{
				struct {
					ID int `json:"id"`
				}{666},
				map[string]json.RawMessage{"id": []byte(`777`)},
			},
			[]byte(`{"id":666}`),
			false,
		},
		{
			"marshal other field",
			args{
				struct {
					ID int `json:"id"`
				}{666},
				map[string]json.RawMessage{"foo": []byte(`777`)},
			},
			[]byte(`{"foo":777,"id":666}`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBts, err := MarshalJSON(tt.args.obj, tt.args.otherFields)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBts, tt.wantBts) {
				t.Errorf("MarshalJSON() = %s, want %s", gotBts, tt.wantBts)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type Obj struct {
		ID int
	}
	type args struct {
		jsonStr []byte
		obj     Obj
	}
	tests := []struct {
		name            string
		args            args
		wantOtherFields map[string]json.RawMessage
		wantObj         Obj
		wantErr         bool
	}{
		{
			"marshal other field",
			args{
				[]byte(`{"id":666,"name":"god"}`),
				Obj{},
			},
			map[string]json.RawMessage{"name": []byte(`"god"`)},
			Obj{666},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOtherFields, err := UnmarshalJSON(tt.args.jsonStr, &tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOtherFields, tt.wantOtherFields) {
				t.Errorf("UnmarshalJSON() = %s, want %s", gotOtherFields, tt.wantOtherFields)
			}
			if !reflect.DeepEqual(tt.wantObj, tt.args.obj) {
				t.Errorf("UnmarshalJSON() = %s, want %s", gotOtherFields, tt.wantOtherFields)
			}
		})
	}
}
