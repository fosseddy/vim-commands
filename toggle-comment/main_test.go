package main

import (
	"testing"
	"fmt"
)

type tcase struct {
	id int
	input string
	want string
}

func TestSingleLineCommentOut(t *testing.T) {
	comment := comment{
		prefix: []byte("//"),
		isMultiLine: false,
	}

	prefix := comment.prefix

	tests := []tcase{
		{1, "A\n", fmt.Sprintf("%sA\n", prefix)},
		{2, "\tA\n", fmt.Sprintf("\t%sA\n", prefix)},
		{3, "\n", fmt.Sprintf("%s\n", prefix)},
		{4, "\t\n", fmt.Sprintf("\t%s\n", prefix)},

		{5, "A\nB\n", fmt.Sprintf("%sA\n%sB\n", prefix, prefix)},
		{6, "\tA\n\tB\n", fmt.Sprintf("\t%sA\n\t%sB\n", prefix, prefix)},
		{7, "\n\n", fmt.Sprintf("%s\n%s\n", prefix, prefix)},
		{8, "\t\n\t\n", fmt.Sprintf("\t%s\n\t%s\n", prefix, prefix)},

		{9, "A\n\tB\n", fmt.Sprintf("%sA\n%s\tB\n", prefix, prefix)},
		{10, "\n\t\n", fmt.Sprintf("%s\n%s\t\n", prefix, prefix)},
		{11, "\nA\n\tB\n", fmt.Sprintf("%s\n%sA\n%s\tB\n", prefix, prefix, prefix)},
		{12, "\t\nA\n\tB\n", fmt.Sprintf("%s\t\n%sA\n%s\tB\n", prefix, prefix, prefix)},

		{13, "  A\n", fmt.Sprintf("  %sA\n", prefix)},
		{14, "  \n", fmt.Sprintf("  %s\n", prefix)},
		{15, "  A\n  B\n", fmt.Sprintf("  %sA\n  %sB\n", prefix, prefix)},
		{16, "  \n  \n", fmt.Sprintf("  %s\n  %s\n", prefix, prefix)},

		{17, "A\n  B\n", fmt.Sprintf("%sA\n%s  B\n", prefix, prefix)},
		{18, "\n  \n", fmt.Sprintf("%s\n%s  \n", prefix, prefix)},
		{19, "\nA\n  B\n", fmt.Sprintf("%s\n%sA\n%s  B\n", prefix, prefix, prefix)},
		{20, "  \nA\n  B\n", fmt.Sprintf("%s  \n%sA\n%s  B\n", prefix, prefix, prefix)},

		{21, "\t\t\tA\n\t\tB\n\tC\n", fmt.Sprintf("\t%s\t\tA\n\t%s\tB\n\t%sC\n", prefix, prefix, prefix)},
		{22, "      A\n    B\n  C\n", fmt.Sprintf("  %s    A\n  %s  B\n  %sC\n", prefix, prefix, prefix)},
	}

	for _, tc := range tests {
		input := parseInput([]byte(tc.input), &comment)
		out := []byte{}

		out = commentOutSingleLine(out, input, &comment)
		if string(out) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, out)
		}
	}
}

