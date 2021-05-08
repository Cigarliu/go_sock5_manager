package httpsocks

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gin-gonic/gin"

)
var MysqlDb *sql.DB
var MysqlDbErr error

const (
	USER_NAME = "cigarliu"
	PASS_WORD = "liuxuejia.123"
	HOST      = "gz-cynosdbmysql-grp-gtbfz5lr.sql.tencentcdb.com"
	PORT      = "29692"
	DATABASE  = "socks5"
	CHARSET   = "utf8"
)

type DBuser struct {
	id int
	user string
	pass string
	y int
	m int
	d int
	timeStamp int
	maxDevice int
}


func InitDB()(interface{}){
	UserPass = make(map[string]DBuser)
	//UserList = make(map[string]UserConn)
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)

	// 打开连接失败
	MysqlDb, MysqlDbErr = sql.Open("mysql", dbDSN)
	//defer MysqlDb.Close();
	if MysqlDbErr != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + MysqlDbErr.Error())
	}

	// 最大连接数
	MysqlDb.SetMaxOpenConns(100)
	// 闲置连接数
	MysqlDb.SetMaxIdleConns(20)
	// 最大连接周期
	MysqlDb.SetConnMaxLifetime(100*time.Second)
	if MysqlDbErr = MysqlDb.Ping(); nil != MysqlDbErr {
		fmt.Print("mysql connect error")
		return MysqlDbErr.Error()
	}
	fmt.Print("mysql connect ok\n")
	return nil
}

var UserPass map[string]DBuser /*创建集合 */

var UserIpList1  map[string]string /*创建集合 */
var UserIpList2  map[string]string /*创建集合 */
var UserIpList3  map[string]string /*创建集合 */

type IpList struct {
	Ip1 string
	Ip2 string
	Ip3 string
}

type UserConn struct {
	ConnNum int
	ConnList IpList
	Lasttimestamp int
}



func GetUserInfo(user string)(DBuser,error){
	var u DBuser
	u.user = user
	sqlStr := "select id,pass,y,m,d,timestamp,max_device from user_info where user=?"
    mapUser,ok :=UserPass[user]
	if (ok){
		//fmt.Println("使用map查询")
		return mapUser,nil
	}

	
	rowObj :=MysqlDb.QueryRow(sqlStr,user)
	err := rowObj.Scan(&u.id,&u.pass,&u.y,&u.m,&u.d,&u.timeStamp,&u.maxDevice)
	if err != nil {
		//fmt.Println(err)
		return u,err
	}
	UserPass[user] = u
	fmt.Println(u)
	return u, nil
}
func CheckUser(user,pass string)error {
	u , err := GetUserInfo(user)
	if err != nil{
		return err
	}
	if u.pass != pass {
		return errors.New("pass check fail")
	}
	return nil
}

func AddUser(c *gin.Context){
	clientIP := c.ClientIP()

	user := c.Query("user")
	pass := c.Query("pass")
	overtime := c.Query("over_time")
	fmt.Print("\nuser:",user)
	fmt.Print("\npass:",pass)
	fmt.Print("\nover_time:",overtime)
	fmt.Print("\n:","------------------")

	//MysqlDb.Prepare("insert socks5 (user,pass,create_time,over_time) value (?,?,?,?)")
	//sql_ex,_ := MysqlDb.Prepare("insert socks5 (user,pass,create_time,over_time) value (?,?,?,?)")
	////
	//res,_ := sql_ex.Exec(user,pass,time.Now().Unix(),overtime)
	//fmt.Print(res.LastInsertId())

	c.JSON(http.StatusOK,gin.H{
		"status":200,
		"msg":clientIP,
	})
}

func LoginHandler(c *gin.Context){
	clientIP := c.ClientIP()
	// request parameter
	//  user
	//  pass

	// response
	// status :  1  = ok   0 = auth fail
	user := c.Query("user")
	pass := c.Query("pass")

	err := CheckUser(user,pass)
	var status  int
	if err != nil {
		status = 0
	}else {
		status = 1
	}

	c.JSON(http.StatusOK,gin.H{
		"status":status,
		"msg":clientIP,
	})
}










func WebStart()  {
	r := gin.Default()
	r.GET("/login",LoginHandler)
	r.GET("/addUserCigar",AddUser)

	err :=r.Run(":8989")
	if(err !=nil) {
		fmt.Println("servevr run fail")
	}
}