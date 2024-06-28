package otelbaggage

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/baggage"
	"testing"
)

func TestInjectIntoBaggage(t *testing.T) {
	amount := 5
	members := testMembers(t, amount)
	ctx := context.Background()

	b, err := InjectIntoBaggage(ctx, members)
	if err != nil {
		t.Errorf("function throwed unexpected error: %s", err)
	}

	if len(b.Members()) != amount {
		t.Errorf("expected %d result but got: %d", amount, len(b.Members()))
	}
	for i := 0; i < amount; i++ {
		if b.Member(fmt.Sprintf("key_%d", i)).Value() != fmt.Sprintf("value_%d", i) {
			t.Error("the expected value of 'value' was not in the baggage")
		}
	}
}

func testMembers(t *testing.T, size ...int) []baggage.Member {
	if size == nil {
		m, err := baggage.NewMember("key_0", "value_0")
		if err != nil {
			t.Fatalf("failed to produce test data")
		}
		return []baggage.Member{m}
	}

	var result []baggage.Member
	for i := 0; i < size[0]; i++ {
		m, err := baggage.NewMember(fmt.Sprintf("%s_%d", "key", i), fmt.Sprintf("%s_%d", "value", i))
		if err != nil {
			t.Fatalf("failed to produce test data")
		}
		result = append(result, m)
	}
	return result
}