func TestSingleLineUncomment(t *testing.T) {
	comment := comment{
		prefix: []byte("//"),
		isMultiLine: false,
	}

	prefix := comment.prefix

	tests := []tcase{
		{1, fmt.Sprintf("%sA\n", prefix), "A\n"},
		{2, fmt.Sprintf("\t%sA\n", prefix), "\tA\n"},
		{3, fmt.Sprintf("%s\n", prefix), "\n"},
		{4, fmt.Sprintf("\t%s\n", prefix), "\t\n"},

		{5, fmt.Sprintf("%sA\n%sB\n", prefix, prefix), "A\nB\n"},
		{6, fmt.Sprintf("\t%sA\n\t%sB\n", prefix, prefix), "\tA\n\tB\n"},
		{7, fmt.Sprintf("%s\n%s\n", prefix, prefix), "\n\n"},
		{8, fmt.Sprintf("\t%s\n\t%s\n", prefix, prefix), "\t\n\t\n"},

		{9, fmt.Sprintf("%sA\n%s\tB\n", prefix, prefix), "A\n\tB\n"},
		{10, fmt.Sprintf("%s\n%s\t\n", prefix, prefix), "\n\t\n"},
		{11, fmt.Sprintf("%s\n%sA\n%s\tB\n", prefix, prefix, prefix), "\nA\n\tB\n"},
		{12, fmt.Sprintf("%s\t\n%sA\n%s\tB\n", prefix, prefix, prefix), "\t\nA\n\tB\n"},

		{13, fmt.Sprintf("  %sA\n", prefix), "  A\n"},
		{14, fmt.Sprintf("  %s\n", prefix), "  \n"},
		{15, fmt.Sprintf("  %sA\n  %sB\n", prefix, prefix), "  A\n  B\n"},
		{16, fmt.Sprintf("  %s\n  %s\n", prefix, prefix), "  \n  \n"},

		{17, fmt.Sprintf("%sA\n%s  B\n", prefix, prefix), "A\n  B\n"},
		{18, fmt.Sprintf("%s\n%s  \n", prefix, prefix), "\n  \n"},
		{19, fmt.Sprintf("%s\n%sA\n%s  B\n", prefix, prefix, prefix), "\nA\n  B\n"},
		{20, fmt.Sprintf("%s  \n%sA\n%s  B\n", prefix, prefix, prefix), "  \nA\n  B\n"},

		{21, fmt.Sprintf("\t%s\t\tA\n\t%s\tB\n\t%sC\n", prefix, prefix, prefix), "\t\t\tA\n\t\tB\n\tC\n"},
		{22, fmt.Sprintf("  %s    A\n  %s  B\n  %sC\n", prefix, prefix, prefix), "      A\n    B\n  C\n"},
	}

	for _, tc := range tests {
		input := parseInput([]byte(tc.input), &comment)
		out := []byte{}

		out = uncommentSingleLine(out, input, &comment)
		if string(out) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, out)
		}
	}
}

func TestMultiLineCommentOut(t *testing.T) {
	comment := comment{
		prefix: []byte("/*"),
		suffix: []byte("*/"),
		isMultiLine: true,
	}

	prefix := comment.prefix
	suffix := comment.suffix

	tests := []tcase{
		{1, "A\n", fmt.Sprintf("%sA%s\n", prefix, suffix)},
		{2, "\tA\n", fmt.Sprintf("\t%sA%s\n", prefix, suffix)},
		{3, "\n", fmt.Sprintf("%s%s\n", prefix, suffix)},
		{4, "\t\n", fmt.Sprintf("\t%s%s\n", prefix, suffix)},

		{5, "A\nB\n", fmt.Sprintf("%sA\nB%s\n", prefix, suffix)},
		{6, "\tA\n\tB\n", fmt.Sprintf("\t%sA\n\tB%s\n", prefix, suffix)},
		{7, "\n\n", fmt.Sprintf("%s\n%s\n", prefix, suffix)},
		{8, "\t\n\t\n", fmt.Sprintf("\t%s\n\t%s\n", prefix, suffix)},

		{9, "A\n\tB\n", fmt.Sprintf("%sA\n\tB%s\n", prefix, suffix)},
		{10, "\n\t\n", fmt.Sprintf("%s\n\t%s\n", prefix, suffix)},
		{11, "\nA\n\tB\n", fmt.Sprintf("%s\nA\n\tB%s\n", prefix, suffix)},
		{12, "\t\nA\n\tB\n", fmt.Sprintf("\t%s\nA\n\tB%s\n", prefix, suffix)},

		{13, "  A\n", fmt.Sprintf("  %sA%s\n", prefix, suffix)},
		{14, "  \n", fmt.Sprintf("  %s%s\n", prefix, suffix)},
		{15, "  A\n  B\n", fmt.Sprintf("  %sA\n  B%s\n", prefix, suffix)},
		{16, "  \n  \n", fmt.Sprintf("  %s\n  %s\n", prefix, suffix)},

		{17, "A\n  B\n", fmt.Sprintf("%sA\n  B%s\n", prefix, suffix)},
		{18, "\n  \n", fmt.Sprintf("%s\n  %s\n", prefix, suffix)},
		{19, "\nA\n  B\n", fmt.Sprintf("%s\nA\n  B%s\n", prefix, suffix)},
		{20, "  \nA\n  B\n", fmt.Sprintf("  %s\nA\n  B%s\n", prefix, suffix)},
	}

	for _, tc := range tests {
		input := parseInput([]byte(tc.input), &comment)
		out := []byte{}

		out = commentOutMultiLine(out, input, &comment)
		if string(out) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, out)
		}
	}
}

