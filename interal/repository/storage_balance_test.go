package repository_test

import (
	"testing"
)

func TestStorageBalance_GetBalanceByLogin(t *testing.T) {
	// mockDB, err := pgxmock.NewConn()
	// if err != nil {
	// 	fmt.Println("failed to open pgxmock database:", err)
	// }
	// defer mockDB.Close(context.Background())
	// // 1. Define expected columns and rows
	// rows := pgxmock.NewRows([]string{"user_id", "current_balance", "withdrawn_balance", "updated_at"}).AddRow(23, 32, 24, 34645645645)
	// // 2. Set expectation: Expect a specific SELECT query and return rows
	// // Use ExpectQuery for SELECT statements
	// mockDB.ExpectQuery("SELECT user_id, current_balance, withdrawn_balance, updated_at FROM balances WHERE user_id = \\?").WithArgs(23).WillReturnRows(rows)

	// // 3. Call your production code
	// storageBalance := repository.NewBalanceStorage(mockDB)
	// uid := "23"
	// balance := storageBalance.GetBalanceByLogin(context.Background(), uid)

	// // 4. Assertions
	// assert.NoError(t, err)
	// assert.Equal(t, "23", balance.UserID)
	// assert.NoError(t, mockDB.ExpectationsWereMet()) // Ensure all mocks were called
}
