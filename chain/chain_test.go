package chain_test

import (
	"strconv"
	"testing"

	"github.com/gopherd/exp/chain"
)

func TestChain2(t *testing.T) {
	// Create a Runnable instance that wraps a function that takes a string and returns an int.
	r1 := chain.Func(func(s string) int {
		return len(s)
	})
	// Create a Runnable instance that wraps a function that takes an int and returns a string.
	r2 := chain.Func(func(i int) string {
		return strconv.Itoa(i)
	})
	// Chain the two Runnable instances together.
	r := chain.Chain2(r1, r2)
	// Invoke the chained Runnable instance.
	out, err := r.Invoke("hello")
	if err != nil {
		t.Fatal(err)
	}
	if out != "5" {
		t.Fatalf("expected: 5, got: %s", out)
	}
}

func TestChain3(t *testing.T) {
	// Create a Runnable instance that wraps a function that takes a string and returns an int.
	r1 := chain.Func(func(s string) int {
		return len(s)
	})
	// Create a Runnable instance that wraps a function that takes an int and returns a string.
	r2 := chain.Func(func(i int) string {
		return strconv.Itoa(i)
	})
	// Create a Runnable instance that wraps a function that takes a string and returns a string.
	r3 := chain.Func2(func(s string) (int, error) {
		return strconv.Atoi(s)
	})
	// Chain the three Runnable instances together.
	r := chain.Chain3(r1, r2, r3)
	// Invoke the chained Runnable instance.
	out, err := r.Invoke("hello")
	if err != nil {
		t.Fatal(err)
	}
	if out != 5 {
		t.Fatalf("expected: 5, got: %d", out)
	}
}
