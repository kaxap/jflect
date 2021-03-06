package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type testFile struct {
	path string
}

var (
	url = "https://api.github.com/repos/str1ngs/gotimer"
)

var testFiles = []testFile{
	{
		path: "testdata/gotimer.json",
	},
}

func TestReflect(t *testing.T) {
	for _, f := range testFiles {
		want, err := readWant(f.path + ".want")
		if err != nil {
			t.Fatal(err)
		}
		fd, err := os.Open(f.path)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fd.Close()
		got := new(bytes.Buffer)
		err = read(fd, got)
		if err != nil {
			t.Error(err)
		}
		strWant := string(want)
		strGot := got.String()
		if strWant != strGot {
			t.Errorf("%s: want %d bytes got %d bytes", f.path, len(strWant), len(strGot))
			minLen := len(strWant)
			if minLen > len(strGot) {
				minLen = len(strGot)
			}
			for i := 0; i < minLen; i++ {
				if strWant[i] != strGot[i] {
					t.Errorf("difference starts at %d, (want: \"%s\", got: \"%s\")", i,
						strWant[i-5:i+5], strGot[i-5: i+5])

					break
				}

			}
		}
		//if !reflect.DeepEqual(want, got.Bytes()) {
		//	t.Errorf("%s: want %d bytes got %d bytes", f.path, len(want), len(got.Bytes()))
		//	fmt.Println(string(want))
		//}
	}
}

func TestReflectWithTags(t *testing.T) {
	for _, f := range testFiles {
		want, err := readWant(f.path + ".want_tags")
		if err != nil {
			t.Fatal(err)
		}
		fd, err := os.Open(f.path)
		if err != nil {
			t.Error(err)
			continue
		}
		defer fd.Close()
		got := new(bytes.Buffer)

		tagsFlag = append(tagsFlag, "db")
		err = read(fd, got)
		if err != nil {
			t.Error(err)
		}
		tagsFlag = tagsFlag[:0]

		strWant := string(want)
		strGot := got.String()
		if strWant != strGot {
			t.Errorf("%s: want %d bytes got %d bytes", f.path, len(strWant), len(strGot))
			minLen := len(strWant)
			if minLen > len(strGot) {
				minLen = len(strGot)
			}
			for i := 0; i < minLen; i++ {
				if strWant[i] != strGot[i] {
					t.Errorf("difference starts at %d, (want: \"%s\", got: \"%s\")", i,
						strWant[i-5:i+5], strGot[i-5: i+5])

					break
				}

			}
		}
	}
}

func TestSliceType(t *testing.T) {
	ty, _ := sliceType([]interface{}{})
	exp := "[]interface{}"
	if ty != exp {
		t.Fatalf("expected %s; got %s", exp, ty)
	}

	ty, _ = sliceType([]interface{}{"a", "b"})
	exp = "[]string"
	if ty != exp {
		t.Fatalf("expected %s; got %s", exp, ty)
	}

	ty, _ = sliceType([]interface{}{float64(1), float64(2)})
	exp = "[]int"
	if ty != exp {
		t.Fatalf("expected %s; got %s", exp, ty)
	}

	ty, _ = sliceType([]interface{}{"a", 1})
	exp = "[]interface{}"
	if ty != exp {
		t.Fatalf("expected %s; got %s", exp, ty)
	}

	ty, _ = sliceType([]interface{}{
		map[string]interface{}{
			"a": "aa",
			"b": "bb",
			"c": "cc",
		},
		map[string]interface{}{
			"a": "aa",
			"b": "bb",
		},
	})
	exp = "[]struct {A string `json:\"a\"`\nB string `json:\"b\"`\nC string `json:\"c\"`\n}"
	if ty != exp {
		t.Fatalf("expected %s; got %s", exp, ty)
	}
}

func TestHyphenError(t *testing.T) {
	var (
		j = `{"some-key": "foo"}`
	)
	bw := bytes.NewBufferString(j)
	out, _ := os.Open(os.DevNull)
	defer out.Close()
	err := read(bw, out)
	if err != nil {
		t.Error(err)
	}
}

func readWant(p string) ([]byte, error) {
	fd, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fd)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
