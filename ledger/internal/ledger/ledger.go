package ledger

import "database/sql"

func Insert(db *sql.DB, order_id, user_id string, amount int64, operation string, transactionDate string) error {
	stmt, err := db.Prepare("INSERT INTO ledger (order_id, user_id, amount, operation, transactionDate) VALUES (?,?,?,?,?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(order_id, user_id, amount, operation, transactionDate)
	if err != nil {
		return err
	}

	return nil
}
