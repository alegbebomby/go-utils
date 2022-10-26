package library

import (
	"bitbucket.org/dnda-tech/api-and-billing/app/constants"
	"bitbucket.org/dnda-tech/api-and-billing/app/logger"
	"bitbucket.org/dnda-tech/api-and-billing/app/models"
	"database/sql"
	"fmt"
	"github.com/logrusorgru/aurora"
	"strings"
	"time"
)

const tableFieldCreated = "created"
const tableFieldClientID = "client_id"
const tableFieldStatus = "status"
const tableFieldUserID = "user_id"


//IsValidSenderID checks if the sender belongs to client
func IsValidSenderID(db *sql.DB, clientID int64, senderID string) string {

	var senderId sql.NullString

	dbUtil := Db{DB: db}
	queryString := "SELECT s.sender_name FROM client_sender cs INNER JOIN sender s ON cs.sender_id = s.id WHERE cs.client_id = ? AND cs.status = 1 AND s.status = 1 AND UPPER(s.sender_name) = ? "
	dbUtil.SetQuery(queryString)
	dbUtil.SetParams(clientID, strings.ToUpper(strings.TrimSpace(senderID)))

	
	err := dbUtil.FetchOne().Scan(&senderId)
	if err != nil {

		logger.Error("error checking if client is assigned sender ID %s",err.Error())
	}

	return senderId.String
}

//checked
//GetMessageMapping gets message to be stripped off
func ISOtp(db *sql.DB, clientID int64, senderID string) bool {

	var otp sql.NullInt64

	dbUtil := Db{DB: db}
	dbUtil.SetQuery("SELECT cs.otp FROM client_sender cs INNER JOIN sender s ON cs.sender_id=s.id WHERE cs.client_id = ? AND UPPER(c.sender_id) = ? ")
	dbUtil.SetParams(clientID, strings.ToUpper(strings.TrimSpace(senderID)))

	err := dbUtil.FetchOne().Scan(&otp)
	if err == sql.ErrNoRows {

		return false
	}

	if otp.Int64 == int64(1) {

		return true
	}

	return false
}

func CreateRootUser(db *sql.DB) error {

	msisdn := 254700000000
	// create client billing
	data := map[string]interface{}{
		"name": "Admin Account",
		"address": "Root",
		"msisdn": msisdn,
		"email": "admin@example.com",
		tableFieldStatus: 1,
		"created_by": 0,
		tableFieldClientID: 0,
		tableFieldCreated: ToMysqlDateTime(time.Now()),
	}

	dbObject := &Db{DB:db}

	clientID, err := dbObject.Upsert("client",data,nil)
	if err != nil {
		logger.GetLog().Errorf("error creating root client %s",err.Error())
		return err
	}

	data = map[string]interface{}{
		tableFieldClientID: clientID,
		"billing_mode": "POSTPAID",
		"credit_limit": 1000000,
		tableFieldCreated: ToMysqlDateTime(time.Now()),
	}

	_, _ = dbObject.Upsert("client_billing",data,nil)

	// create client billing
	data = map[string]interface{}{
		tableFieldClientID: clientID,
		"balance":   1000000,
		tableFieldCreated:   ToMysqlDateTime(time.Now()),
	}

	_, _ = dbObject.Upsert("client_balance", data, nil)

	data = map[string]interface{}{
		"full_name":    "Root",
		"email":        "admin@example.com",
		"msisdn":       msisdn,
		"country_code": "ke",
		tableFieldStatus: 1,
		"role_id":      1,
		tableFieldCreated:      ToMysqlDateTime(time.Now()),
	}

	userID, err := dbObject.Upsert("user", data, nil)
	if err != nil {
		logger.GetLog().Errorf("error creating root client %s",err.Error())
		return err
	}

	pass := "500t@P1ssM0rd"
	hash, _ := Hash(pass)

	data = map[string]interface{}{
		tableFieldUserID: userID,
		"password": hash,
		"password_reset_code":  "",
		"next_password_change": ToMysqlDateTime(time.Now().Add(time.Duration(1000) * 24 * time.Hour)),
		tableFieldCreated:  ToMysqlDateTime(time.Now()),
	}

	_, err = dbObject.Upsert("user_login", data,[]string{tableFieldUserID})
	if err != nil {
		logger.GetLog().Errorf("error creating root client %s",err.Error())
		return err
	}

	dbObject.Upsert("password_history", map[string]interface{}{
		tableFieldUserID:  userID,
		"password": MD5S(pass),
		tableFieldCreated:  ToMysqlDateTime(time.Now()),
	}, nil)

	data = map[string]interface{}{
		tableFieldClientID: clientID,
		tableFieldUserID:   userID,
		"role_id":   1,
		tableFieldStatus:    1,
		tableFieldCreated:   ToMysqlDateTime(time.Now()),
	}

	_, err = dbObject.Upsert("client_user", data, nil)
	if err != nil {
		logger.GetLog().Errorf("error creating root client %s",err.Error())
		return err
	}

	logger.GetLog().Infof("successfully created root account, root username %d password %s",msisdn,pass)
	return nil

}

