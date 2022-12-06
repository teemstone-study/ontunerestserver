package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ontunerestserver/client"
	"ontunerestserver/data"
	"ontunerestserver/teemcache"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

var rd *render.Render = render.New()

type AppHandler struct {
	//handler http.Handler
	http.Handler
	db     data.DBHandler
	chdata *chan string
	clientRoom *client.ClientRoom
	id uint64
}

type Success struct {
	Success bool `json:"success"`
}

func cache_map_search(aType string, aIds string) (string, bool) {
	// Cache Map Search & Data Return
	cachemap := teemcache.GetCachemap()

	if cachemap[aType].Id == aIds {
		// Data 존재 시, 캐시 데이터 전송
		if cachemap[aType].Data != "" {
			return cachemap[aType].Data, true
		}
	}
	return "", false
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusTemporaryRedirect) //인덱스 경로로 들어와도 todo.html 리다이렉션 해라
}

func (a *AppHandler) getPagesHandler(w http.ResponseWriter, r *http.Request) {
	list := a.db.GetPages()
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) GetPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index, _ := strconv.Atoi(vars["id"])
	//ok := model.RemoveTodo(id)
	//log.Println("index: " + strconv.Itoa(index))
	page := a.db.GetPage(index)

	if page != nil {
		rd.JSON(w, http.StatusOK, page)
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)
	}
}

func (a *AppHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := client.Upgrader.Upgrade(w, r, nil)		
	if err != nil {
		log.Println(err)
		return
	}
	a.id ++
	client := client.NewClient(a.id, a.clientRoom, conn) //Client conn 을 가지는 Clinet 생성 -ReadPump, WritePump..set
	client.Room.ChanEnter <- client //Clinet 의 주소값을 넘긴다. 
}

// func (a *AppHandler) addPageHandler(w http.ResponseWriter, r *http.Request) {
// 	page := new(data.Page)
// 	err := json.NewDecoder(r.Body).Decode(page)
// 	if err != nil {
// 		//rd.Text(w, http.StatusBadRequest, err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err)
// 		return
// 	}
// 	//log.Println("test  >>> " + strconv.Itoa( page.Index) + " " + page.Contents)
// 	a.db.AddPage(page)
// 	rd.JSON(w, http.StatusCreated, page)
// }

// func (a *AppHandler) UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
// 	page := new(data.Page)
// 	err := json.NewDecoder(r.Body).Decode(page)
// 	if err != nil {
// 		//rd.Text(w, http.StatusBadRequest, err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprint(w, err)
// 		return
// 	}

// 	ok := a.db.UpdatePage(page)

// 	if ok {
// 		rd.JSON(w, http.StatusOK, Success{true})
// 	} else {
// 		rd.JSON(w, http.StatusOK, Success{false})
// 	}
// }

func (a *AppHandler) addPageHandler(w http.ResponseWriter, r *http.Request) {
	pages := []*data.Page{}
	err := json.NewDecoder(r.Body).Decode(&pages)
	if err != nil {
		//rd.Text(w, http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	//log.Println("test  >>> " + strconv.Itoa( page.Index) + " " + page.Contents)
	var errNum int = 0
	for _, p := range pages {
		if !a.db.AddPage(p) {
			errNum++
		}
	}

	if errNum == 0 {
		rd.JSON(w, http.StatusCreated, pages)
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)
	}
	//rd.JSON(w, http.StatusCreated, pages)
}

func (a *AppHandler) UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	pages := []*data.Page{}
	err := json.NewDecoder(r.Body).Decode(&pages)
	if err != nil {
		//rd.Text(w, http.StatusBadRequest, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	var errNum int = 0
	for _, p := range pages {
		ok := a.db.UpdatePage(p)
		if !ok {
			errNum++
		}
	}

	if errNum == 0 {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (a *AppHandler) DeletePageHandler(w http.ResponseWriter, r *http.Request) {
	ok := a.db.DeletePage()
	if ok {
		rd.JSON(w, http.StatusOK, nil)
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)
	}
}

func (a *AppHandler) GetRealTimePerfHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// ids  , _ := strconv.Atoi(vars["ids"])
	ids := vars["ids"]
	stime, _ := strconv.Atoi(vars["stime"])
	etime, _ := strconv.Atoi(vars["etime"])
	//ok := model.RemoveTodo(id)
	//log.Println("index: " + strconv.Itoa(index))
	real := a.db.GetRealTimePerf(ids, stime, etime)

	if real != nil {
		rd.JSON(w, http.StatusOK, real)
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)
	}
}