func TestMultiLineUncomment(t *testing.T) {
	comment := comment{
		prefix: []byte("/*"),
		suffix: []byte("*/"),
		isMultiLine: true,
	}

	prefix := comment.prefix
	suffix := comment.suffix

	tests := []tcase{
		{1, fmt.Sprintf("%sA%s\n", prefix, suffix), "A\n"},
		{2, fmt.Sprintf("\t%sA%s\n", prefix, suffix), "\tA\n"},
		{3, fmt.Sprintf("%s%s\n", prefix, suffix), "\n"},
		{4, fmt.Sprintf("\t%s%s\n", prefix, suffix), "\t\n"},

		{5, fmt.Sprintf("%sA\nB%s\n", prefix, suffix), "A\nB\n"},
		{6, fmt.Sprintf("\t%sA\n\tB%s\n", prefix, suffix), "\tA\n\tB\n"},
		{7, fmt.Sprintf("%s\n%s\n", prefix, suffix), "\n\n"},
		{8, fmt.Sprintf("\t%s\n\t%s\n", prefix, suffix), "\t\n\t\n"},

		{9, fmt.Sprintf("%sA\n\tB%s\n", prefix, suffix), "A\n\tB\n"},
		{10, fmt.Sprintf("%s\n\t%s\n", prefix, suffix), "\n\t\n"},
		{11, fmt.Sprintf("%s\nA\n\tB%s\n", prefix, suffix), "\nA\n\tB\n"},
		{12, fmt.Sprintf("\t%s\nA\n\tB%s\n", prefix, suffix), "\t\nA\n\tB\n"},

		{13, fmt.Sprintf("  %sA%s\n", prefix, suffix), "  A\n"},
		{14, fmt.Sprintf("  %s%s\n", prefix, suffix), "  \n"},
		{15, fmt.Sprintf("  %sA\n  B%s\n", prefix, suffix), "  A\n  B\n"},
		{16, fmt.Sprintf("  %s\n  %s\n", prefix, suffix), "  \n  \n"},

		{17, fmt.Sprintf("%sA\n  B%s\n", prefix, suffix), "A\n  B\n"},
		{18, fmt.Sprintf("%s\n  %s\n", prefix, suffix), "\n  \n"},
		{19, fmt.Sprintf("%s\nA\n  B%s\n", prefix, suffix), "\nA\n  B\n"},
		{20, fmt.Sprintf("  %s\nA\n  B%s\n", prefix, suffix), "  \nA\n  B\n"},
	}

	for _, tc := range tests {
		input := parseInput([]byte(tc.input), &comment)
		out := []byte{}

		out = uncommentMultiLine(out, input, &comment)
		if string(out) != tc.want {
			t.Errorf("\ncase: %d\ninput: %q\nwant:  %q\ngot:   %q\n", tc.id, tc.input, tc.want, out)
		}
	}
}
