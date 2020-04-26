package index

import (
	"reflect"
	"testing"
)

func getTestInvMap() InvMap {
	newMap := NewInvMap()
	newMap["love"] = []WordInfo{{
		Filename:  "first",
		Positions: []int{0},
	}, {
		Filename:  "second",
		Positions: []int{0},
	}}
	newMap["cats"] = []WordInfo{{
		Filename:  "first",
		Positions: []int{1},
	}}
	return newMap
}

func TestInvMap_InvertIndex(t *testing.T) {
	in := "love cats."
	filename := "filename"
	expected := NewInvMap()
	expected["love"] = []WordInfo{{
		Filename:  filename,
		Positions: []int{0},
	}}
	expected["cats"] = []WordInfo{{
		Filename:  filename,
		Positions: []int{1},
	}}
	actual := NewInvMap()
	actual.InvertIndex(in, filename)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	expected["love"] = []WordInfo{{
		Filename:  filename,
		Positions: []int{0, 2, 3},
	}}
	actual = NewInvMap()
	actual.InvertIndex("love cats love love", filename)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestGetDocStrSlice(t *testing.T) {
	in := []WordInfo{{Filename: "first_text"}, {Filename: "second_text"}}
	expected := []string{"first_text", "second_text"}
	actual := GetDocStrSlice(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestInvMap_Searcher(t *testing.T) {
	in := []string{"love"}
	expected := []MatchList{{
		Matches:  1,
		Filename: "first",
	}, {
		Matches:  1,
		Filename: "second",
	}}
	i := getTestInvMap()
	actual := i.Searcher(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestIsWordInList(t *testing.T) {
	i := getTestInvMap()
	actual, _ := i.isWordInList("love", "second")
	expected := 1
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
	actual, _ = i.isWordInList("cats", "first")
	expected = 0
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestPrepareText(t *testing.T) {
	in := "I like 254 cats, they are AWESOME!! !"
	expected := []string{"like", "cats", "awesome"}
	actual := prepareText(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}

func TestCleanText(t *testing.T) {
	in := []string{"I", "like", "cats", "they", "are", "AWESOME"}
	expected := []string{"like", "cats", "awesome"}
	actual := cleanText(in)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%v is not equal to expected %v", actual, expected)
	}
}
