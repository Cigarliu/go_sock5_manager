package httpsocks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"

)




func Get(url string)(string,interface{}){
	res, err :=http.Get(url)
	if err != nil {
		//fmt.Println("get fail ")
		return "",errors.New("fail")
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "",errors.New("fail")
	}
	//fmt.Print(string(robots))
	return string(robots),nil
}

type HttpResCheck struct {
	Status int
	Msg string
}
func CheckUser(user,pass string) (interface{}){
	var url ="http://flv.comeboy.cn:8989/login"
	var urlUserPass = url + "?" +"user=" + user +"&pass=" +pass
	//fmt.Println(urlUserPass)
	res,err := Get(urlUserPass)
	if err !=nil {
		return err
	}
    status :=HttpResCheck{}
	//fmt.Println(res)

	err = json.Unmarshal([]byte(res),&status)
	if err != nil {
		fmt.Print(err)
		return err
	}
	//fmt.Println(status.Status)
	//fmt.Println(status.Msg)
	if status.Status != 1{
		//fmt.Print("pass error")
		return errors.New("pass error")
	}
	return nil

}
