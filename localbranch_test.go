package main

import (
	"testing"
)

func TestLocalBranchCount(t *testing.T) {
	t.Run("Local branch count test", func(t *testing.T) {
		count := 1
		npmObj := &npmObject{}
		localBranchCount(count, npmObj)

		if npmObj.numBranches != "2" {
			t.Errorf("Expected branch count to be 2, but got %s", npmObj.numBranches)
		}
	})
}
