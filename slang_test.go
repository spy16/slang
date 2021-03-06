package slang_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spy16/sabre"
	"github.com/spy16/slang"
)

const (
	testDir = "./lib"
	libDir  = "./lib"
)

var _ sabre.Scope = (*slang.Slang)(nil)

func TestSlang_Bind(t *testing.T) {
	sl := slang.New()

	tests := []struct {
		name    string
		symbol  string
		ns      string
		wantErr bool
	}{
		{
			name:    "CrossNamespaceBindingValidation",
			symbol:  "core/not",
			ns:      "user",
			wantErr: true,
		},
		{
			name:    "BindingInCurrentNS",
			symbol:  "hello",
			ns:      "user",
			wantErr: false,
		},
		{
			name:    "UserBinding",
			symbol:  "user/hello",
			ns:      "user",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sl.Bind(tt.symbol, sabre.Nil{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSlang_Resolve(t *testing.T) {
	sl := slang.New()
	sl.Bind("pi", sabre.Float64(3.1412))

	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{
			name:   "CoreBinding",
			symbol: "core/impl?",
		},
		{
			name:    "UserBinding",
			symbol:  "hello",
			wantErr: true,
		},
		{
			name:    "MissingUserBinding",
			symbol:  "hello",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sl.Resolve(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resolve() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSlang(t *testing.T) {
	if testing.Short() {
		return
	}

	t.Parallel()

	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}

	for _, fi := range files {
		if !strings.HasSuffix(fi.Name(), "_test.lisp") {
			continue
		}

		t.Run(fi.Name(), func(t *testing.T) {
			testFile(t, filepath.Join(testDir, fi.Name()))
		})
	}
}

func testFile(t *testing.T, file string) {
	fh, err := os.Open(file)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer fh.Close()

	sl, err := initSlang()
	if err != nil {
		t.Fatalf("failed to init slang: %v", err)
	}

	_, err = sl.ReadEval(fh)
	if err != nil {
		t.Errorf("execution failed for '%s': %v", file, err)
	}
}

func initSlang() (*slang.Slang, error) {
	di, err := ioutil.ReadDir(testDir)
	if err != nil {
		return nil, err
	}

	sl := slang.New()
	for _, fi := range di {
		if !strings.HasSuffix(fi.Name(), ".lisp") ||
			strings.HasSuffix(fi.Name(), "_test.lisp") {
			continue
		}

		fh, err := os.Open(filepath.Join(testDir, fi.Name()))
		if err != nil {
			return nil, err
		}
		defer fh.Close()

		if _, err := sl.ReadEval(fh); err != nil {
			return nil, err
		}
	}

	return sl, nil
}
