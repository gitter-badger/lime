// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package commands

import (
	. "github.com/limetext/lime/backend"
	. "github.com/quarnster/util/text"
	"reflect"
	"testing"
)

type findTest struct {
	in  []Region
	exp []Region
}

func runFindTest(command string, tests *[]findTest, t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	defer w.Close()

	v := w.NewFile()
	defer func() {
		v.SetScratch(true)
		v.Close()
	}()

	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	for i, test := range *tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, command, nil)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Test %d failed: %v, %+v", i, sr, test)
		}
	}
}

func TestSingleSelection(t *testing.T) {
	tests := []findTest{
		{
			[]Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}},
			[]Region{{1, 1}},
		},
		{
			[]Region{{2, 2}, {3, 3}, {6, 6}},
			[]Region{{2, 2}},
		},
		{
			[]Region{{5, 5}},
			[]Region{{5, 5}},
		},
	}

	runFindTest("single_selection", &tests, t)
}

func TestFindUnderExpand(t *testing.T) {
	tests := []findTest{
		{
			[]Region{{0, 0}},
			[]Region{{0, 5}},
		},
		{
			[]Region{{19, 20}},
			[]Region{{19, 20}, {22, 23}},
		},
	}

	runFindTest("find_under_expand", &tests, t)
}

func TestSelectAll(t *testing.T) {
	tests := []findTest{
		{
			[]Region{{1, 1}, {2, 2}, {3, 3}, {6, 6}},
			[]Region{{0, 36}},
		},
		{
			[]Region{{2, 2}, {3, 3}, {6, 6}},
			[]Region{{0, 36}},
		},
		{
			[]Region{{5, 5}},
			[]Region{{0, 36}},
		},
	}

	runFindTest("select_all", &tests, t)
}
