package dwt

import (
	"bufio"
	"fmt"
	"os"
)

type File struct {
	Path  string
	Lines int
}

type WordlistFiles []File

type WordlistPermutations struct {
	WordlistFiles WordlistFiles
	Count         int
	endState      []int
}

func (wlp *WordlistPermutations) Initialize(wordlistPaths []string) {
	for _, w := range wordlistPaths {
		wl := File{Path: w, Lines: CountLinesInFile(w)}
		wlp.WordlistFiles = append(wlp.WordlistFiles, wl)
	}

	perms := 1
	for _, w := range wlp.WordlistFiles {
		perms = perms * w.Lines
	}
	wlp.Count = perms
	wlp.endState = wlp.WordlistFiles.getPermutesEndState()
}
func (wlp WordlistPermutations) GetPermuteByState(pair []int) ([]string, error) {
	// Very quick!
	for i := 0; i < len(pair); i++ {
		if pair[i] > wlp.endState[i]-1 {
			return nil, fmt.Errorf("Bad permutation: %v. File: %s have only %d lines. Max value must be: %d", pair, wlp.WordlistFiles[i].Path, wlp.WordlistFiles[i].Lines, (wlp.endState[i] - 1))
		}
	}
	return wlp.WordlistFiles.getPermuteByState(pair), nil
}
func (wlp WordlistPermutations) GetPermuteByNumber(number int) ([]string, error) {
	// Danger: Don't use it. Very slow!
	if wlp.Count-1 < number {
		return nil, fmt.Errorf("Number of permutations is %d from 0 to %d. You want to get: %d", wlp.Count, wlp.Count-1, number)
	}
	p := make(chan []int, 0)
	go wlp.Permute(p, number, number+1)
	value, err := wlp.GetPermuteByState(<-p)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func (wlp WordlistPermutations) PermuteAll(linePair chan []int) {
	wlp.Permute(linePair, 0, 0)
}
func (wlp WordlistPermutations) Permute(linePair chan []int, from int, before int) {
	{
		var numberPermutations = 1
		for _, countLines := range wlp.endState {
			numberPermutations *= countLines
		}
	}
	var positionCounter = make([]int, len(wlp.endState))
	permuteNumber := 0
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

			var permute []int
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
func (wlp WordlistPermutations) EndState() []int {
	var state []int
	for i := 0; i < len(wlp.endState); i++ {
		state = append(state, wlp.endState[i]-1)
	}
	return state
}

func (f WordlistFiles) getPermutesEndState() []int {
	var wlines []int
	for _, wf := range f {
		wlines = append(wlines, wf.Lines)
	}
	return wlines
}
func (f WordlistFiles) getPermuteByState(pair []int) []string {
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

func CountLinesInFile(fileName string) int {
	f, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Sprintf("File not found: %s", fileName))
	}
	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}
func GetLine(wordlist File, from int, before int) ([]string, error) {
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

	count := 0
	for scanner.Scan() {

		if count >= from && count < before {
			lines = append(lines, scanner.Text())
		}
		count++
	}

	return lines, nil
}
