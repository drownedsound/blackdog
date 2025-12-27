package infra

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	app "github.com/drownedsound/blackdog/internal/features/application"
)

var fixedNow time.Time

func TestMain(m *testing.M) {
	// Freeze time at Jan 1, 2025, 12:00 UTC
	fixedNow = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	os.Exit(m.Run())
}

func TestMockDb_SaveValidation(t *testing.T) {
	repo := NewMockDb()
	ctx := context.Background()

	t.Run("Save Returns Error When Id Is Invalid", func(t *testing.T) {
		// 1. Create app with Id = 0 (default)
		invalidApp := &app.Application{
			MemberReferenceNo: "REF-001",
			// Id is explicitly 0
		}

		// 2. Expect ErrInvalidUUID
		// This hits the 'if a.Id <= 0' branch in mockdb.go
		if err := repo.Save(ctx, invalidApp); err != app.ErrInvalidUUID {
			t.Errorf(
				"Save Error: Expected = %v, Actual = %v",
				app.ErrInvalidUUID,
				err,
			)
		}
	})

	t.Run("Save Auto-Populates CreatedAt", func(t *testing.T) {
		// 1. Create app with Zero CreatedAt
		freshApp := &app.Application{
			Id:                101,
			MemberReferenceNo: "REF-002",
			// CreatedAt is left as time.Time{} (Zero)
		}

		if err := repo.Save(ctx, freshApp); err != nil {
			t.Fatalf("Save Error: %v", err)
		}

		// 2. Verify CreatedAt was set
		// This hits the 'if a.CreatedAt.IsZero()' branch in mockdb.go
		if freshApp.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be populated, got Zero")
			t.Errorf(
				"Save Error on Refetched App: "+
					"CreatedAt Expected = 0, Actual = %v",
				freshApp.CreatedAt,
			)
		}
	})
}

func TestMockDb_SaveAndGetApplication(t *testing.T) {
	repo := NewMockDb()
	ctx := context.Background()

	newApp := app.Application{
		CreatedAt:         fixedNow,
		UpdatedAt:         fixedNow,
		OtherApplicants:   []app.Applicant{},
		MemberReferenceNo: "ABCDE12345",
		Applicant: app.Applicant{
			Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
			ContactNumbers: []app.ContactNumber{
				{Value: "9171234567", Type: app.TypeMobile},
			},
			LastName:    "Smith",
			FirstName:   "John",
			MiddleName:  "Doe",
			IsPrincipal: true,
		},
		CreditCard: app.CreditCard{
			ProfileId:    1,
			CurrencyId:   1,
			CreditLimit:  1_000_000,
			InterestRate: 300,
		},
		Id:              7410680147088510976,
		RequestedAmount: 100_000_000,
		Status:          app.StatusCreated,
	}

	t.Run("Save Creates Valid Application", func(t *testing.T) {
		if err := repo.Save(ctx, &newApp); err != nil {
			t.Errorf("Save Error on New App: %v", err)
		}
	})

	var fetchedApp app.Application

	t.Run("GetById Retrieves Saved Application", func(t *testing.T) {
		var err error
		fetchedApp, err = repo.GetById(ctx, newApp.Id)
		if err != nil {
			t.Errorf("GetById Error on New App: %v", err)
		}

		if fetchedApp.RequestedAmount != 100_000_000 {
			t.Errorf(
				"GetById Error on Fetched App: "+
					"RequestedAmount Expected = %v, Actual = %v",
				100_000_000,
				fetchedApp.RequestedAmount,
			)
		}
	})

	// Verify Memory Isolation (Deep Copy)
	// This ensures that modifying the struct returned by GetById does not
	// corrupt the data inside the MockDb (common bug in in-memory mocks).
	t.Run("Verify_Memory_Isolation", func(t *testing.T) {
		if fetchedApp.Id == 0 {
			t.Skip("Skipping Isolation Test Because Previous Fetch Failed")
		}

		// Mutate the local copy
		fetchedApp.RequestedAmount = 99999

		// Fetch a fresh copy from the repo
		refetch, err := repo.GetById(ctx, newApp.Id)
		if err != nil {
			t.Errorf("GetById Error on Refetched App: %v", err)
		}

		// The repo should still have the OLD value
		if refetch.RequestedAmount == 99999 {
			t.Error(
				"Violation of Isolation: Modifying Local Struct Affected" +
					" The Repository Store",
			)
		}
	})

	// Verify Update
	t.Run("Persist_And_Verify_Update", func(t *testing.T) {
		if fetchedApp.Id == 0 {
			t.Skip("Skipping Update Test Because Struct is Missing")
		}

		// Save the mutated 'fetchedApp' (which has Amount = 99999)
		if err := repo.Save(ctx, &fetchedApp); err != nil {
			t.Errorf("Save Error on Fetched App: %v", err)
		}

		// Fetch again to verify persistence
		refetchAfterUpdate, err := repo.GetById(ctx, newApp.Id)
		if err != nil {
			t.Errorf("GetById Error Refetched App: %v", err)
		}

		if refetchAfterUpdate.RequestedAmount != 99_999 {
			t.Errorf(
				"GetById Error on Refetched App: "+
					"RequestedAmount Expected = %v, Actual = %v",
				99_999,
				refetchAfterUpdate.RequestedAmount,
			)
		}
	})
}

