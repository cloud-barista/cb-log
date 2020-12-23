package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
)

var cblogger *logrus.Logger

type DBCONN struct {
	name string
}

func init() {
	// cblog is a global variable.
	filePath := filepath.Join("..", "conf", "log_conf.yaml")
	cblogger = cblog.GetLogger("CB-SPIDER", filePath)
}

func main() {

	for {
		cblogger.Info("start.........")

		err := createUser3("newUser")
		cblogger.Debug("msg for debugging msg!!")
		if err != nil {
			cblogger.Error(err)
		}

		err = createUser4("newUser")
		if err != nil {
			cblogger.Error(err)
		}

		cblogger.Info("end.........")

		time.Sleep(time.Second * 2)
		fmt.Print("\n")
	}
}

func createUser3(newUser string) error {
	cblogger.Info("start creating user.")

	var db *DBCONN
	db = new(DBCONN)
	if db == nil {
		cblogger.Error("DBMS Session is closed!!")
		//panic(fmt.Errorf("panic by ..... !"))
	}

	isExist, err := checkUser(newUser)
	cblogger.Debug("msg for debugging msg!!")
	if isExist {
		return err
	}

	cblogger.Info("finish creating user.")
	return nil
}

func createUser4(newUser string) error {
	cblogger.Info("start creating user.")

	var db *DBCONN
	//      db := new(DBCONN)
	if db == nil {
		cblogger.Error("DBMS Session is closed!!")
		//panic(fmt.Errorf("panic by ..... !"))
	}

	isExist, err := checkUser(newUser)
	if isExist {
		return err
	}

	cblogger.Info("finish creating user.")
	return nil
}

func checkUser(user string) (bool, error) {
	return false, fmt.Errorf("%s: already existed User!!", user)
}
