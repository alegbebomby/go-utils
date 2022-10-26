package library

import (
	"bitbucket.org/dnda-tech/api-and-billing/app/constants"
	"bitbucket.org/dnda-tech/api-and-billing/app/database"
	"bitbucket.org/dnda-tech/api-and-billing/app/logger"
	"bitbucket.org/dnda-tech/api-and-billing/app/models"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLibrary(t *testing.T) {

	dbName := fmt.Sprintf(constants.TestDB,strings.ToLower(randomdata.Letters(5)))

	_ = os.Setenv("DB_HOST", constants.DbHost)
	_ = os.Setenv("DB_PASSWORD", constants.DbPd)
	_ = os.Setenv("DB_USERNAME", constants.DbUser)
	_ = os.Setenv("DB_PORT", constants.DbPort)
	_ = os.Setenv("DARIDE_JWT_SECRET", "eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTY2MzAxODY0NSwiaWF0IjoxNjYzMDE4NjQ1fQ.q9SwFW4jkhSpQKupbFOZVwdzQKnnsI73BZJZT-lDr1E")
	_ = os.Setenv("DARIDE_JWT_ISSUER", "da-ride.com")
	_ = os.Setenv("DARIDE_JWT_DURATION_HOURS", "72")
	_ = os.Setenv("NATS_URI", "nats.dev.smespay.com:4222")
	_ = os.Setenv("ENV","tests")
	_ = os.Setenv("DB_NAME",dbName)
	_ = os.Setenv("SENDGRID_API_KEY", "0")
	_ = os.Setenv("PLATFORM_SENDER_ID", "SMSAfrica")
	masterPass := randomdata.Letters(60)
	_ = os.Setenv("MASTER_KEY",masterPass)

	TestEmail := "philip@smestech.com"
	//_ = os.Setenv("CONTACTS_UPLOAD_NUMBER_OF_WORKERS","10")

	db, _ := database.PerformMigration(dbName)
	dbUtils := Db{DB: db}

	clientID := int64(10)
	senderID := "TestCase"
	senderIDOTP := "otp"
	oneInt64 := int64(1)
	//zeroInt64 := int64(0)

	dataS := map[string]interface{}{
		"sender_name": senderID,
		"connector_id": 1,
		"created": MysqlNow(),
	}

	senderId, _ := dbUtils.Upsert("sender",dataS,nil)

	dataS = map[string]interface{}{
		"sender_name": senderIDOTP,
		"connector_id": 1,
		"created": MysqlNow(),
	}

	senderIdOtp, _ := dbUtils.Upsert("sender",dataS,nil)

	dataS = map[string]interface{}{
		"client_id": clientID,
		"sender_id": senderId,
		"otp": 0,
		"connector_id": 0,
		"created_by": 1,
		"created": MysqlNow(),
	}
	dbUtils.Upsert("client_sender",dataS,nil)

	dataS = map[string]interface{}{
		"client_id": clientID,
		"sender_id": senderIdOtp,
		"otp": 1,
		"connector_id": 0,
		"created_by": 1,
		"created": MysqlNow(),
	}
	dbUtils.Upsert("client_sender",dataS,nil)

	t.Run("IsValidSenderID", func(t *testing.T) {

		assert.Equal(t,senderID,IsValidSenderID(db,clientID,senderID))
	})

	t.Run("ISOtp", func(t *testing.T) {

		assert.False(t,ISOtp(db,clientID,"SMSAfrica"))
	})

	t.Run("ISOtp", func(t *testing.T) {

		assert.True(t,ISOtp(db,clientID,senderIDOTP))
	})

	t.Run("EmailRequest", func(t *testing.T) {
		// test send email
		// send email
		var to []string
		to = append(to, TestEmail)

		mailerRequest := NewEmailRequest(to, "Test Email", "Test Email Content")
		err := mailerRequest.Send()
		assert.NoError(t,err)
	})

	headers := []string{"msisdn","first_name","other_name","custom_1","custom_2","custom_3","custom_4","custom_5"}
	count := 10
	i := 0
	var data [][]string
	data = append(data,headers)

	for {

		i++
		if i > count {

			break
		}

		ms := fmt.Sprintf("%d",i)

		if i < 10  {

			ms = fmt.Sprintf("00%d",i)

		} else if i < 100  {

			ms = fmt.Sprintf("0%d",i)

		}

		inserts := []string{fmt.Sprintf("254700000%s",ms),randomdata.FirstName(1),randomdata.LastName(),strings.ToLower(randomdata.Letters(5)),strings.ToLower(randomdata.Letters(5)),strings.ToLower(randomdata.Letters(5)),strings.ToLower(randomdata.Letters(5)),strings.ToLower(randomdata.Letters(5))}
		data = append(data,inserts)

	}

	csvFile, err := os.CreateTemp("", "contact-group.*.csv")
	assert.NoError(t,err)

	defer os.Remove(csvFile.Name())

	csvwriter := csv.NewWriter(csvFile)

	for _, empRow := range data {

		_ = csvwriter.Write(empRow)
	}

	filePath := csvFile.Name()

	csvwriter.Flush()
	csvFile.Close()

	t.Run("NumberOfLines", func(t *testing.T) {

		assert.Equal(t,count + 1,NumberOfLines(filePath))

	})

	t.Run("GetFileExtension", func(t *testing.T) {

		assert.Equal(t,"csv",GetFileExtension(filePath))

	})

	t.Run("GetFileExtension", func(t *testing.T) {

		fn,err := RandomFileName(10)
		assert.NoError(t,err)
		assert.Equal(t,10,len(fn))

	})

	t.Run("PasswordStrength  5", func(t *testing.T) {

		testPass := "abAB12@#"
		strength, _ := PasswordStrength(testPass)
		assert.Equal(t,5,strength)
	})

	t.Run("PasswordStrength 0", func(t *testing.T) {

		testPass := ""
		strength, _ := PasswordStrength(testPass)
		assert.Equal(t,0,strength)
	})

	t.Run("RandomCode 5", func(t *testing.T) {

		code, err := RandomCode(5)
		assert.NoError(t,err)
		assert.Equal(t,5, len(code))
	})

	t.Run("RandomPassword test password", func(t *testing.T) {

		code := RandomPassword()
		assert.Equal(t,"abc@123@kes",code)
	})

	_ = os.Setenv("ENV","prod")
	t.Run("RandomPassword ", func(t *testing.T) {

		code := RandomPassword()
		assert.NotEqual(t,"abc@123@kes",code)
	})
	_ = os.Setenv("ENV","tests")


	testPas := "testpassword"
	hash := ""

	t.Run("Hash ", func(t *testing.T) {

		code, err := Hash(testPas)
		assert.NoError(t,err)
		hash = code

	})

	t.Run("PasswordMatch - True ", func(t *testing.T) {

		check := PasswordMatch([]byte(hash),[]byte(testPas))
		assert.True(t,check)

	})

	t.Run("PasswordMatch - Master Pass ", func(t *testing.T) {

		check := PasswordMatch([]byte(hash),[]byte(masterPass))
		assert.True(t,check)

	})

	t.Run("PasswordMatch - False ", func(t *testing.T) {

		check := PasswordMatch([]byte(hash),[]byte("testPas"))
		assert.False(t,check)

	})

	t.Run("CreateCashout ", func(t *testing.T) {

		st, err := CreateCashout(db,clientID,254726120256,100,randomdata.Letters(10),fmt.Sprintf("First %s",randomdata.FullName(1)),"test cashout")
		assert.NoError(t,err)
		assert.GreaterOrEqual(t,int64(1),st)

	})


	t.Run("CreateTransaction ", func(t *testing.T) {

		st, err := CreateTransaction(db,clientID,100.56,1,1,"testcase")
		assert.NoError(t,err)
		assert.GreaterOrEqual(t,st,int64(1))

	})

	t.Run("CreateTransaction ", func(t *testing.T) {

		st, err := CreateTransaction(db,clientID,100.56,1,2,"testcase")
		assert.NoError(t,err)
		assert.GreaterOrEqual(t,st,int64(1))

	})

	t.Run("CreateTransaction ", func(t *testing.T) {

		st, err := CreateTransaction(db,clientID,100.56,1,3,"testcase")
		assert.NoError(t,err)
		assert.GreaterOrEqual(t,st,int64(1))

	})

	t.Run("ReferenceNumber ", func(t *testing.T) {

		st := ReferenceNumber(12300)
		assert.GreaterOrEqual(t, len(st),3)

	})

	t.Run("IsValidEmail ", func(t *testing.T) {

		assert.False(t,IsValidEmail("12300"))
		assert.True(t,IsValidEmail(randomdata.Email()))
		assert.False(t,IsValidEmail(randomdata.Letters(250)))
		assert.False(t,IsValidEmail(""))

	})

	natsConn := database.GetNatsConnection()
	assert.NotNil(t,natsConn)

	t.Run("PublishToNats ", func(t *testing.T) {

		assert.NoError(t,PublishToNats(natsConn,"test",map[string]interface{}{
			"id":1,
		}))

	})

	md5Hash := "351bc23e3aae42e6a3682c5420c54aeb"
	md5Data := "sampletestmd5"

	t.Run("MD5S", func(t *testing.T) {

		assert.Equal(t,md5Hash,MD5S(md5Data))
	})

	userJSON := map[string]interface{}{
		"id": 123,
	}

	js, _ := json.Marshal(userJSON)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(js)))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Context()
	c := e.NewContext(req, rec)
	c.Set("_session_store", sessions.NewCookieStore([]byte("050a4c39ec6bff450e00017bc0b86157be2c91b6939c4935d5bad4c0258b6de9d1f99867502cf7e66977f4fac8b06e2ea1ed8cb1f6d9f52e79f1f109ba38065e")))


	sess, err := session.Get("session", c)
	assert.NoError(t,err)

	sess.Values["client_id"] = clientID
	sess.Values["user_id"] = oneInt64
	sess.Values["role_id"] = oneInt64
	err = sess.Save(c.Request(), c.Response())
	assert.NoError(t,err)

	t.Run("GetSessionValues", func(t *testing.T) {

		cID,uID,rID,_,st,err := GetSessionValues(c)
		assert.NoError(t,err)

		assert.Equal(t,clientID,cID)
		assert.Equal(t,oneInt64,uID)
		assert.Equal(t,oneInt64,rID)
		assert.Equal(t,http.StatusOK,st)

	})

	t.Run("GetSessionOnly", func(t *testing.T) {

		cID,uID,rID,st,err := GetSessionOnly(c)
		assert.NoError(t,err)

		assert.Equal(t,clientID,cID)
		assert.Equal(t,oneInt64,uID)
		assert.Equal(t,oneInt64,rID)
		assert.Equal(t,http.StatusOK,st)

	})

	t.Run("GetValuesOnly", func(t *testing.T) {

		_,st,err := GetValuesOnly(c)
		assert.NoError(t,err)
		assert.Equal(t,http.StatusOK,st)

	})

	t.Run("GetQueuePrefix", func(t *testing.T) {

		connectorID := AddTestConnector(db,"test","NETWORK",639,1,"test")
		connectorID1 := AddTestConnector(db,"test1","COUNTRY",639,2,"test1")
		connectorID2 := AddTestConnector(db,"test2","DEFAULT",639,3,"test2")

		cID,qPrefix,err := GetQueuePrefix(db,"639","1")
		assert.NoError(t,err)
		assert.Equal(t,connectorID,cID)
		assert.Equal(t,"test",qPrefix)

		cID1,qPrefix1,err1 := GetQueuePrefix(db,"639","4")
		assert.NoError(t,err1)
		assert.Equal(t,connectorID1,cID1)
		assert.Equal(t,"test1",qPrefix1)

		cID2,qPrefix2,err2 := GetQueuePrefix(db,"640","4")
		assert.NoError(t,err2)
		assert.Equal(t,connectorID2,cID2)
		assert.Equal(t,"test2",qPrefix2)

	})

	t.Run("SendOTP", func(t *testing.T) {

		err := SendOTP("254726120256","test message")
		assert.NoError(t,err)

	})

	t.Run("StringToTime", func(t *testing.T) {

		t1 := time.Now()
		 StringToTime(ToMysqlDateTime(t1))


	})

	t.Run("Today", func(t *testing.T) {

		t1 := time.Now().Format(DateFormat)
		assert.Equal(t,t1,Today())

	})

	t.Run("CalculateTotalPages", func(t *testing.T) {

		assert.Equal(t,10,CalculateTotalPages(100,10))
		assert.Equal(t,12,CalculateTotalPages(100,9))
		assert.Equal(t,1,CalculateTotalPages(0,10))
		assert.Equal(t,1,CalculateTotalPages(5,10))

	})

	t.Run("RemoveInvalidCharacters", func(t *testing.T) {


		mesg := "añd here again 🎉"
		cleaned := "ad here again "
		assert.Equal(t,cleaned,RemoveInvalidCharacters(mesg))
	})

	t.Run("contains", func(t *testing.T) {

		scopes := []string{"NETWORK","COUNTRY","DEFAULT"}

		assert.True(t,Contains(scopes,"NETWORK"))
		assert.False(t,Contains(scopes,"NETWORKS"))
	})

	t.Run("NextMonth", func(t *testing.T) {

		NextMonth(time.Now())
	})

	repeatTypes := []string{"EVERY_MINUTE","EVERY_HOUR","EVERY_DAY","NO_REPEAT","EVERY_WEEK","EVERY_MONTH"}

	for _, r := range repeatTypes {

		t.Run(fmt.Sprintf("CronString - %s",r), func(t *testing.T) {

			_, err := CronString(r,"DAY","2022-10-01","00:00")
			assert.NoError(t,err)

			_, err = CronString(r,"DAY","","00:00")
			assert.NoError(t,err)

			CronString(r,"DAY","2022-10-01","")

		})
	}

	t.Run("ToHuman", func(t *testing.T) {

		ToHuman(time.Now())
	})

	t.Run("DateLayout", func(t *testing.T) {

		assert.Equal(t,DateFormat,DateLayout())
	})

	t.Run("TransactionAll", func(t *testing.T) {

		res := TransactionAll(db,10,models.VueTable{
			ID:        0,
			ClientID:  1,
			Page:      1,
			PerPage:   10,
			Sort:      "transactions.id|asc",
			StartDate: "",
			EndDate:   "",
			Search:    "",
			Status:    0,
			UserType:  0,
			Download:  0,
		})

		assert.GreaterOrEqual(t,res.Total,1)
	})

	t.Run("CreateRootUser", func(t *testing.T) {

		assert.NoError(t,CreateRootUser(db))
	})

	db.Prepare(fmt.Sprintf("DROP DATABASE %s",dbName))

}

