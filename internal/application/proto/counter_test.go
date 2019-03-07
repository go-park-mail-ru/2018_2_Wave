package proto

import "testing"

type counterTestPair struct {
	token      UserToken
	expectedID int64
}

func TestSimpleCountingRegister(t *testing.T) {
	var (
		counter = MakeCounter(Counting)
		pairs   = []counterTestPair{
			{UserToken("a"), 1},
			{UserToken("c"), 2},
			{UserToken("d"), 3},
			{UserToken("f"), 4},
		}
	)
	for _, pair := range pairs {
		counter.Register(pair.token)
		takenID, err := counter.GetTokenID(pair.token)
		if err != nil {
			t.Fatal(err.Error())
		}
		if takenID != pair.expectedID {
			t.Fatalf("Unexpected id %v for token %v, expected: %v", takenID, pair.token, pair.expectedID)
		}
	}
}

func TestCountingRegister(t *testing.T) {
	type counterTestPair struct {
		token      UserToken
		expectedID int64
	}
	var (
		counter = MakeCounter(Counting)
		pairs1  = []counterTestPair{
			{UserToken("a"), 1},
			{UserToken("c"), 2},
		}
		rmToken = UserToken("c")
		pairs2  = []counterTestPair{
			{UserToken("c"), 3},
			{UserToken("g"), 4},
		}
	)
	for _, pair := range pairs1 {
		counter.Register(pair.token)
		takenID, err := counter.GetTokenID(pair.token)
		if err != nil {
			t.Fatal(err.Error())
		}
		if takenID != pair.expectedID {
			t.Fatalf("Unexpected id %v for token %v, expected: %v", takenID, pair.token, pair.expectedID)
		}
	}

	counter.Unregister(rmToken)

	for _, pair := range pairs2 {
		counter.Register(pair.token)
		takenID, err := counter.GetTokenID(pair.token)
		if err != nil {
			t.Fatal(err.Error())
		}
		if takenID != pair.expectedID {
			t.Fatalf("Unexpected id %v for token %v, expected: %v", takenID, pair.token, pair.expectedID)
		}
	}
}

func TestSimpleFillGapsRegister(t *testing.T) {
	type counterTestPair struct {
		token      UserToken
		expectedID int64
	}
	var (
		counter = MakeCounter(FillGaps)
		pairs   = []counterTestPair{
			{UserToken("a"), 1},
			{UserToken("c"), 2},
			{UserToken("d"), 3},
			{UserToken("f"), 4},
		}
	)
	for _, pair := range pairs {
		counter.Register(pair.token)
		takenID, err := counter.GetTokenID(pair.token)
		if err != nil {
			t.Fatal(err.Error())
		}
		if takenID != pair.expectedID {
			t.Fatalf("Unexpected id %v for token %v, expected: %v", takenID, pair.token, pair.expectedID)
		}
	}
}

func TestFillGapsRegister(t *testing.T) {
	type counterTestPair struct {
		token      UserToken
		expectedID int64
	}
	var (
		counter = MakeCounter(FillGaps)
		pairs1  = []counterTestPair{
			{UserToken("a"), 1},
			{UserToken("c"), 2},
			{UserToken("d"), 3},
		}
		rmTokens = []UserToken{
			UserToken("d"),
			UserToken("a"),
		}
		pairs2 = []counterTestPair{
			{UserToken("f"), 1},
			{UserToken("g"), 3},
		}
	)
	for _, pair := range pairs1 {
		counter.Register(pair.token)
		takenID, err := counter.GetTokenID(pair.token)
		if err != nil {
			t.Fatal(err.Error())
		}
		if takenID != pair.expectedID {
			t.Fatalf("Unexpected id %v for token %v, expected: %v", takenID, pair.token, pair.expectedID)
		}
	}

	for _, token := range rmTokens {
		counter.Unregister(token)
	}

	for _, pair := range pairs2 {
		counter.Register(pair.token)
		takenID, err := counter.GetTokenID(pair.token)
		if err != nil {
			t.Fatal(err.Error())
		}
		if takenID != pair.expectedID {
			t.Fatalf("Unexpected id %v for token %v, expected: %v", takenID, pair.token, pair.expectedID)
		}
	}
}
