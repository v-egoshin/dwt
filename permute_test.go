package distributed_word_tool

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func TestPermuteAll(t *testing.T) {
	type args struct {
		parts    []int
		pair     chan []int
		expected [][]int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get all permutes for 2 files with two lines",
			args: args{
				parts:    []int{2, 2},
				pair:     nil,
				expected: [][]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			},
		}, {
			name: "Get all permutes for 1 file with 3 lines",
			args: args{
				parts:    []int{3},
				pair:     nil,
				expected: [][]int{{0}, {1}, {2}},
			},
		}, {
			name: "Get all permutes for 3 files",
			args: args{
				parts:    []int{3, 2, 1},
				pair:     nil,
				expected: [][]int{{0, 0, 0}, {0, 1, 0}, {1, 0, 0}, {1, 1, 0}, {2, 0, 0}, {2, 1, 0}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(chan []int, 0)
			var result [][]int
			go PermuteAll(tt.args.parts, p)
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
		wordlistLines []int
		expected      [][]int
		from          int
		before        int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get permutes from 0 before 2",
			args: args{
				wordlistLines: []int{3, 2, 1},
				//permutes number	 0        1		  2		   3        4        5
				//expected: [][]int{{0,0,0}, {0,1,0}, {1,0,0}, {1,1,0}, {2,0,0}, {2,1,0}},
				expected: [][]int{{0, 0, 0}, {0, 1, 0}},
				from:     0,
				before:   2,
			},
		},
		{
			name: "Get permutes from 10000 before 10005",
			args: args{
				wordlistLines: []int{10, 10, 10, 10, 10},
				from:          1000,
				before:        1005,
				expected:      [][]int{{0, 1, 0, 0, 0}, {0, 1, 0, 0, 1}, {0, 1, 0, 0, 2}, {0, 1, 0, 0, 3}, {0, 1, 0, 0, 4}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := make(chan []int, 0)
			var result [][]int
			go Permute(tt.args.wordlistLines, p, tt.args.from, tt.args.before)
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
		want int
	}{
		{
			name: "Count lines in file test/wl1.txt",
			args: args{
				"test/wl1.txt",
			},
			want: 2,
		},
		{
			name: "Count lines in file test/wl4.txt",
			args: args{
				"test/wl4.txt",
			},
			want: 4,
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
		from     int
		before   int
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
