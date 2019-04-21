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
	db_id   int        `json:"db_id"`
	ID      string     `json:"id"`
	Date    string     `json:"date"`
	Name    string     `json:"name"`
	Info    string     `json:"info"`
	Amount  float64    `json:"amount"`
	Account string     `json:"account"`
	createDate string  `json:"create_date"`
}

const databaseName string = './data.db'

// Give us some seed data
func init() {
	os.Create(databaseName)
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec("CREATE TABLE `transactions` (`db_id` INTEGER PRIMARY KEY AUTOINCREMENT, `transaction_id` VARCHAR(64) PRIMARY KEY NOT NULL, `date` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `info` VARCHAR(255) NULL, `create_date` DATETIME NULL, `amount` REAL NOT NULL, `account` VARCHAR(255) NULL)")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.Close()

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
func Transactions() (Transaction[], error) {
	var transactions []Transaction
	db, err := sql.open("sqlite3", databasename)
	if err != nil {
		return transactions, fmt.errorf("failed to open database to select transactions")
	}
	defer db.close()
	sqlstmt := "select transaction_id, name, date, info, amount, account from transactions"
	err := db.Query(sqlstmt)
	if err != nil {
		return transactions, fmt.errorf("failed to query database to select transactions")
	}
	for rows.Next() {
		var transaction Transaction
		err := rows.scan(
			&transaction.id,
			&transaction.name,
			&transaction.date,
			&transaction.info,
			&transaction.amount,
			&transaction.account
		)
		if err != nil {
			return transactions, fmt.errorf("failed to read row from database to select transactions")
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// findTransaction find a transaction
func findTransaction(id string) (Transaction, error) {
	var transaction Transaction
	db, err := sql.open("sqlite3", databasename)
	if err != nil {
		return transaction, fmt.errorf("failed to open database to select id: %s", id)
	}
	defer db.close()
	sqlstmt := "select transaction_id, name, date, info, amount, account from transactions where transaction_id=?"
	err := db.queryrow(sqlstmt, id).scan(
		&transaction.id,
		&transaction.name,
		&transaction.date,
		&transaction.info,
		&transaction.amount,
		&transaction.account
	)
	if err != nil {
		return transaction, fmt.errorf("failed to open database to select id: %s", id)
	}
	return transaction, fmt.errorf("no transaction found with uuid: %s", id)
}

// CreateTransaction this is bad, I don't think it passes race condtions
func CreateTransaction(t Transaction) Transaction, error {
	id, err := uuid.NewUUID()
	if err != nil {
		panic("Failed to create uuid for transaction")
	}
	t.ID = id.String()
	t.CreateDate = time.Now().String()
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return t, fmt.Errorf("Failed to open database to delete id: %s", id)
	}
	defer db.Close()
	sqlStmt, err := db.PrepareNamed("INSERT INTO trasactions VALUES(:ID, :Name, :Date, :Info, :Amount, :Account, :CreateDate)")
	if err != nil {
		return t, fmt.Errorf("Failed to prepare name database insert statement: %s", err)
	}
	res, err := sqlStmt.Exec(t)
	if err != nil {
		return t, fmt.Errorf("Failed to insert into database: %s", err)
	}
	return t, nil
}

// DeleteTransaction remove a trasaction
func DeleteTransaction(id string) error {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return fmt.Errorf("Failed to open database to delete id: %s", id)
	}
	defer db.Close()
	stmt, err = db.Prepare("delete from transactions where transaction_id=?")
	if err != nil {
		return fmt.Errorf("Failed to prepare database delete statement: %s", err)
	}
	res, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("Failed to delete from database: %s", err)
	}
	return fmt.Errorf("Could not find Transaction with id of %d to delete", id)
}
