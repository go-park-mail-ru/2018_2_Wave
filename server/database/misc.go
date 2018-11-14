package database

import (
	"log"
	"strconv"
)

func (model *DatabaseModel) present(tableName string, colName string, target string) (fl bool, err error) {
	var exists string
	row := model.Database.QueryRowx("SELECT EXISTS (SELECT true FROM " + tableName + " WHERE " + colName + "='" + target + "');")
	err = row.Scan(&exists)

	if err != nil {
		//log.Fatal(err)
		panic(err)
		return false, err
	}

	fl, err = strconv.ParseBool(exists)
	if err != nil {
		//log.Fatal(err)
		panic(err)
		return false, err
	}

	log.Println(fl)
	return fl, nil
}

func validateCredentials(target string) bool {
	/*
		// http://regexlib.com/REDetails.aspx?regexp_id=2298
		reg, _ := regexp.Compile("^([a-zA-Z])[a-zA-Z_-]*[\\w_-]*[\\S]$|^([a-zA-Z])[0-9_-]*[\\S]$|^[a-zA-Z]*[\\S]$")

		if reg.MatchString(target) {
			return true
		}
		log.Println("bad username or/and password")

		return false
	*/
	return true
}
