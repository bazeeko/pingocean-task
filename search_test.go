package main

import (
	"context"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	testCases := []struct {
		searchStr      string
		urls           []string
		expectedOutput map[string]int
	}{
		{
			searchStr:      "google",
			urls:           []string{"https://google.com/", "https://yandex.ru/"},
			expectedOutput: map[string]int{"https://google.com/": 76, "https://yandex.ru/": 0},
		},
	}

	for _, tc := range testCases {
		searchResult := Search(tc.searchStr, tc.urls, context.Background())

		if !reflect.DeepEqual(searchResult, tc.expectedOutput) {
			t.Errorf("TestSearch: got %v, want %v", searchResult, tc.expectedOutput)
		}
	}
}
