package dwt

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func TestPermuteAll(t *testing.T) {
	type args struct {
		pair     chan []uint32
		expected [][]uint32
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get all permutes",
			args: args{
				pair:     nil,
				expected: [][]uint32{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}, {0, 0, 3}, {0, 1, 0}, {0, 1, 1}, {0, 1, 2}, {0, 1, 3}, {0, 2, 0}, {0, 2, 1}, {0, 2, 2}, {0, 2, 3}, {0, 3, 0}, {0, 3, 1}, {0, 3, 2}, {0, 3, 3}, {1, 0, 0}, {1, 0, 1}, {1, 0, 2}, {1, 0, 3}, {1, 1, 0}, {1, 1, 1}, {1, 1, 2}, {1, 1, 3}, {1, 2, 0}, {1, 2, 1}, {1, 2, 2}, {1, 2, 3}, {1, 3, 0}, {1, 3, 1}, {1, 3, 2}, {1, 3, 3}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wp WordlistPermutations

			fpaths := []string{
				"./dwt/test/wl1.txt",
				"./dwt/test/wl2.txt",
				"./dwt/test/wl3.txt",
			}
			_, filename, _, _ := runtime.Caller(0)
			// The ".." may change depending on you folder structure
			dir := path.Join(path.Dir(filename), "..")
			_ = os.Chdir(dir)

			wp.Initialize(fpaths)
			p := make(chan []uint32, 0)
			var result [][]uint32
			go wp.PermuteAll(p)
			for {
				pair, ok := <-p
				if !ok {
					break
				}
				result = append(result, pair)
			}
			if !reflect.DeepEqual(tt.args.expected, result) {
				t.Errorf("expected error containing %v, got %v", tt.args.expected, result)
			}
		})
	}
}

func TestPermuteRanges(t *testing.T) {
	type args struct {
		wordlistLines []uint32
		expected      [][]uint32
		from          uint32
		before        uint32
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get permutes from 0 before 2",
			args: args{
				//permutes number	 0        1		  2		   3        4        5
				//expected: [][]int{{0,0,0}, {0,1,0}, {1,0,0}, {1,1,0}, {2,0,0}, {2,1,0}},
				expected: [][]uint32{{0, 0, 0, 0}, {0, 0, 0, 1}},
				from:     0,
				before:   2,
			},
		},
		{
			name: "Get permutes from 10000 before 10005",
			args: args{
				from:     10,
				before:   15,
				expected: [][]uint32{{0, 0, 1, 2}, {0, 0, 1, 3}, {0, 0, 1, 4}, {0, 0, 1, 5}, {0, 0, 1, 6}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wp WordlistPermutations

			fpaths := []string{
				"./dwt/test/wl1.txt",
				"./dwt/test/wl2.txt",
				"./dwt/test/wl3.txt",
				"./dwt/test/wl4.txt",
			}

			wp.Initialize(fpaths)
			p := make(chan []uint32, 0)
			var result [][]uint32
			go wp.Permute(p, tt.args.from, tt.args.before)
			for {
				pair, ok := <-p
				if !ok {
					break
				}
				result = append(result, pair)
			}
			if !reflect.DeepEqual(tt.args.expected, result) {
				t.Errorf("expected error containing %v, got %v", tt.args.expected, result)
			}
		})
	}
}

func TestCountLinesInFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "Count lines in file test/wl1.txt",
			args: args{
				"./dwt/test/wl1.txt",
			},
			want: 2,
		},
		{
			name: "Count lines in file test/wl4.txt",
			args: args{
				"./dwt/test/wl4.txt",
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountLinesInFile(tt.args.fileName); got != tt.want {
				_, filename, _, _ := runtime.Caller(0)
				// The ".." may change depending on you folder structure
				dir := path.Join(path.Dir(filename), "..")
				err := os.Chdir(dir)
				t.Errorf("CountLinesInFile() = %v, want %v %v", got, tt.want, err)
			}
		})
	}
}

func TestGetLine(t *testing.T) {
	type args struct {
		wordlist File
		from     uint32
		before   uint32
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Get lines from ../test/wl4.txt",
			args: args{
				wordlist: File{
					Path:  "/test/wl4.txt",
					Lines: 4,
				},
				from:   0,
				before: 3,
			},
			want:    []string{"1", "2", "3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, filename, _, _ := runtime.Caller(0)
			// The ".." may change depending on you folder structure
			dir := path.Join(path.Dir(filename), ".")
			err := os.Chdir(dir)
			tt.args.wordlist.Path = dir + tt.args.wordlist.Path
			got, err := GetLine(tt.args.wordlist, tt.args.from, tt.args.before)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}
