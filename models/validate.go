package models 

import "fmt"

func Unique(table, key, val string) (bool, error) {

	con := Connect()
		
	sql := fmt.Sprintf("select COUNT(*) from %s where %s = ?", table, key)

	var rows int

	erro := con.QueryRow( sql, val ).Scan( &rows )

	if erro != nil {

		return false, erro

	}

	if rows != 0 {

		return false, nil
	
	}

	defer con.Close()

	return true, nil

}

func UniqueUpdate(table, key, val , ignoreKey string, ignoreValue int) (bool, error) {

	con := Connect()

	var rows int

	sql := fmt.Sprintf("select count(*) from %s where %s = ? and %s != ?", table, key, ignoreKey)

	erro := con.QueryRow( sql, val, ignoreValue ).Scan(&rows)

	if erro != nil {

		return false, erro

	}

	if rows > 0 {

		return false, nil

	}

	defer con.Close()

	return true, nil
}