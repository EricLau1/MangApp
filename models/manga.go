package models

import (
	"errors"
	"fmt"
)

var (
	ErrMangaDescriptionTake = errors.New("Descrição ja existe")
)

type Manga struct {
	Id         int 	   `json:"id"`
	Descricao  string  `json:"descricao"`  
	Formato    string  `json:"formato"`
	Quantidade int     `json:"quantidade"`
	Volumes    string  `json:"volumes"`
	Status     string  `json:"status"`
	Valor      float64 `json:"valor"`
}

func NewManga(manga Manga) (bool, error) {

	isValid, erro := Unique( "mangas", "descricao", manga.Descricao )

	if erro != nil {

		return false, erro

	}

	if !isValid {
		return false, ErrMangaDescriptionTake
	}


	return registerManga(manga)

}

func registerManga(manga Manga) (bool, error) {

	con := Connect()

	sql := "insert into mangas (descricao, formato, quantidade, volumes, status, valor ) values (?,?,?,?,?,?)"

	stmt, erro := con.Prepare(sql)

	if erro != nil {

		return false, erro

	}

	_, erro = stmt.Exec(
			manga.Descricao,
			manga.Formato,
			manga.Quantidade,
			manga.Volumes,
			manga.Status,
			manga.Valor)

	if erro != nil {

		//panic(erro.Error())

		return false, erro

	}

	defer stmt.Close()
	defer con.Close()

	return true, nil

}

func GetMangas() ([]Manga, error) {

	con := Connect()

	sql := "select * from mangas"

	rs, erro := con.Query( sql )

	if erro != nil {

		return nil, erro

	}

	var mangas []Manga

	for rs.Next() {

		var manga Manga

		erro := rs.Scan( &manga.Id, 
						 &manga.Descricao, 
						 &manga.Formato, 
						 &manga.Quantidade, 
						 &manga.Volumes,
						 &manga.Status,
						 &manga.Valor )

		if erro != nil {

			return nil, erro

		}

		mangas = append(mangas, manga)

	}

	defer rs.Close()
	defer con.Close()

	return mangas, nil
}

func UpdateManga(manga Manga) (bool, error) {

	con := Connect()

	sql := "update mangas set descricao = ?, formato = ?, quantidade = ?, volumes = ?, status = ?, valor = ? where id = ?"

	stmt, erro := con.Prepare( sql )

	if erro != nil {

		return false, erro

	}

	_, erro = stmt.Exec( manga.Descricao, manga.Formato, manga.Quantidade, manga.Volumes, manga.Status, manga.Valor, manga.Id )

	if erro != nil {

		return false, erro

	}

	defer con.Close()
	defer stmt.Close()

	return true, nil

}

func GetMangaById(id int) (Manga, error) {

	con := Connect()

	sql := "select * from mangas where id = ?"

	rs, erro := con.Query( sql, id )

	if erro != nil {

		return Manga{}, nil

	}

	var manga Manga

	if rs.Next() {

		erro := rs.Scan(&manga.Id,
			&manga.Descricao,
			&manga.Formato,
			&manga.Quantidade,
			&manga.Volumes,
			&manga.Status,
			&manga.Valor)

		if erro != nil {

			return Manga{}, nil

		}
	}

	defer rs.Close()
	defer con.Close()

	return manga, nil
	
}

func DeleteManga(id int) (int64, error ) {

	con := Connect()

	sql := "delete from mangas where id = ?"

	stmt, erro := con.Prepare( sql )

	if erro != nil {

		return -1, erro

	}

	rs, erro := stmt.Exec( id )

	if erro != nil {

		return -1, erro

	}

	rows, erro := rs.RowsAffected()

	if erro != nil {

		return -1, erro

	}

	defer stmt.Close()
	defer con.Close()

	return rows, nil

}

func SearchManga(search string) ([]Manga, error) {

	con := Connect()

	sql := "select * from mangas where descricao like ?"

	rs, erro := con.Query( sql, search )
	
	if erro != nil {

		return nil, erro

	}

	var mangas []Manga

	for rs.Next() {

		var manga Manga

		erro := rs.Scan( &manga.Id, &manga.Descricao, &manga.Formato, &manga.Quantidade, &manga.Volumes, &manga.Status, &manga.Valor,  )

		if erro != nil {

			return nil, erro

		}

		mangas = append(mangas, manga)

	}

	defer rs.Close()
	defer con.Close()

	return mangas, nil

}

func (manga *Manga) ToPrintf() {

	fmt.Printf("{ Id:%d, Descricao:%s, Formato:%s, Quantidade:%d, Volumes:%s, Status:%s, Valor:%f }\n",
		manga.Id, manga.Descricao, manga.Formato, manga.Quantidade, manga.Volumes, manga.Status, manga.Valor)
	
}