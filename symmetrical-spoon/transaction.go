package main

import (
	"os"
	"fmt"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)

// Transaction the datastructure to hold a transaction
type Transaction struct {
	id              int        `json:"id"`
	Identifier      string     `json:"identifier"`
	Date            string     `json:"date"`
	Name            string     `json:"name"`
	Info            string     `json:"info"`
	Amount          float64    `json:"amount"`
	Account         string     `json:"account"`
	createDate      string     `json:"create_date"`
}

// Give us some seed data
func init() {
	databaseName := os.Getenv("DATABASE")
	os.Create(fmt.Sprintf("./%s", databaseName))
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE `transactions` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `identifier` VARCHAR(64) NOT NULL, `date` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `info` VARCHAR(255) NULL, `amount` REAL NOT NULL, `account` VARCHAR(255) NULL, `create_date` DATETIME NULL)")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var t Transaction = Transaction{
		Name: "symmetrical-spoon",
		Date: time.Now().String(),
		Info: "VISA",
		Amount: 0.0,
		Account: "checking",
	}
	CreateTransaction(t)
}

// Transactions show all transactions
func Transactions() ([]Transaction, error) {
	var transactions []Transaction
	db, err := sql.Open("sqlite3", os.Getenv("database_name"))
	if err != nil {
		return transactions, fmt.Errorf("failed to open database to select transactions")
	}
	defer db.Close()
	sqlstmt := "select transaction_id, name, date, info, amount, account from transactions"
	rows, err := db.Query(sqlstmt)
	if err != nil {
		return transactions, fmt.Errorf("failed to query database to select transactions")
	}
	defer rows.Close()
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.id,
			&transaction.Identifier,
			&transaction.Name,
			&transaction.Date,
			&transaction.Info,
			&transaction.Amount,
			&transaction.Account,
		)
		if err != nil {
			return transactions, fmt.Errorf("failed to read row from database to select transactions")
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// findTransaction find a transaction
func findTransaction(identifier string) (Transaction, error) {
	var transaction Transaction
	db, err := sql.Open("sqlite3", os.Getenv("database_name"))
	if err != nil {
		return transaction, fmt.Errorf("failed to open database to select identifier: %s", identifier)
	}
	defer db.Close()
	sqlstmt := "select identifier, name, date, info, amount, account from transactions where identifier=?"
	err = db.QueryRow(sqlstmt, identifier).Scan(
		&transaction.Identifier,
		&transaction.Name,
		&transaction.Date,
		&transaction.Info,
		&transaction.Amount,
		&transaction.Account,
	)
	if err != nil {
		return transaction, fmt.Errorf("failed to open database to select identifier: [%s] %s", identifier, err)
	}
	return transaction, nil
}

// CreateTransaction this is bad, I don't think it passes race condtions
func CreateTransaction(t Transaction) (Transaction, error) {
	identifier, err := uuid.NewUUID()
	if err != nil {
		panic("Failed to create uuid for transaction")
	}
	t.Identifier = identifier.String()
	t.createDate = time.Now().String()
	db, err := sql.Open("sqlite3", os.Getenv("database_name"))
	if err != nil {
		return t, fmt.Errorf("Failed to open database to delete id: %s", t.Identifier)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return t, fmt.Errorf("failed to begin database to insert identifier: [%s] %s", t.Identifier, err)
	}
	defer tx.Rollback()
	sqlStmt, err := tx.Prepare("INSERT INTO transactions VALUES(:Identifier, :Name, :Date, :Info, :Amount, :Account, :createDate)")
	if err != nil {
		return t, fmt.Errorf("Failed to prepare name database insert statement: %s", err)
	}
	_, err = sqlStmt.Exec(t)
	if err != nil {
		return t, fmt.Errorf("Failed to insert into database: %s", err)
	}
	err = tx.Commit()
	if err != nil {
		return t, fmt.Errorf("failed to commit database to insert identifier: %s", t.Identifier)
	}
	return t, nil
}

// DeleteTransaction remove a transaction
func DeleteTransaction(identifier string) error {
	db, err := sql.Open("sqlite3", os.Getenv("database_name"))
	if err != nil {
		return fmt.Errorf("Failed to open database to delete id: %s", identifier)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin database to insert identifier: %s", identifier)
	}
	defer tx.Rollback()
	stmt, err := db.Prepare("delete from transactions where transaction_id=?")
	if err != nil {
		return fmt.Errorf("Failed to prepare database delete statement: %s", err)
	}
	_, err = stmt.Exec(identifier)
	if err != nil {
		return fmt.Errorf("Failed to delete from database: %s", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit database to insert identifier: %s", identifier)
	}
	return nil
}
