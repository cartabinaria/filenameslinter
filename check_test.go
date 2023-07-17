package filenameslinter

import (
	"testing"

	"github.com/csunibo/synta"
	"github.com/liamg/memoryfs"
	"github.com/stretchr/testify/assert"
)

// Don't want to use storage for tests, so we use /tmp
// but it's only linux compatible.

var rootFS = memoryfs.New()

var _ = rootFS.MkdirAll("dir1", 0777)
var _ = rootFS.WriteFile("abc", []byte(""), 0777)
var _ = rootFS.WriteFile("APPUNTI.md", []byte(""), 0777)
var _ = rootFS.WriteFile("basic.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove", 0777)
var _ = rootFS.WriteFile("prove/basic123.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove2", 0777)
var _ = rootFS.WriteFile("prove2/a1a.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove3", 0777)
var _ = rootFS.WriteFile("prove3/basi___c123.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove3/pippo.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove4", 0777)
var _ = rootFS.MkdirAll("prove4/prova4", 0777)
var _ = rootFS.WriteFile("prove4/prova4/asdiw_sdfws.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove4/prova4/__addi11223sjdK___.txt", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove5", 0777)
var _ = rootFS.MkdirAll("prove5/pro__asd", 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/samu.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/luca.md", []byte(""), 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/fabio123.md", []byte(""), 0777)
var _ = rootFS.WriteFile("prove5/pro__asd/ANDREA.pdf", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove6", 0777)
var _ = rootFS.MkdirAll("prove6/cart-ella", 0777)
var _ = rootFS.WriteFile("prove6/cart-ella/b.txt", []byte(""), 0777)
var _ = rootFS.WriteFile("prove6/cart-ella/f.md", []byte(""), 0777)
var _ = rootFS.WriteFile("prove6/cart-ella/x.pdf", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove7", 0777)
var _ = rootFS.MkdirAll("prove7/.test-dir!!", 0777)
var _ = rootFS.WriteFile("prove7/.test-file!!", []byte(""), 0777)
var _ = rootFS.MkdirAll("prove7/subdir", 0777)
var _ = rootFS.MkdirAll("prove7/subdir/.test-dir!!", 0777)
var _ = rootFS.WriteFile("prove7/subdir/.test-file!!", []byte(""), 0777)

func TestDirEmpty(t *testing.T) {
	synta := synta.MustSynta(`useless = a|dir1
> useless.useless`)
	err := CheckDir(synta, rootFS, "dir1", &Options{true, true, false})
	assert.Nil(t, err)
}

func TestDirOneFileCorrect(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove", &Options{true, true, false})
	assert.Nil(t, err)
}

func TestDirOneFileNotCorrect(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove2", &Options{true, true, false})
	assert.NotNil(t, err)
}

func TestDirFileCorrectAndNotCorrect(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove3", &Options{true, true, false})
	assert.NotNil(t, err)
	matchErr := err.(RegexMatchError)
	assert.Equal(t, "basi___c123.txt", matchErr.Filename)
}

func TestDirRecursiveWithCorrectDirName(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove4", &Options{true, true, false})
	assert.Nil(t, err)
}

func TestDirRecursiveWithNotCorrectDirName(t *testing.T) {
	input := `word = [a-zA-Z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove5", &Options{true, true, false})
	assert.NotNil(t, err)
}

func TestDirRecursiveWithNotCorrectButKebabDirName(t *testing.T) {
	input := `word = [a-z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove6", &Options{true, true, false})
	assert.Nil(t, err)
}

func TestDirIgnoreDots(t *testing.T) {
	input := `word = [a-z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove7", &Options{false, true, true})
	assert.Nil(t, err)
	err = CheckDir(synta, rootFS, "prove7", &Options{false, true, false})
	assert.NotNil(t, err)
}

func TestDirIgnoreDotsRecursive(t *testing.T) {
	input := `word = [a-z]+[0-9]*
ext = pdf|txt|tex|md
> word.ext`

	synta := synta.MustSynta(input)
	err := CheckDir(synta, rootFS, "prove7", &Options{true, true, true})
	assert.Nil(t, err)
	err = CheckDir(synta, rootFS, "prove7", &Options{false, true, true})
	assert.Nil(t, err)
	err = CheckDir(synta, rootFS, "prove7", &Options{true, true, false})
	assert.NotNil(t, err)
}
