package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type inputTest struct {
	Name, Input string
}

var inputTests = []inputTest{
	{
		Name:  "ascii: symbols",
		Input: "$%'()*+,-./<>:;=",
	},
	{
		Name:  "unicode: emoji",
		Input: "😇🤖💝🙇🏿‍♂️",
	},
	{
		Name:  "unicode: text",
		Input: "Testing «ταБЬℓσ»: 1<2 & 4+1>3, now 20% off!",
	},
}

func InputTests(delim byte, unsafeChars ...byte) []inputTest {
	t := make([]inputTest, len(inputTests))
	copy(t, inputTests)
	for i := 1; i <= 8; i++ {
		t = append(t, inputTest{
			Name:  fmt.Sprintf("delim=%q;i=%d;beginning", delim, i),
			Input: strings.Repeat(string(delim), i) + "bar",
		})
		t = append(t, inputTest{
			Name:  fmt.Sprintf("delim=%q;i=%d;middle", delim, i),
			Input: "foo" + strings.Repeat(string(delim), i) + "bar",
		})
		t = append(t, inputTest{
			Name:  fmt.Sprintf("delim=%q;i=%d;end", delim, i),
			Input: "foo" + strings.Repeat(string(delim), i),
		})
	}
	for _, c := range append([]byte("\t\n \"'\\"), unsafeChars...) {
		for i := 1; i <= 4; i++ {
			cs := strings.Repeat(string(c), i)

			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;beginning", cs),
				Input: cs + "bar",
			})
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;middle", cs),
				Input: "foo" + cs + "bar",
			})
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;end", cs),
				Input: "foo" + cs,
			})

			ds := string(delim)
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;ds=%q;beginning", cs, ds),
				Input: cs + ds + "bar",
			})
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;ds=%q;middle", cs, ds),
				Input: "foo" + cs + ds + "bar",
			})
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;ds=%q;end", cs, ds),
				Input: "foo" + cs + ds,
			})

			ds += string(delim)
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;ds=%q;beginning", cs, ds),
				Input: cs + ds + "bar",
			})
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;ds=%q;middle", cs, ds),
				Input: "foo" + cs + ds + "bar",
			})
			t = append(t, inputTest{
				Name:  fmt.Sprintf("cs=%q;ds=%q;end", cs, ds),
				Input: "foo" + cs + ds,
			})
		}
	}
	return t
}

func TestExecOutput(t *testing.T, expected, name string, args ...string) {
	out, cmd, err := Output(name, args...)
	if err != nil {
		t.Fatalf("Cmd.Output() = _, %v; want nil\nCmd: %v", err, cmd)
	}
	TestOutput(t, cmd, expected, string(out))
}

func TestOutput(t *testing.T, cmd []string, expected, actual string) {
	if sdiff := cmp.Diff(expected, actual); sdiff != "" {
		bdiff := cmp.Diff([]byte(expected), []byte(actual))
		t.Errorf("Output mismatch:\nCmd: %v\nString (-want +got):\n%s\nBytes (-want +got):\n%s", cmd, sdiff, bdiff)
	}
}

func TestDiff(t *testing.T, name, expected, actual string) {
	if sdiff := cmp.Diff(expected, actual); sdiff != "" {
		bdiff := cmp.Diff([]byte(expected), []byte(actual))
		t.Errorf("%s mismatch:\nString (-want +got):\n%s\nBytes (-want +got):\n%s", name, sdiff, bdiff)
	}
}

func init() {
	b := make([]byte, 255)
	for i := 1; i <= 255; i++ {
		b[i-1] = byte(i)
	}
	inputTests = append(inputTests, inputTest{
		Name:  "bytes: 1-255",
		Input: string(b),
	})
}