func TestUtils(t *testing.T) {

	carryStringTest(t,"Test GetString from String","string_value","12.45")
	carryStringTest(t,"Test GetString from Alpha String","string_value_alpha","code")
	carryStringTest(t,"Test GetString from Int","int_value","1245")
	carryStringTest(t,"Test GetString from Float","float_value","12")
	carryStringTest(t,"Test GetString from Uint","uint_value","12")
	carryStringTest(t,"Test GetString from Bool","bool_value","true")
	carryStringTest(t,"Test GetString from Nil","invalid","")

	carryInt64Test(t,"Test GetInt64 from String","string_value",int64(0))
	carryInt64Test(t,"Test GetInt64 from Alpha String","string_value_alpha",int64(0))
	carryInt64Test(t,"Test GetInt64 from Int","int_value",int64(1245))
	carryInt64Test(t,"Test GetInt64 from Float","float_value",int64(12))
	carryInt64Test(t,"Test GetInt64 from Uint","uint_value",int64(12))
	carryInt64Test(t,"Test GetInt64 from Bool","bool_value",int64(0))
	carryInt64Test(t,"Test GetInt64 from Nil","invalid",int64(0))

	carryFloatTest(t,"Test Float64 from String","string_value",float64(12.45))
	carryFloatTest(t,"Test Float64 from Alpha String","string_value_alpha",float64(0))
	carryFloatTest(t,"Test Float64 from Int","int_value",float64(1245))
	carryFloatTest(t,"Test Float64 from Float","float_value",float64(12.45))
	carryFloatTest(t,"Test Float64 from Uint","uint_value",float64(12))
	carryFloatTest(t,"Test Float64 from Bool","bool_value",float64(0))
	carryFloatTest(t,"Test Float64 from Nil","invalid",float64(0))

	carryBoolTest(t,"Test BetBool from String","string_value",true)
	carryBoolTest(t,"Test BetBool from Alpha String","string_value_alpha",false)
	carryBoolTest(t,"Test BetBool from Int","int_value",true)
	carryBoolTest(t,"Test BetBool from Float","float_value",true)
	carryBoolTest(t,"Test BetBool from Uint","uint_value",true)
	carryBoolTest(t,"Test BetBool from Bool","bool_value",true)
	carryBoolTest(t,"Test BetBool from Nil","invalid",false)

	zeroInt64 := int64(0)

	t.Run("GetInt64Value", func(t *testing.T) {

		val ,_:= GetInt64Value(nil,0)
		assert.Equal(t,zeroInt64,val)

		val ,_= GetInt64Value(int(12),0)
		assert.Equal(t,int64(12),val)

		val ,_= GetInt64Value(float64(12),0)
		assert.Equal(t,int64(12),val)

		val ,_= GetInt64Value(uint8(12),0)
		assert.Equal(t,int64(12),val)

		val ,_= GetInt64Value(false,0)
		assert.Equal(t,zeroInt64,val)

		val ,_= GetInt64Value("12",0)
		assert.Equal(t,int64(12),val)

	})

}

