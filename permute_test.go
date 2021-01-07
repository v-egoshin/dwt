package dwt

import (
	"bufio"
	"fmt"
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

			_, filename, _, _ := runtime.Caller(0)
			// The ".." may change depending on you folder structure
			dir := path.Join(path.Dir(filename), "..")
			_ = os.Chdir(dir)
			wl := ListWordlists("./test")
			wp.Initialize(wl)
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
				expected: [][]uint32{{0, 0, 0, 0, 0}, {0, 0, 0, 0, 1}},
				from:     0,
				before:   2,
			},
		},
		{
			name: "Get permutes from 10000 before 10005",
			args: args{
				from:     10,
				before:   15,
				expected: [][]uint32{{0, 0, 0, 1, 2}, {0, 0, 0, 1, 3}, {0, 0, 0, 1, 4}, {0, 0, 0, 1, 5}, {0, 0, 0, 1, 6}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wp WordlistPermutations
			wl := ListWordlists("./test")
			wp.Initialize(wl)
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
			count, _ := CountLinesInFile(tt.args.fileName)
			if got := count; got != tt.want {
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

func CountLinesInFile2(fileName string) (uint32, map[uint32]uint32) {
	f, err := os.Open(fileName)
	defer f.Close()

	if err != nil {
		panic(fmt.Sprintf("File not found: %s", fileName))
	}
	scanner := bufio.NewScanner(f)

	index := make(map[uint32]uint32, 1)
	scanner.Split(bufio.ScanLines)

	var count, offset, offsetThousand uint32
	count = 0
	offset = 0
	offsetThousand = 1
	for scanner.Scan() {
		offset += uint32(len(scanner.Text()) + 1)
		if count%1000 == 0 {
			index[offsetThousand] = offset
			offsetThousand += 1
		}
		count++
	}
	return count, index
}

func TestCountLinesInFileWithIndex(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name  string
		args  args
		count uint32
		index uint32
	}{
		{
			name:  "File with one million lines",
			args:  args{fileName: "./test/1_000_000.txt"},
			count: 1000000,
			index: 65531,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, index := CountLinesInFile2(tt.args.fileName)
			if count != tt.count && index[10] != tt.index {
				t.Errorf("CountLinesInFile() = %v, want %v", index[10], tt.index)
			}
		})
	}
}
