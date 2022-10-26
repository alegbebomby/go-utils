package library

import (
	"bitbucket.org/dnda-tech/api-and-billing/app/logger"
	"database/sql"
	"strconv"
	"strings"
)

func CreateCashout(db *sql.DB, clientID, msisdn, amount int64, reference, name, description string) (int64, error) {

	paymentMethod := "DIRECT-TOPUP"
	names := strings.Split(name, " ")

	firstName := ""
	middleName := ""
	lastName := ""

	if len(names) > 0 {

		firstName = names[0]
	}

	if len(names) > 1 {

		middleName = names[1]
	}

	if len(names) > 2 {

		lastName = strings.Join(names[2:], " ")
	}

	data := map[string]interface{}{
		"client_id":               clientID,
		"amount":                  amount,
		"reference":               reference,
		"msisdn":                  msisdn,
		"first_name":              firstName,
		"middle_name":             middleName,
		"other_name":              lastName,
		"payment_method":          paymentMethod,
		"description":             description,
		"status":                  1,
		"callback_url":            "",
		"external_transaction_id": "",
		"created": MysqlNow(),
	}

	dbUtil := Db{DB: db}
	id, err :=dbUtil.Upsert("cashout",data,nil)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func CreateTransaction(db *sql.DB, clientID int64, amount float64, createdBy, transactionType int64, description string) (int64, error) {

	data := map[string]interface{}{
		"client_id":        clientID,
		"transaction_type": transactionType,
		"amount":           amount,
		"description":      description,
		"created_by":       createdBy,
		"created": MysqlNow(),
	}

	dbUtil := Db{DB: db}
	id, err := dbUtil.Upsert("transactions",data,nil)
	if err != nil {
		logger.GetLog().Errorf("error creating transaction %s",err.Error())
		return 0, err
	}

	return id, nil
}

func ReferenceNumber(int642 int64) string {

	int642 = int642 + 1000*1000 // make 3 digits
	hex := strings.ToUpper(strconv.FormatInt(int642, 36))
	return strings.ToUpper(hex)

}