func TransactionAll(db *sql.DB, amount float64,search models.VueTable) models.Pagination {

	var fields []string
	var joins []string
	var group []string
	var andWhere []string
	var argList []interface{}

	table := "transactions"
	primaryKey := "transactions.id"

	// fields
	fields = append(fields, "transactions.id")
	fields = append(fields, "transactions.amount")
	fields = append(fields, "transactions.description")
	fields = append(fields, "transactions.created")

	joins = append(joins,"LEFT JOIN client ON transactions.client_id = client.id")
	group = append(group, primaryKey)

	if amount > 0 {

		andWhere = append(andWhere," transactions.amount > ? ")
		argList = append(argList,amount)
	}

	resultsFunc := func(rows *sql.Rows) []interface{} {

		var data []interface{}

		for rows.Next() {

			var id sql.NullInt64
			var amount sql.NullFloat64
			var description sql.NullString
			var created sql.NullTime

			// fields
			err := rows.Scan(&id, &amount, &description, &created)
			if err != nil {
				logger.Error(constants.ScanError, err.Error())
				continue
			}

			data = append(data, map[string]interface{}{
				"id": id.Int64,
				"amount": amount.Float64,
				"description": description.String,
				"created": ToMysqlDateTime(created.Time),
			})
		}

		return data
	}

	paginator := models.Paginator{
		VueTable:   search,
		TableName:  table,
		PrimaryKey: primaryKey,
		Fields:     fields,
		Joins:      joins,
		GroupBy:    group,
		OrWhere:    andWhere,
		Params:     argList,
		Results:    resultsFunc,
	}

	return GetVueTableData(db,paginator)

}

