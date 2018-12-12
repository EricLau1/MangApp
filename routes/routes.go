package routes;

import (
	"fmt"
	"strconv"
	"strings"
	"html"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"../utils"
	"../models"
	"../sessions"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", indexGetHandler).Methods("GET")

	router.HandleFunc("/manga", mangaCadastroGetHandler).Methods("GET")

	router.HandleFunc("/manga", mangaCadastroPostHandler).Methods("POST")

	router.HandleFunc("/manga/{id}", mangaGetHandler).Methods("GET")

	router.HandleFunc("/manga/edit/{id}", mangaEditGetHandler).Methods("GET")

	router.HandleFunc("/manga/edit/{id}", mangaEditPostHandler).Methods("POST")

	router.HandleFunc("/manga/edit-validate", mangaEditAjaxValidateHandler).Methods("POST")

	router.HandleFunc("/delete/manga", mangaDeleteGetHandler).Methods("GET")

	router.HandleFunc("/search", mangaSearchGetHandler).Methods("GET")

	fileServer := http.FileServer(http.Dir("./assets"))

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	return router
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {

	mangas, erro := models.GetMangas()

	if erro != nil {

		utils.InternalServerError(w)
		return
		
	}

	var totalValor float64 = 0.0
	var totalMangas int = 0

	for _, manga := range mangas {

		totalValor += manga.Valor * float64(manga.Quantidade)
		totalMangas += manga.Quantidade
	}

	session, _ := sessions.Store.Get(r, "session")
	untypedMessage := session.Values["message"]

	msg, ok := untypedMessage.(string) 

	if ok {
		delete(session.Values, "message")
		session.Save(r, w)
	} 

	utils.ExecuteTemplate(w, "index.html", struct{
			Title       string
			Mangas      []models.Manga
			CountMangas int
			ValorMangas float64
			Message     string
		}{
			Title: "Mangás",
			Mangas: mangas,
			CountMangas: totalMangas,
			ValorMangas: totalValor,
			Message: msg,	
		})

}

func mangaCadastroGetHandler(w http.ResponseWriter, r *http.Request) {

	utils.ExecuteTemplate(w, "manga_cadastro.html", nil)

}

func mangaCadastroPostHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	quantidade, _  := strconv.Atoi( r.PostForm.Get("quantidade") )
	valor, _       := strconv.ParseFloat( r.PostForm.Get("valor") , 64)

	var manga models.Manga

	manga.Id = 0;
	manga.Descricao  = strings.TrimSpace( r.PostForm.Get("descricao") )
	manga.Formato    = r.PostForm.Get("formato")
	manga.Quantidade = quantidade 
	manga.Volumes    = r.PostForm.Get("volumes")
	manga.Status     = r.PostForm.Get("status")
	manga.Valor      = valor    

	//fmt.Println("POST VALUES: ", manga)

	_, erro := models.NewManga(manga)

	if erro != nil {

		switch(erro) {

		case models.ErrMangaDescriptionTake:
			utils.ExecuteTemplate(w, "manga_cadastro.html", "Descrição já existe!")
			return
		default :
			utils.InternalServerError(w)
			return
		}

	}

	http.Redirect(w, r, "/", 302)

}

func mangaGetHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//fmt.Println(vars)

	id, _ := strconv.Atoi( vars["id"] )

	manga, erro := models.GetMangaById(id)

	if erro != nil {

		utils.InternalServerError(w)
		return

	}

	 // manga.ToPrintf()

	var display bool = true
	var total float64 = ( manga.Valor * float64( manga.Quantidade ) )

	if manga.Id == 0 {
		
		display = false
	}

	//utils.ToJson(w, manga)

	// verificando se existe uma session de mensagem
	session, _ := sessions.Store.Get( r, "session")
	untypedMessage := session.Values["message"]

	msg, _ := untypedMessage.(string)

	// excluindo a mensagem da sessão
	delete(session.Values, "message")

	// salvando estado da sessão
	session.Save(r, w)

	//fmt.Println(msg, ok)

	utils.ExecuteTemplate(w, "manga_detalhes.html", struct{
			Manga          models.Manga
			DisplayDetails bool
			Total          float64
			Message        string
		}{
			Manga:          manga,
			DisplayDetails: display,
			Total:          total,
			Message:        msg,
		})
}

func mangaEditGetHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, _ := strconv.Atoi( vars["id"] ) // convertendo string para int

	manga, erro := models.GetMangaById(id)

	//fmt.Println(manga)

	if erro != nil {

		utils.InternalServerError(w)
		return

	}

	var display bool = false
	var options string = ""

	if manga.Id != 0 {

		display = true

		formatos := []string{ "tanko", "meio-tanko", "livro" }

		//fmt.Println(formatos)

		for _, formato := range formatos {

			if formato == manga.Formato {

				options += fmt.Sprintf("<option value='%s' selected> %s </option>", formato, formato)

			} else {

				options += fmt.Sprintf("<option value='%s'> %s </option>", formato, formato)	

			}
		} // end for
		//fmt.Println(options)
	} // end if

	utils.ExecuteTemplate(w, "manga_edit.html", struct{
			Manga        models.Manga
			Formatos     template.HTML // este dado irá se comportar como um componente HTML
			AlertMessage string
			DisplayForm  bool
		}{
			Manga: manga,
			// UnescapeString não formata os caracteres especiais do HTML como Strings
			// template.HTML formata a string como uma TAG HTML
			Formatos: template.HTML( html.UnescapeString(options) ),
			AlertMessage: "",
			DisplayForm: display,
		})
}



func mangaEditPostHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id, _          := strconv.Atoi( r.PostForm.Get("id") )
	quantidade, _  := strconv.Atoi( r.PostForm.Get("quantidade") )
	valor, _       := strconv.ParseFloat( r.PostForm.Get("valor") , 64) // convertendo para Float64

	var manga models.Manga

	manga.Id         = id;
	manga.Descricao  = strings.TrimSpace( r.PostForm.Get("descricao") ) // removendo espaços em branco do inicio e do fim de strings
	manga.Formato    = r.PostForm.Get("formato")
	manga.Quantidade = quantidade 
	manga.Volumes    = r.PostForm.Get("volumes")
	manga.Status     = r.PostForm.Get("status")
	manga.Valor      = valor    

	// fmt.Println("POST VALUES: ", manga)

	ok, erro := models.UpdateManga( manga )

	if erro != nil {

		utils.InternalServerError(w)
		return

	}

	var message string = ""

	if ok {

		message = "Atualizado com sucesso!"

	} else {

		message = "Nenhuma informação foi atualizada."

	}

	session, _ := sessions.Store.Get(r, "session")

	session.Values["message"] = message

	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/manga/%d", manga.Id ) , 302)
}

func mangaEditAjaxValidateHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id, _     := strconv.Atoi( r.PostForm.Get("id") )
	descricao := r.PostForm.Get("descricao")

	//fmt.Printf("Ajax Request Values: %d, %s\n", id, descricao)

	ok, erro := models.UniqueUpdate("mangas", "descricao", descricao, "id", id)

	if erro != nil {

		utils.ToJson(w, struct {
			IsValid      bool
			AlertMessage string
		}{
			IsValid: ok,
			AlertMessage: "Ocorreu um erro interno.",	
		})
		return

	}

	utils.ToJson(w, struct {
		IsValid      bool
		AlertMessage string
	}{
		IsValid: ok,
		AlertMessage: "OK",	
	})

}

func mangaDeleteGetHandler(w http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()

	//fmt.Println("KEYS:")
	//fmt.Println(keys)

	// pegando valores por GET enviados pela URL
	id, _  := strconv.Atoi( keys.Get("id") )  //Get return empty string if key not found
	confirm := keys.Get("confirm")
	//fmt.Printf("ID: %d", id)

	var message string = "Nenhuma informação foi deletada"

	session, _ := sessions.Store.Get(r, "session")

	if confirm == "true" {

		//fmt.Printf("Manga com id:%d será deletado.\n", id)

		rows, erro := models.DeleteManga(id)

		if erro != nil {

			utils.InternalServerError(w)
			return

		}

		if rows > 0 {

			message = fmt.Sprintf("Mangá #%d foi excluido. %d linhas foram deletadas.", id, rows)

			session.Values["message"] = message
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)
			return

		}

	}

	session.Values["message"] = message
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/manga/%d", id), 302)

}

func mangaSearchGetHandler(w http.ResponseWriter, r *http.Request) {

	keys := r.URL.Query()

	search := keys.Get("search")

	search = strings.TrimSpace(search)

	param :=  fmt.Sprintf("%%%s%%", search)

	//fmt.Println(param)

	results, erro := models.SearchManga(param)

	if erro != nil {

		utils.InternalServerError(w)
		return

	}

	var count int = len(results)

	utils.ExecuteTemplate(w, "search_results.html", struct{
			Search  string
			Results []models.Manga
			Count   int
		}{
			Results: results,
			Count: count,
			Search: search,
		})

	//utils.ToJson(w, results)

}