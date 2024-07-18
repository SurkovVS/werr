package werr

import (
	"errors"
	"fmt"
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
				status: 0,
				err:    errors.New("test"),
			},
		},
		{
			name: "existing werr",
			args: args{
				err: &werror{
					status: 99,
					err:    errors.New("test"),
				},
			},
			want: &werror{
				status: 99,
				err:    errors.New("test"),
			},
		},
		{
			name: "nil",
			args: args{
				err: nil,
			},
			want: &werror{
				status: 0,
				err:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.err)
			if got.err != nil && tt.want.err != nil {
				if got.status != tt.want.status ||
					got.err.Error() != tt.want.err.Error() {
					t.Errorf("New() = %v, want %v", got, tt.want)
				}
			} else {
				if got.err != tt.want.err {
					t.Error("unexpected nil result")
				}

			}
		})
	}
}

func Test_werror_Error(t *testing.T) {
	type fields struct {
		status   int
		err      error
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
				err:      fmt.Errorf("paper: %w", errors.New("rock")),
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
				err:      fmt.Errorf("paper: %w", errors.New("rock")),
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
				err:      nil,
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
				status: tt.fields.status,
				err:    tt.fields.err,
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
		err      error
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
				status: 100,
				err:    nil,
			},
		},
		{

			name: "filled",
			fields: fields{
				status:   200,
				err:      errors.New("test"),
				original: errors.New("test"),
			},
			args: args{
				status: 100,
			},
			want: &werror{
				status: 100,
				err:    errors.New("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status: tt.fields.status,
				err:    tt.fields.err,
			}
			if got := werr.SetStatus(tt.args.status); got.status != tt.want.status {
				t.Errorf("werror.SetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_Status(t *testing.T) {
	type fields struct {
		status   int
		err      error
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
				err:      errors.New("test"),
				original: errors.New("test"),
			},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status: tt.fields.status,
				err:    tt.fields.err,
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
		err      error
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
				err:      errors.New("test"),
				original: nil,
			},
			args: args{
				t: "paper",
			},
			want: &werror{
				status: 0,
				err:    fmt.Errorf("%s: %w", "paper", errors.New("test")),
			},
		},
		{
			name:   "empty struct",
			fields: fields{},
			args: args{
				t: "paper",
			},
			want: &werror{
				status: 0,
				err:    fmt.Errorf("paper: %w", nil),
			},
		},
		{
			name: "empty text",
			fields: fields{
				status:   0,
				err:      errors.New("test"),
				original: nil,
			},
			args: args{
				t: "",
			},
			want: &werror{
				status: 0,
				err:    fmt.Errorf(": %w", errors.New("test")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status: tt.fields.status,
				err:    tt.fields.err,
			}
			if got := werr.Wrap(tt.args.t); got.err.Error() != tt.want.err.Error() {
				t.Errorf("werror.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_werror_Unwrap(t *testing.T) {
	type fields struct {
		status   int
		err      error
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
				err:      fmt.Errorf("%s: %w", "paper", errors.New("test")),
				original: nil,
			},
			want: errors.New("test"),
		},
		{
			name: "wrap werr",
			fields: fields{
				status:   0,
				err:      fmt.Errorf("%s: %w", "paper", New(errors.New("test"))),
				original: nil,
			},
			want: &werror{
				status: 0,
				err:    errors.New("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			werr := &werror{
				status: tt.fields.status,
				err:    tt.fields.err,
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