func GetClientBillingByClientID(db *sql.DB, clientID int64) (*models.ClientBilling, error) {

	var account models.ClientBilling

	rows := db.QueryRow("SELECT cb.id,cb.client_id,cb.billing_mode,cb.credit_limit,cb.billing_type,b.balance,b.bonus FROM client_billing cb INNER JOIN client_balance b ON cb.client_id = b.client_id WHERE cb.client_id = ? ", clientID)

	err := rows.Scan(
		&account.Id,
		&account.ClientID,
		&account.BillingMode,
		&account.CreditLimit,
		&account.BillingType,
		&account.Balance,
		&account.Bonus)

	if err != nil {

		logger.GetLog().Errorf("database error %s ", err.Error())
		return nil, err
	}

	return &account, nil
}

func ChargeClient(db *sql.DB, clientID int64, amount float64) (status int, statusText string) {

	clientBilling, err := GetClientBillingByClientID(db, clientID)
	if err != nil {

		logger.GetLog().Errorf(" got error retrieving client details %s ", aurora.Red(err.Error()))
		return constants.StatusInternalServerError, err.Error()
	}

	if clientBilling.BillingMode.String == "PREPAID" && amount > clientBilling.Balance.Float64 {

		msg := fmt.Sprintf("Your balance of %0.f is insufficient to send sms at a total cost of KES %0.f ", clientBilling.Balance.Float64, amount)
		return constants.StatusInsufficientBalance, msg
	}

	if amount > clientBilling.Balance.Float64 && amount > clientBilling.CreditLimit.Float64 {

		msg := fmt.Sprintf("Your balance of %0.f is insufficient to send sms at a total cost of KES %0.f ", clientBilling.Balance.Float64, amount)
		return constants.StatusInsufficientBalance, msg
	}

	dbutils := Db{DB: db}

	dbutils.SetQuery("UPDATE client_balance SET balance = balance - ? WHERE client_id = ? LIMIT 1 ")
	dbutils.SetParams(amount, clientID)

	_, err = dbutils.UpdateQuery()
	if err != nil {

		logger.GetLog().Errorf(" got error executing database query %s ", aurora.Red(err.Error()))
		return constants.StatusInternalServerError, err.Error()
	}

	return constants.StatusSend, "client billed successfully"

}

func GetSMSCostPerConnector(db *sql.DB, clientID, connectorID int64) float64 {

	var billingRate sql.NullFloat64

	err := db.QueryRow("SELECT billing_rate FROM client_sms_billing WHERE client_id = ? AND connector_id = ? ", clientID, connectorID).Scan(&billingRate)

	if err != nil {

		logger.GetLog().Errorf(" got error retrieving billing rate %s ", err.Error())
		billingRate.Float64 = 1
	}

	return billingRate.Float64

}
