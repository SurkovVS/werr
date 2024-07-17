package werr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *werror
	}{
		{
			name: "new error",
			args: args{
				err: errors.New("test"),
			},
			want: &werror{
				status:   0,
				wrapping: errors.New("test"),
				original: errors.New("test"),
			},
		},
		{
			name: "existing werr",
			args: args{
				err: &werror{
					status:   99,
					wrapping: errors.New("test"),
					original: errors.New("test"),
				},
			},
			want: &werror{
				status:   99,
				wrapping: errors.New("test"),
				original: errors.New("test"),
			},
		},
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: &werror{
				status:   0,
				wrapping: nil,
				original: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_Error(t *testing.T) {
	type fields struct {
		status   int
		wrapping error
		original error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		recF   func()
	}{
		{
			name: "filled",
			fields: fields{
				status:   0,
				wrapping: fmt.Errorf("paper: %w", errors.New("rock")),
				original: nil,
			},
			want: fmt.Errorf("paper: %w", errors.New("rock")).Error(),
			recF: func() {
				if r := recover(); r != nil {
					t.Fatal("unexpected panic")
				}
			},
		},
		{
			name: "filled with status",
			fields: fields{
				status:   200,
				wrapping: fmt.Errorf("paper: %w", errors.New("rock")),
				original: nil,
			},
			want: fmt.Errorf("(werr status - 200 OK) paper: %w", errors.New("rock")).Error(),
			recF: func() {
				if r := recover(); r != nil {
					t.Fatal("unexpected panic")
				}
			},
		},
		{
			name: "empty",
			fields: fields{
				status:   0,
				wrapping: nil,
				original: nil,
			},
			want: errors.New("").Error(),
			recF: func() {
				if r := recover(); r == nil {
					t.Fatal("expected panic")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status:   tt.fields.status,
				wrapping: tt.fields.wrapping,
				original: tt.fields.original,
			}
			defer tt.recF()
			if got := werr.Error(); got != tt.want {
				t.Errorf("werror.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_SetStatus(t *testing.T) {
	type fields struct {
		status   int
		wrapping error
		original error
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *werror
	}{
		{
			name:   "empty",
			fields: fields{},
			args: args{
				status: 100,
			},
			want: &werror{
				status:   100,
				wrapping: nil,
				original: nil,
			},
		},
		{

			name: "filled",
			fields: fields{
				status:   200,
				wrapping: errors.New("test"),
				original: errors.New("test"),
			},
			args: args{
				status: 100,
			},
			want: &werror{
				status:   100,
				wrapping: errors.New("test"),
				original: errors.New("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status:   tt.fields.status,
				wrapping: tt.fields.wrapping,
				original: tt.fields.original,
			}
			if got := werr.SetStatus(tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("werror.SetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_Status(t *testing.T) {
	type fields struct {
		status   int
		wrapping error
		original error
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "empty",
			fields: fields{},
			want:   0,
		},
		{
			name: "filled",
			fields: fields{
				status:   100,
				wrapping: errors.New("test"),
				original: errors.New("test"),
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status:   tt.fields.status,
				wrapping: tt.fields.wrapping,
				original: tt.fields.original,
			}
			if got := werr.Status(); got != tt.want {
				t.Errorf("werror.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_Wrap(t *testing.T) {
	type fields struct {
		status   int
		wrapping error
		original error
	}
	type args struct {
		t string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *werror
	}{
		{
			name: "filled",
			fields: fields{
				status:   0,
				wrapping: errors.New("test"),
				original: nil,
			},
			args: args{
				t: "paper",
			},
			want: &werror{
				status:   0,
				wrapping: fmt.Errorf("%s: %w", "paper", errors.New("test")),
				original: nil,
			},
		},
		{
			name:   "empty struct",
			fields: fields{},
			args: args{
				t: "paper",
			},
			want: &werror{
				status:   0,
				wrapping: fmt.Errorf("paper: %w", nil),
				original: nil,
			},
		},
		{
			name: "empty text",
			fields: fields{
				status:   0,
				wrapping: errors.New("test"),
				original: nil,
			},
			args: args{
				t: "",
			},
			want: &werror{
				status:   0,
				wrapping: fmt.Errorf(": %w", errors.New("test")),
				original: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status:   tt.fields.status,
				wrapping: tt.fields.wrapping,
				original: tt.fields.original,
			}
			if got := werr.Wrap(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("werror.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_Unwrap(t *testing.T) {
	type fields struct {
		status   int
		wrapping error
		original error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "simple",
			fields: fields{
				status:   0,
				wrapping: fmt.Errorf("%s: %w", "paper", errors.New("test")),
				original: nil,
			},
			want: errors.New("test"),
		},
		{
			name: "wrap werr",
			fields: fields{
				status:   0,
				wrapping: fmt.Errorf("%s: %w", "paper", New(errors.New("test"))),
				original: nil,
			},
			want: &werror{
				status:   0,
				wrapping: errors.New("test"),
				original: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status:   tt.fields.status,
				wrapping: tt.fields.wrapping,
				original: tt.fields.original,
			}
			got := werr.Unwrap()
			// if !ok {
			// 	t.Error("unexpected type, want *werr")
			// }
			if got.Error() != tt.want.Error() {
				t.Errorf("werror.Wrap() = %v, want %v", got.Error(), tt.want.Error())
			}
		})
	}
}

func Test_werror_Original(t *testing.T) {
	type fields struct {
		status   int
		wrapping error
		original error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "filled",
			fields: fields{
				original: errors.New("test"),
			},
			wantErr: true,
		},
		{
			name:    "empty",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status:   tt.fields.status,
				wrapping: tt.fields.wrapping,
				original: tt.fields.original,
			}
			if err := werr.Original(); (err != nil) != tt.wantErr {
				t.Errorf("werror.Original() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
