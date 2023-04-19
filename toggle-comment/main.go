package main

import (
	"fmt"
	"os"
	"io"
	"bytes"
	"log"
)

type comment struct {
	prefix []byte
	suffix []byte
	isMultiLine bool
}

type input struct {
	src []byte
	len int
	lines [][]byte
	lineCount int
	smallestIndentLevel int
	isUncommenting bool
}

func main() {
	args := os.Args[1:]
	argslen := len(args)

	if argslen < 1 {
		log.Fatalf("need at least 1 argument\n")
	}

	comment := comment{
		prefix: []byte(args[0]),
		isMultiLine: false,
	}

	if argslen > 1 {
		comment.isMultiLine = true
		comment.suffix = []byte(args[1])
	}

	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	input := parseInput(src, &comment)
	output := []byte{}

	if input.isUncommenting {
		output = make([]byte, 0, input.len)

		if comment.isMultiLine {
			output = uncommentMultiLine(output, input, &comment)
		} else {
			output = uncommentSingleLine(output, input, &comment)
		}
	} else {
		output = make([]byte, 0, input.len * 3)

		if comment.isMultiLine {
			output = commentOutMultiLine(output, input, &comment)
		} else {
			output = commentOutSingleLine(output, input, &comment)
		}
	}

	fmt.Print(string(output))
}

func parseInput(src []byte, comment *comment) *input {
	input := &input{
		src: src,
		len: len(src),
		smallestIndentLevel: -1,
		isUncommenting: true,
	}

	pad := 0
	line := make([]byte, 0, 80)
	isLineStart := true
	isModeSet := false

	for i := 0; i < input.len; i++ {
		if isLineStart {
			for input.src[i] == ' ' || input.src[i] == '\t' {
				line = append(line, input.src[i])
				pad++
				i++
			}

			if !isModeSet {
				if comment.isMultiLine {
					input.isUncommenting = bytes.HasPrefix(input.src[i:], comment.prefix)
					isModeSet = true
				} else {
					// if at least one line is not commented out then we are not uncommenting
					if !bytes.HasPrefix(input.src[i:], comment.prefix) {
						input.isUncommenting = false
						isModeSet = true
					}
				}
			}

			isLineStart = false
		}

		line = append(line, input.src[i])

		if input.src[i] == '\n' {
			if input.smallestIndentLevel == -1 {
				input.smallestIndentLevel = pad
			} else {
				if !comment.isMultiLine {
					if pad < input.smallestIndentLevel {
						input.smallestIndentLevel = pad
					}
				}
			}

			input.lines = append(input.lines, line)

			pad = 0
			line = make([]byte, 0, 80)
			isLineStart = true
		}
	}

	input.lineCount = len(input.lines)

	return input
}

func commentOutMultiLine(output []byte, input *input, comment *comment) []byte {
	isPrefixInserted := false

	for lineNum, line := range input.lines {
		for i, ch := range line {
			if !isPrefixInserted && i == input.smallestIndentLevel {
				output = append(output, comment.prefix...)
				isPrefixInserted = true
			}

			if lineNum == input.lineCount - 1 && ch == '\n' {
				output = append(output, comment.suffix...)
			}

			output = append(output, ch)
		}
	}

	return output
}

func uncommentMultiLine(output []byte, input *input, comment *comment) []byte {
	isPrefixDeleted := false

	for lineNum, line := range input.lines {
		for i := 0; i < len(line); i++ {
			if !isPrefixDeleted && i == input.smallestIndentLevel {
				i += len(comment.prefix)
				isPrefixDeleted = true
			}

			if lineNum == input.lineCount - 1 && line[i] == '\n' {
				end := len(output) - len(comment.suffix)
				output = output[:end]
			}

			output = append(output, line[i])
		}
	}

	return output
}

func commentOutSingleLine(output []byte, input *input, comment *comment) []byte {
	for _, line := range input.lines {
		for i, ch := range line {
			if i == input.smallestIndentLevel {
				output = append(output, comment.prefix...)
			}
			output = append(output, ch)
		}
	}

	return output
}

func uncommentSingleLine(output []byte, input *input, comment *comment) []byte {
	for _, line := range input.lines {
		for i := 0; i < len(line); i++ {
			if i == input.smallestIndentLevel {
				i += len(comment.prefix)
			}
			output = append(output, line[i])
		}
	}

	return output
}
