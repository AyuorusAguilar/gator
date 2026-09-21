package rss

import (
	"context"
	"fmt"
	"testing"
)

func TestRSS(t *testing.T) {
	result, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		t.Fatalf("Error:\n\t%v\n", err)
	}
	fmt.Printf("Result:\n\t%v\n", result)
}