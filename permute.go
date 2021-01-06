package dwt

import (
	"bufio"
	"fmt"
	"os"
)

type File struct {
	Lines   uint32
	Indexes map[uint32]uint32 //TODO: collect offset each 1000 line for quick search
	Handler *File
	Path    string
}

type WordlistFiles []File

type WordlistPermutations struct {
	WordlistFiles WordlistFiles
	Count         uint32
	endState      []uint32
}

func (wlp *WordlistPermutations) Initialize(wordlistPaths []string) {
	for _, w := range wordlistPaths {
		wl := File{Path: w, Lines: CountLinesInFile(w)}
		wlp.WordlistFiles = append(wlp.WordlistFiles, wl)
	}
	var perms uint32
	perms = 1
	for _, w := range wlp.WordlistFiles {
		perms = perms * w.Lines
	}
	wlp.Count = perms
	wlp.endState = wlp.WordlistFiles.getPermutesEndState()
}
func (wlp WordlistPermutations) GetPermuteByState(pair []uint32) ([]string, error) {
	// Very quick!
	for i := 0; i < len(pair); i++ {
		if pair[i] > wlp.endState[i]-1 {
			return nil, fmt.Errorf("Bad permutation: %v. File: %s have only %d lines. Max value must be: %d", pair, wlp.WordlistFiles[i].Path, wlp.WordlistFiles[i].Lines, (wlp.endState[i] - 1))
		}
	}
	return wlp.WordlistFiles.getPermuteByState(pair), nil
}
func (wlp WordlistPermutations) GetPermuteByNumber(number uint32) ([]string, error) {
	// Danger: Don't use it. Very slow!
	if wlp.Count-1 < number {
		return nil, fmt.Errorf("Number of permutations is %d from 0 to %d. You want to get: %d", wlp.Count, wlp.Count-1, number)
	}
	p := make(chan []uint32, 0)
	go wlp.Permute(p, number, number+1)
	value, err := wlp.GetPermuteByState(<-p)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func (wlp WordlistPermutations) PermuteAll(linePair chan []uint32) {
	wlp.Permute(linePair, 0, 0)
}
func (wlp WordlistPermutations) Permute(linePair chan []uint32, from uint32, before uint32) {
	{
		var numberPermutations uint32
		numberPermutations = 1
		for _, countLines := range wlp.endState {
			numberPermutations *= countLines
		}
	}
	var positionCounter = make([]uint32, len(wlp.endState))
	var permuteNumber uint32
	permuteNumber = 0
	{
	loop:
		for {
			for index := len(wlp.endState) - 1; index >= 0; index-- {
				if positionCounter[index] > 0 && positionCounter[index] >= wlp.endState[index] {
					if index == 0 || (index == 1 && positionCounter[index-1] == wlp.endState[0]-1) {
						break loop
					}
					positionCounter[index] = 0
					positionCounter[index-1]++
				}
			}

			var permute []uint32
			for index, countLines := range wlp.endState {
				var position = positionCounter[index]
				if position >= 0 && position < countLines {
					permute = append(permute, position)
					if len(permute) == len(wlp.endState) {
						permuteNumber += 1
						if before > 0 {
							if permuteNumber > from && permuteNumber <= before {
								linePair <- permute
							}
							if permuteNumber == before {
								break loop
							}
						} else {
							linePair <- permute
						}
					}
				}
			}

			positionCounter[len(wlp.endState)-1]++
		}
	}
	close(linePair)
}
func (wlp WordlistPermutations) EndState() []uint32 {
	var state []uint32
	for i := 0; i < len(wlp.endState); i++ {
		state = append(state, wlp.endState[i]-1)
	}
	return state
}

func (f WordlistFiles) getPermutesEndState() []uint32 {
	var wlines []uint32
	for _, wf := range f {
		wlines = append(wlines, wf.Lines)
	}
	return wlines
}
func (f WordlistFiles) getPermuteByState(pair []uint32) []string {
	var stringPair []string
	for i, file := range f {
		lines, err := GetLine(file, pair[i], pair[i]+1)
		if err != nil {
			panic(err)
		}
		if len(lines) > 0 {
			stringPair = append(stringPair, lines[0])
		}
	}

	return stringPair
}

func CountLinesInFile(fileName string) uint32 {
	f, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("File not found: %s", fileName))
	}
	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	var count uint32
	count = 0
	for scanner.Scan() {
		count++
	}

	return count
}
func GetLine(wordlist File, from uint32, before uint32) ([]string, error) {
	var lines []string
	if before > wordlist.Lines {
		return nil, fmt.Errorf("File %s have %d lines. Requested: %d", wordlist.Path, wordlist.Lines, before)
	}
	f, err := os.Open(wordlist.Path)
	if err != nil {
		panic(fmt.Sprintf("File not found: %s", wordlist.Path))
	}
	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	var count uint32
	count = 0
	for scanner.Scan() {

		if count >= from && count < before {
			lines = append(lines, scanner.Text())
		}
		count++
	}

	return lines, nil
}
