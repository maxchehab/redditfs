package main

import (
	"fmt"
)

func main() {
	// session := GetSessionWithCredentials(TEST_USERNAME, TEST_PASSWORD)
	manifest, err := RetreiveManifestFromReddit("77346c3e708a")

	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("%+v\n", manifest)
	}
}
