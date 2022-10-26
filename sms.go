package library

import (
	"bitbucket.org/dnda-tech/api-and-billing/app/database"
	"bitbucket.org/dnda-tech/api-and-billing/app/models"
	"database/sql"
	"errors"
	phonenumbers "github.com/mudphilo/phonenumber"
	"os"
)

func SendOTP(msisdn string, message string) error {

	senderId := os.Getenv("PLATFORM_SENDER_ID")

	mnc := phonenumbers.GetMSISDN(msisdn)

	if mnc.Msisdn == 0 {

		return errors.New("invalid phone number")

	}

	outbox := models.Outbox{
		UserID:      1,
		ClientID:    1,
		Msisdn:      msisdn,
		Message:     message,
		ShortCode:   senderId,
		Priority:    1,
		SMSCharged:  0,
		NumberOfSMS: int64(CalculateTotalPages(len(message), 160)),
	}

	queueName := "transactional_sms"

	conn := database.GetNatsConnection()
	defer conn.Close()

	// completed
	return PublishToNats(conn,"sms", queueName, outbox)

}

func GetQueuePrefix(db *sql.DB,mcc,mnc string) (connector int64, queue string, err error)  {

	// GET Connector ID
	var connectorID1 sql.NullInt64
	var queuePrefix1 sql.NullString

	dbUtil := Db{DB: db}
	dbUtil.SetQuery("SELECT id,queue_prefix FROM sms_connector WHERE mcc = ? AND mnc = ? ")
	dbUtil.SetParams(mcc,mnc)

	err = dbUtil.FetchOne().Scan(&connectorID1,&queuePrefix1)
	if err != nil && err != sql.ErrNoRows {

		return 0,  "", err

	}

	if connectorID1.Int64 == 0 {

		dbUtil.SetQuery("SELECT id,queue_prefix FROM sms_connector WHERE mcc = ? AND scope = 'COUNTRY' ")
		dbUtil.SetParams(mcc)

		// get country connecor
		err = dbUtil.FetchOne().Scan(&connectorID1,&queuePrefix1)
		if err != nil && err != sql.ErrNoRows {

			return 0,  "", err
		}

		// get default connector
		if connectorID1.Int64 == 0 {

			dbUtil.SetQuery("SELECT id,queue_prefix FROM sms_connector WHERE scope = 'DEFAULT' ")
			dbUtil.SetParams()

			err = dbUtil.FetchOne().Scan(&connectorID1,&queuePrefix1)
			if err != nil  {

				return 0,  "", err

			}
		}

	}

	return connectorID1.Int64, queuePrefix1.String, nil

}