func TestMockDb_GetUsingInvalidId(t *testing.T) {
	repo := NewMockDb()
	_, err := repo.GetById(context.Background(), -100)
	if err != app.ErrInvalidUUID {
		t.Errorf("GetById Error: %v", err)
	}
}

func TestMockDb_ApplicationNotFound(t *testing.T) {
	repo := NewMockDb()
	_, err := repo.GetById(context.Background(), 999)
	if err != app.ErrNotFound {
		t.Errorf("GetById Error: %v", err)
	}
}

func TestMockDb_ConcurrentSaveAndGetApplication(t *testing.T) {
	repo := NewMockDb()
	ctx := context.Background()

	app := app.Application{
		CreatedAt:         fixedNow,
		UpdatedAt:         fixedNow,
		OtherApplicants:   []app.Applicant{},
		MemberReferenceNo: "ABCDE12345",
		Applicant: app.Applicant{
			Birthday: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
			ContactNumbers: []app.ContactNumber{
				{Value: "9171234567", Type: app.TypeMobile},
			},
			LastName:    "Smith",
			FirstName:   "John",
			MiddleName:  "Doe",
			IsPrincipal: true,
		},
		CreditCard: app.CreditCard{
			ProfileId:    1,
			CurrencyId:   1,
			CreditLimit:  1_000_000,
			InterestRate: 300,
		},
		Id:              7410680147088510976,
		RequestedAmount: 100_000_000,
		Status:          app.StatusCreated,
	}

	repo.Save(ctx, &app)
	targetId := app.Id

	// Run parallel readers and writers
	start := make(chan struct{})
	var wg sync.WaitGroup

	// 10 Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			a, _ := repo.GetById(ctx, targetId)
			a.UpdatedAt = time.Now()
			repo.Save(ctx, &a)
		}()
	}

	// 10 Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			repo.GetById(ctx, targetId)
		}()
	}

	close(start)
	wg.Wait()
}

func TestMockDb_SerializationErrors(t *testing.T) {
	repo := NewMockDb()
	ctx := context.Background()

	t.Run("Save Returns Error on Json.Marshal Failure", func(t *testing.T) {
		// 1. Construct an app with a Year outside [0, 9999]
		//    to trigger MarshalJSON error.
		//
		// 2. Set Id > 0 to simulate an UPDATE.
		//    If Id is 0, Save() overwrites CreatedAt with time.Now(),
		//    masking the error.
		invalidApp := &app.Application{
			Id:                1,
			CreatedAt:         time.Date(10_001, 1, 1, 0, 0, 0, 0, time.UTC),
			MemberReferenceNo: "INVALID",
		}

		if err := repo.Save(ctx, invalidApp); err == nil {
			t.Error(
				"Expected Json.Marshal Error Due to Invalid Year, Got Nil",
			)
		}
	})

	t.Run(
		"GetById Returns Error on Json.Unmarshal Failure",
		func(t *testing.T) {
			// 1. Manually inject corrupted JSON into the private store.
			badID := int64(999)
			repo.store[badID] = []byte(`{"truncated_json":`)

			// 2. Attempt to retrieve it.
			_, err := repo.GetById(ctx, badID)
			if err == nil {
				t.Error(
					"Expected Json.Unmarshal Error" +
						"Due to Corrupt Data, Got Nil",
				)
			}
		})
}

func TestMockDb_SimulatedConnectionError(t *testing.T) {
	repo := NewMockDb()
	ctx := context.Background()

	// Trigger the specific "ERR-100" condition
	a := &app.Application{
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		MemberReferenceNo: "ERR-100",
		Status:            app.StatusCreated,
		RequestedAmount:   10_0000_000,
	}

	err := repo.Save(ctx, a)
	if err != app.ErrConnectionRefused {
		t.Errorf(
			"Expected ErrConnectionRefused for ERR-100, Got %v",
			err,
		)
	}
}
