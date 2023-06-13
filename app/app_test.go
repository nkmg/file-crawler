package app

import (
	"log"
	"os"
	"testing"
)

type testCase struct {
	searchingWord string
	filesExpected []string
}

func TestSearching(t *testing.T) {
	if err := os.Mkdir("app_test", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	os.WriteFile("app_test/apple.txt", []byte("banana"), 0644)
	os.WriteFile("app_test/banana.txt", []byte("pineapple"), 0644)
	os.WriteFile("app_test/pineapple.txt", []byte("grape"), 0644)
	os.WriteFile("app_test/grape.txt", []byte("strawberry"), 0644)
	os.WriteFile("app_test/grape_apple.txt", []byte("pineapple"), 0644)
	os.WriteFile("app_test/strawberry.txt", []byte("papaya"), 0644)
	os.WriteFile("app_test/papaya.txt", []byte("banana"), 0644)
	os.WriteFile("app_test/papaya_apple.txt", []byte("grape"), 0644)
	os.WriteFile("app_test/melon.txt", []byte("grape banana apple"), 0644)

	defer os.RemoveAll("app_test/")

	testCases := []testCase{
		{"banana", []string{"apple.txt", "banana.txt", "papaya.txt", "melon.txt"}},
		{"apple", []string{"apple.txt", "grape_apple.txt", "papaya_apple.txt", "melon.txt", "pineapple.txt"}},
		{"grape", []string{"grape_apple.txt", "papaya_apple.txt", "pineapple.txt", "grape.txt", "melon.txt"}},
		{"melon", []string{"melon.txt"}},
		{"papaya", []string{"strawberry.txt", "papaya.txt", "papaya_apple.txt"}},
	}

	for _, testingCases := range testCases {
		t.Run(testingCases.searchingWord, func(t *testing.T) {
			filesFounded := searching("app_test/", testingCases.searchingWord)

			if len(filesFounded) != len(testingCases.filesExpected) {
				t.Error("Files founded is not the same the files expected. Fail in search for: ", testingCases.searchingWord)
			} else {
				countFiles := 0
				for _, filename := range testingCases.filesExpected {
					for _, filesAux := range filesFounded {
						if filesAux == filename {
							countFiles++
						}
					}
				}
				if countFiles != len(testingCases.filesExpected) {
					t.Error("Files founded is not the same the files expected. Fail in search for: ", testingCases.searchingWord)
				}
			}
		})
	}
}