func (a *AppHandler) GetRealTimeCpuHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// ids  , _ := strconv.Atoi(vars["ids"])
	ids := vars["ids"]
	stime, _ := strconv.Atoi(vars["stime"])
	etime, _ := strconv.Atoi(vars["etime"])
	//ok := model.RemoveTodo(id)
	//log.Println("index: " + strconv.Itoa(index))

	tmp_data, tmp_result := cache_map_search("realtimecpu", ids)

	if tmp_result {
		var data map[string]interface{}         // JSON 문서의 데이터를 저장할 공간을 맵으로 선언
		json.Unmarshal([]byte(tmp_data), &data) // Data를 바이트 슬라이스로 변환하여 넣고, data의 포인터를 넣어줌

		rd.JSON(w, http.StatusOK, data)
		// fmt.Println("1111111")
		// fmt.Println(data)

		return
	}

	// Data가 존재하지 않으므로, DB에 조회해야 함.
	real := a.db.GetRealTimeCpu(ids, stime, etime)

	if real != nil {
		rd.JSON(w, http.StatusOK, real)

		// fmt.Println("222222")
		// fmt.Println(real)

		// ch에 최신 데이터 저장시킴.
		data := make(map[string]string)
		data["datatype"] = "realtimecpu"
		data["id"] = ids
		out, _ := json.Marshal(real)
		data["data"] = string(out)
		*a.chdata <- data["datatype"] + "-" + data["id"] + "-" + data["data"]
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)

		// 캐시 데이터 삭제
		delete(teemcache.GetCachemap(), "realtimecpu")
	}
}

func (a *AppHandler) GetRealTimeNetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// ids  , _ := strconv.Atoi(vars["ids"])
	ids := vars["ids"]
	stime, _ := strconv.Atoi(vars["stime"])
	etime, _ := strconv.Atoi(vars["etime"])
	//ok := model.RemoveTodo(id)
	//log.Println("index: " + strconv.Itoa(index))
	real := a.db.GetRealTimeNet(ids, stime, etime)

	if real != nil {
		rd.JSON(w, http.StatusOK, real)
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)
	}
}

func (a *AppHandler) GetRealTimeDiskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// ids  , _ := strconv.Atoi(vars["ids"])
	ids := vars["ids"]
	stime, _ := strconv.Atoi(vars["stime"])
	etime, _ := strconv.Atoi(vars["etime"])
	//ok := model.RemoveTodo(id)
	//log.Println("index: " + strconv.Itoa(index))
	real := a.db.GetRealTimeDisk(ids, stime, etime)

	if real != nil {
		rd.JSON(w, http.StatusOK, real)
	} else {
		rd.JSON(w, http.StatusBadRequest, nil)
	}
}

func (a *AppHandler) Close() {
	a.db.Close()
}

// func MakeHandler() http.Handler {
func MakeHandler(dbConn string, ch *chan string, cRoom *client.ClientRoom) *AppHandler {

	mux := mux.NewRouter()
	//mux.Use(corsMiddleware)

	c := cors.New(cors.Options{
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT", "OPTIONS"},
	})

	a := &AppHandler{
		Handler: c.Handler(mux),
		db:      data.NewDBHandler(dbConn),
		chdata:  ch,
		id: 0,
		clientRoom: cRoom,
	}

	mux.HandleFunc("/ws", a.wsHandler)
	mux.HandleFunc("/pages", a.getPagesHandler).Methods("GET")
	mux.HandleFunc("/pages", a.addPageHandler).Methods("POST")
	mux.HandleFunc("/pages", a.UpdatePageHandler).Methods("PUT")
	mux.HandleFunc("/pages", a.DeletePageHandler).Methods("DELETE")
	mux.HandleFunc("/pages/{id:[0-9]+}", a.GetPageHandler).Methods("GET")
	mux.HandleFunc("/", a.indexHandler)

	// Handler RealTimePerf
	// mux.HandleFunc("/realtimeperf", a.getAllRealtimePerfHandler).Methods("GET")
	mux.HandleFunc("/realtimeperf/{ids}/{stime:[0-9]+},{etime:[0-9]+}", a.GetRealTimePerfHandler).Methods("GET")

	// Handler RealTimeCpu
	mux.HandleFunc("/realtimecpu/{ids}/{stime:[0-9]+},{etime:[0-9]+}", a.GetRealTimeCpuHandler).Methods("GET")

	// Handler RealTimeNet
	mux.HandleFunc("/realtimenet/{ids}/{stime:[0-9]+},{etime:[0-9]+}", a.GetRealTimeNetHandler).Methods("GET")

	// Handler RealTimeDisk
	mux.HandleFunc("/realtimedisk/{ids}/{stime:[0-9]+},{etime:[0-9]+}", a.GetRealTimeDiskHandler).Methods("GET")

	return a
}

// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Header().Set("Access-Control-Allow-Origin", "*")
// 			w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
// 			w.Header().Add("Access-Control-Allow-Credentials", "true")
// 			w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		    w.Header().Add("Access-Control-Expose-Headers", "Origin,  X-Auth-Token, Authorization, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
// 			w.Header().Set("content-type", "application/json;charset=UTF-8")

// 			if r.Method == "OPTIONS" {
// 					w.WriteHeader(http.StatusNoContent)
// 					return
// 			}
// 			next.ServeHTTP(w, r)
// 	})
// }
