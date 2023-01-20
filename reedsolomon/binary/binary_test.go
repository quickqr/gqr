package binary

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		booleans []bool
	}
	tests := []struct {
		name string
		args args
		want *Binary
	}{
		{
			name: "case 0",
			args: args{
				// 10001 0x88
				booleans: []bool{true, false, false, false, true},
			},
			want: &Binary{
				bits:    []byte{0x88},
				lenBits: 5,
			},
		},
		{
			name: "case 1",
			args: args{
				// 1111 1111 1111 0xff 0xf0
				booleans: []bool{true, true, true, true, true, true, true, true, true, true, true, true},
			},
			want: &Binary{
				bits:    []byte{0xff, 0xf0},
				lenBits: 12,
			},
		},
		{
			name: "case 2",
			args: args{
				// 0
				booleans: []bool{},
			},
			want: &Binary{
				bits:    []byte{},
				lenBits: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.booleans...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromBinaryString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Binary
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				s: "10011000 1111",
			},
			want: &Binary{
				bits:    []byte{0x98, 0xf0},
				lenBits: 12,
			},
			wantErr: false,
		},
		{
			name: "case 1, invalid char",
			args: args{
				s: "10011000 c111",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromBinaryString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromBinaryString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromBinaryString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinary_ensureCapacity(t *testing.T) {
	type args struct {
		numBits int
	}
	tests := []struct {
		name string
		b    *Binary
		args args
	}{
		{
			name: "case 0, 0 numBits",
			b:    New(),
			args: args{
				numBits: 0,
			},
		},
		{
			name: "case 1, 8 numBits",
			b:    New(),
			args: args{
				numBits: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.ensureCapacity(tt.args.numBits)
			if cap(tt.b.bits) < tt.args.numBits {
				t.Errorf("could not ensureCapcity")
			}
		})
	}
}

func TestBinary_At(t *testing.T) {
	type args struct {
		pos int
	}
	tests := []struct {
		name string
		b    *Binary
		args args
		want bool
	}{
		{
			name: "case 0, pos 0 of len 4",
			b:    New(true, false, true, false), // 1010
			args: args{pos: 0},
			want: true,
		},
		// {
		// 	name: "case 1, pos 4 of len 4",
		// 	b:    New(true, false, true, false), // 1010, this should panic or not?
		// 	args: args{pos: 4},
		// 	want: false,
		// },
		{
			name: "case 2, pos 3 of len 4",
			b:    New(true, false, true, false), // 1010
			args: args{pos: 3},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.At(tt.args.pos); got != tt.want {
				t.Errorf("Binary.At() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinary_Subset(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name    string
		b       *Binary
		args    args
		want    *Binary
		wantErr bool
	}{
		{
			name: "case 0",
			b:    New(true, true, true, true, false, false, false, false, true, false), // 0xf0 0x80, (0b1111000010)
			args: args{
				start: 0,
				end:   10,
			},
			want:    New(true, true, true, true, false, false, false, false, true, false),
			wantErr: false,
		},
		{
			name: "case 1",
			b:    New(true, true, true, true, false, false, false, false, true, false), // 0xf0 0x80, (0b1111000010)
			args: args{
				start: 4,
				end:   9,
			},
			want:    New(false, false, false, false, true),
			wantErr: false,
		},
		{
			name: "case 2",
			b:    New(true, true, true, true, false, false, false, false, true, false), // 0xf0 0x80, (0b1111000010)
			args: args{
				start: 3,
				end:   8,
			},
			want:    New(true, false, false, false, false),
			wantErr: false,
		},
		// {
		// 	name: "case 2",
		// 	// 0xf0 0x80, (0b1111000010)
		// 	b: New(true, true, true, true, false, false, false, false, true, false),
		// 	args: args{
		// 		start: 3,
		// 		end:   14,
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Subset(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Binary.Subset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.EqualTo(tt.want) {
				t.Errorf("Binary.Subset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinary_Append(t *testing.T) {
	type args struct {
		other *Binary
	}
	tests := []struct {
		name    string
		b       *Binary
		args    args
		equalTo string
	}{
		{
			name: "case 0",
			b:    New(true),
			args: args{
				other: New(false, false, false),
			},
			equalTo: "1000",
		},
		{
			name: "case 0",
			b:    New(true),
			args: args{
				other: New(false, false, false),
			},
			equalTo: "1000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Append(tt.args.other)
			e, _ := NewFromBinaryString(tt.equalTo)
			if !tt.b.EqualTo(e) {
				t.Errorf("Binary.Append() = %v, want %v", tt.b, e)
			}
		})
	}
}

func TestBinary_AppendBytes(t *testing.T) {
	type args struct {
		byts []byte
	}
	tests := []struct {
		name string
		b    *Binary
		args args
		want *Binary
	}{
		{
			name: "case 0",
			b: &Binary{
				bits:    []byte{0xff, 0xc0},
				lenBits: 12,
			},
			args: args{
				byts: []byte{0xff, 0x00},
			},
			want: &Binary{
				bits:    []byte{0xff, 0xcf, 0xf0, 0x00},
				lenBits: 28,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.AppendBytes(tt.args.byts...)
			if !tt.want.EqualTo(tt.b) {
				t.Errorf("Binary.AppendBytes want: %v, got: %v", tt.want, tt.b)
			}
		})
	}
}

func TestBinary_AppendByte(t *testing.T) {
	type args struct {
		byt     byte
		numBits int
	}
	tests := []struct {
		name    string
		b       *Binary
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "case 0",
			b:    New(),
			args: args{
				byt:     0xff,
				numBits: 8,
			},
			want:    "1111 1111",
			wantErr: false,
		},
		{
			name: "case 1",
			b:    New(),
			args: args{
				byt:     0xff,
				numBits: 0,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "case 2",
			b:    New(),
			args: args{
				byt:     0xff,
				numBits: 9,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "case 3",
			b:    New(true, false, true),
			args: args{
				byt:     0xff,
				numBits: 3,
			},
			want:    "101 111",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.AppendByte(tt.args.byt, tt.args.numBits); (err != nil) != tt.wantErr {
				t.Errorf("Binary.AppendByte() error = %v, wantErr %v", err, tt.wantErr)
			}
			if wantB, _ := NewFromBinaryString(tt.want); !wantB.EqualTo(tt.b) {
				t.Errorf("Binary.AppendByte() got = %v, want %v", tt.b, wantB)
			}
		})
	}
}

func TestBinary_AppendBools(t *testing.T) {
	type args struct {
		booleans []bool
	}
	tests := []struct {
		name    string
		b       *Binary
		args    args
		wantLen int
	}{
		{
			name: "case 0",
			b:    New(),
			args: args{
				booleans: []bool{true, true, false},
			},
			wantLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.b.AppendBools(tt.args.booleans...); tt.b.Len() != tt.wantLen {
				t.Errorf("Binary.AppendBools() wantLen: %d, but %d", tt.wantLen, tt.b.Len())
			}
		})
	}
}

func TestBinary_AppendNumBools(t *testing.T) {
	type args struct {
		num     int
		boolean bool
	}
	tests := []struct {
		name string
		b    *Binary
		args args
		want string
	}{
		{
			name: "case 0",
			b:    New(false, true, false),
			args: args{
				num:     4,
				boolean: true,
			},
			want: "0101111",
		},
		{
			name: "case 0",
			b:    New(false, true, false),
			args: args{
				num:     3,
				boolean: false,
			},
			want: "010000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.AppendNumBools(tt.args.num, tt.args.boolean)
			if wantB, _ := NewFromBinaryString(tt.want); !wantB.EqualTo(tt.b) {
				t.Errorf("Binary.AppendNumBools() want: %v, but %v", wantB, tt.b)
			}
		})
	}
}

func TestBinary_VisitAll(t *testing.T) {
	type args struct {
		f IterFunc
	}
	tests := []struct {
		name string
		b    *Binary
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.VisitAll(tt.args.f)
		})
	}
}

func TestBinary_String(t *testing.T) {
	tests := []struct {
		name string
		b    *Binary
		want string
	}{
		{
			name: "case 0",
			b:    New(true, false, true, false),
			want: "Binary length: 4, bits: 1010",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Binary.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinary_Len(t *testing.T) {
	tests := []struct {
		name string
		b    *Binary
		want int
	}{
		{
			name: "case 0",
			b:    New([]bool{true, false, true, false, true, false}...),
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Len(); got != tt.want {
				t.Errorf("Binary.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinary_EqualTo(t *testing.T) {
	type args struct {
		other string
	}
	tests := []struct {
		name string
		b1   *Binary
		b2   *Binary
		want bool
	}{
		{
			name: "case 0",
			b1: &Binary{
				bits:    []byte{0x88, 0xc0},
				lenBits: 10,
			},
			b2: &Binary{
				bits:    []byte{0x88, 0xc1},
				lenBits: 10,
			},
			want: true,
		},
		{
			name: "case 1",
			b1: &Binary{
				bits:    []byte{0x88, 0xc1},
				lenBits: 16,
			},
			b2: &Binary{
				bits:    []byte{0xff, 0xc1},
				lenBits: 16,
			},
			want: false,
		},
		{
			name: "case 2",
			b1: &Binary{
				bits:    []byte{0x88, 0xc1},
				lenBits: 16,
			},
			b2: &Binary{
				bits:    []byte{0x88, 0xc1},
				lenBits: 10,
			},
			want: false,
		},
		{
			name: "case 2",
			b1: &Binary{
				bits:    []byte{0x88, 0xc1},
				lenBits: 16,
			},
			b2: &Binary{
				bits:    []byte{0x88, 0xcf},
				lenBits: 16,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b1.EqualTo(tt.b2); got != tt.want {
				t.Errorf("Binary.EqualTo() want: %v, got: %v, b1: %v, b2: %v", tt.want, got, tt.b1, tt.b2)
			}
		})
	}
}

func TestBinary_Copy(t *testing.T) {
	type fields struct {
		bits    []byte
		lenBits int
	}
	tests := []struct {
		name   string
		fields fields
		want   *Binary
	}{
		{
			name: "case 0",
			fields: fields{
				bits:    []byte{0x12, 0x34},
				lenBits: 16,
			},
			want: &Binary{
				bits:    []byte{0x12, 0x34},
				lenBits: 16,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binary{
				bits:    tt.fields.bits,
				lenBits: tt.fields.lenBits,
			}
			got := b.Copy()
			b.AppendByte(0x33, 8)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Binary.Copy() = %v, want %v", got, tt.want)
			}

			if b.EqualTo(got) {
				t.Errorf("Binary.Copy() changed b: %v, but copy %v changed too", b, got)
			}
		})
	}
}

func TestBinary_AppendUint32(t *testing.T) {
	type fields struct {
		bits    []byte
		lenBits int
	}
	type args struct {
		value   uint32
		numBits int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Binary
	}{
		{
			name: "case 0",
			fields: fields{
				bits:    []byte{0x22},
				lenBits: 8,
			},
			args: args{
				value:   0xf12f,
				numBits: 9,
			},
			want: &Binary{
				bits:    []byte{0x22, 0x97, 0x80},
				lenBits: 17,
			},
		},
		{
			name: "case 2",
			fields: fields{
				bits:    []byte{0x22},
				lenBits: 8,
			},
			args: args{
				value:   0xf12f,
				numBits: 13,
			},
			want: &Binary{
				bits:    []byte{0x22, 0x89, 0x78},
				lenBits: 21,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binary{
				bits:    tt.fields.bits,
				lenBits: tt.fields.lenBits,
			}
			t.Logf("origin binary: %v", b)
			b.AppendUint32(tt.args.value, tt.args.numBits)
			if !b.EqualTo(tt.want) {
				t.Errorf("Binary.AppendUint32(): %v, want: %v", b, tt.want)
			}
		})
	}
}
