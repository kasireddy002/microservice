package visitservice

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

//KINDNAME is the table name to store the values
const KINDNAME = "Visit"

//NAMESPACENAME is the Namespace
const NAMESPACENAME = "-kashyak-"

//VisitEntiry is struct to hold the vist details
type VisitEntiry struct {
	ID        int64  `datastore:"-"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	PassWord  string `json:"password"`
	Email     string `json:"email"`
}

var visit *VisitEntiry

//SuccessResponse store response
type SuccessResponse struct {
	//	visit   VisitEntiry `json:"entity"`
	ID      int64  `json:"Id"`
	Message string `json:"message"`
}

func init() {

	http.HandleFunc("/api/getallvisits/", restHandler)
	http.HandleFunc("/api/postavisit/", restHandler)
	http.HandleFunc("/api/deleteavisit/", restHandler)
}

func restHandler(w http.ResponseWriter, r *http.Request) {

	var v VisitEntiry
	_ = json.NewDecoder(r.Body).Decode(&v)
	json.NewEncoder(w).Encode(v)

	visit = &VisitEntiry{
		ID:        v.ID,
		FirstName: v.FirstName,
		LastName:  v.LastName,
		UserName:  v.UserName,
		PassWord:  v.PassWord,
		Email:     v.Email,
	}

	switch r.Method {
	case "GET":
		getallvisitshandler(w, r)
		return
	case "POST":
		putavisthandler(w, r)
		return
	case "DELETE":
		deletevisithandler(w, r)
		return
	default:
		//respondErr(w, r, http.StatusNotFound, "is not supported HTTP methods")
	}
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func getallvisitshandler(w http.ResponseWriter, r *http.Request) {
	var visitslist []VisitEntiry
	var ctx context.Context
	ctx = appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}
	x := strings.Split(r.URL.Path, "/")[3]
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	log.Printf("%#v Getting values url - y ", y)

	if y > 0 {
		key := datastore.NewKey(ctx, KINDNAME, "", int64(y), nil)
		err = datastore.Get(ctx, key, visit)

		var models VisitEntiry

		models.ID = key.IntID()
		models.FirstName = visit.FirstName
		models.LastName = visit.LastName
		models.UserName = visit.UserName
		models.PassWord = visit.PassWord
		models.Email = visit.Email

		json.NewEncoder(w).Encode(models)

	} else {

		q := datastore.NewQuery(KINDNAME)

		keys, _ := q.GetAll(ctx, &visitslist)

		models := make([]VisitEntiry, len(visitslist))

		for i := 0; i < len(visitslist); i++ {
			models[i].ID = keys[i].IntID()
			models[i].FirstName = visitslist[i].FirstName
			models[i].LastName = visitslist[i].LastName
			models[i].UserName = visitslist[i].UserName
			models[i].PassWord = visitslist[i].PassWord
			models[i].Email = visitslist[i].Email
		}
		json.NewEncoder(w).Encode(models)
	}

}
func putavisthandler(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context
	keys := make([]*datastore.Key, 1)
	ctx = appengine.NewContext(r)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}
	x := strings.Split(r.URL.Path, "/")[3]
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	log.Printf("%#v Getting values url - y ", y)
	if err != nil {
		keys[0] = datastore.NewIncompleteKey(ctx, KINDNAME, nil)
	} else {
		keys[0] = datastore.NewKey(ctx, KINDNAME, "", int64(y), nil)
	}
	k, err := datastore.Put(ctx, keys[0], visit)

	var models VisitEntiry
	models.ID = k.IntID()
	models.FirstName = visit.FirstName
	models.LastName = visit.LastName
	models.UserName = visit.UserName
	models.PassWord = visit.PassWord
	models.Email = visit.Email
	json.NewEncoder(w).Encode(models)

}

func deletevisithandler(w http.ResponseWriter, r *http.Request) {
	var visitslist []VisitEntiry
	var ctx context.Context
	ctx = appengine.NewContext(r)
	keys := make([]*datastore.Key, 1)
	ctx, err := appengine.Namespace(ctx, NAMESPACENAME)
	if err != nil {
		return
	}
	x := strings.Split(r.URL.Path, "/")[3]
	log.Printf("%#v Getting values url - x ", x)
	y, err := strconv.Atoi(x)
	log.Printf("%#v Getting values url - y ", y)
	if y > 0 {

		keys[0] = datastore.NewKey(ctx, KINDNAME, "", int64(y), nil)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
	} else {
		q := datastore.NewQuery(KINDNAME)
		keys, _ = q.GetAll(ctx, &visitslist)
		option := &datastore.TransactionOptions{XG: true}
		err = datastore.RunInTransaction(ctx, func(c context.Context) error {
			return datastore.DeleteMulti(c, keys)
		}, option)
	}

}