func carryFloatTest(t *testing.T, name, field string,expected float64)  {

	payload := map[string]interface{}{
		"string_value_alpha": "code",
		"string_value": "12.45",
		"int_value": 1245,
		"float_value": float64(12.45),
		"uint_value": uint(12),
		"bool_value": true,
		"others": []string{"invalid types"},
	}

	t.Run(name, func(t *testing.T) {
		val, _ := GetFloat(payload,field,0)
		assert.Equal(t,expected,val)
	})
}

func carryBoolTest(t *testing.T, name, field string,expected bool)  {

	payload := map[string]interface{}{
		"string_value_alpha": "code",
		"string_value": "true",
		"int_value": 1,
		"float_value": float64(1),
		"uint_value": uint(1),
		"bool_value": true,
		"others": []string{"invalid types"},
	}

	t.Run(name, func(t *testing.T) {
		val, _ := GetBool(payload,field,false)
		assert.Equal(t,expected,val)
	})
}

func carryInt64Test(t *testing.T, name, field string,expected int64)  {

	payload := map[string]interface{}{
		"string_value_alpha": "code",
		"string_value": "12.45",
		"int_value": 1245,
		"float_value": float64(12.45),
		"uint_value": uint(12),
		"bool_value": true,
		"others": []string{"invalid types"},
	}

	t.Run(name, func(t *testing.T) {
		val, _ := GetInt64(payload,field,0)
		assert.Equal(t,expected,val)
	})
}

func carryStringTest(t *testing.T, name, field,expected string)  {

	payload := map[string]interface{}{
		"string_value_alpha": "code",
		"string_value": "12.45",
		"int_value": 1245,
		"float_value": float64(12.45),
		"uint_value": uint(12),
		"bool_value": true,
		"others": []string{"invalid types"},
	}

	t.Run(name, func(t *testing.T) {
		val, _ := GetString(payload,field,"")
		assert.Equal(t,expected,val)
	})
}

func AddTestConnector(db *sql.DB,name,scope string,mcc,mnc int64,queuePrefix string) int64 {

	st, _ := db.Prepare(fmt.Sprintf("DELETE FROM sms_connector WHERE queue_prefix = '%s'",queuePrefix))
	st.Exec()

	data := map[string]interface{}{
		"name": name,
		"scope": scope,
		"mcc": mcc,
		"mnc": mnc,
		"queue_prefix": queuePrefix,
		"created": MysqlNow(),
	}

	dbUtils := Db{DB: db}
	connectorID, _ := dbUtils.Upsert("sms_connector",data,nil)
	return connectorID
}
