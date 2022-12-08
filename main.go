package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"ontunerestserver/app"
	"ontunerestserver/client"
	"ontunerestserver/data"
	"ontunerestserver/teemcache"
	"ontunerestserver/types"
	"runtime"
	"time"

	"github.com/urfave/negroni"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // 모든 CPU 사용
	
	var chend = make(chan string)
	var ch = make(chan string)

	clientRoom := client.NewClientRoom()
	go clientRoom.Run()	
	// go clientRoom.HostLastQueue.DataProducer(types.HOST_CODE)
	// go clientRoom.BasicQueue.DataProducer(types.BASIC_CODE)
	// go clientRoom.IOQueue.DataProducer(types.NET_CODE)
	go clientRoom.HostLastQueue.DataListen()
	go clientRoom.RunDatapump()
	
	//var mapch = make(chan map[string]teemcache.CacheData)

	//dbInfo := app.CreateDBInfo()
	//mux := app.MakeHandler("./test.db") //flag.Args 이런걸로 사용하자. 설정인자는 최대한 바깥으로 빼자
	//os.Getenv("DATABASE_URL")
	teemcache.MakeCacheServer()
	go teemcache.MakeCache(&ch, &chend)

	env := app.CreateDBInfo("server.env")
	mux := app.MakeHandler(env.GetDBConnString(), &ch, clientRoom)
	defer mux.Close() //finally 개념

	
	//conn, err := net.Dial("tcp", "192.168.0.140:8088")	 	
	conn, err := net.Dial("tcp", env.GetManagerConnString())	 	
	if err != nil {
		fmt.Println("Faield to Dial : ", err)
	}
	defer conn.Close()
	go writeHandler(conn) 
	go readHandler(conn, clientRoom.HostLastQueue) 

	ngri := negroni.Classic()
	ngri.UseHandler(mux)

	log.Println("Started App")
	err = http.ListenAndServe(":"+env.GetWPort(), ngri)
	if err != nil {
		panic(err)
	}
}

func writeHandler(conn net.Conn) {

	dataKey := types.DataKey{}
	dataKey.Code = types.DATAKEY_CODE
	dataKey.Key = 254	

	send, err := json.Marshal(dataKey)
	if err != nil {
		log.Printf("error: %s", err)
	}
	//send := "Hello2"
	for {
		_, err := conn.Write([]byte(send))
		if err != nil {
			fmt.Println("Failed to write data : ", err)
			break;
		}
		time.Sleep(6000 * time.Second)
	}
}

func readHandler(conn net.Conn, dQ *data.DataQueue) {
	header := make([]byte, 4)
	for {
		n, err := conn.Read(header)
		if err != nil {
			if io.EOF == err {
				return
			}
			log.Printf("Failed Connection: %v\n", err)
			return
		}
		if 0 < n {
			fmt.Printf("header %d\n", header)

			recv := make([]byte, header[0])
		n, err = conn.Read(recv)
			if err != nil {
				if io.EOF == err {
					return
				}
				log.Printf("Failed Connection: %v\n", err)
				return
			}
			if 0 < n {
				str_data := recv[:n]
				fmt.Printf("%s\n", str_data)
				// fmt.Println("---")				
				//dQ.DataChan <- message
				//fmt.Println(message)
				databuf := types.RealData{}
				err = json.Unmarshal(str_data, &databuf)
				if err != nil {
					log.Printf("error: %s", err)				
					fmt.Println(str_data)
					continue				
				}	
				// dQ.DataQueue.Push(message)		
					fmt.Println("code: ", databuf.Code, " data: ", databuf.Data)
					dQ.DataChan <- str_data					
			}
		}

		
	}
}

// func readHandler(conn net.Conn, dQ *data.DataQueue) {
// 	recv := make([]byte, 4096)
// 	databuf := types.RealData{}

// 	for {
// 		n, err := conn.Read(recv)
// 		if err != nil {
// 			fmt.Println("Failed to Read data : ", err)
// 			break
// 		}
// 		var message = recv[:n]
// 		//dQ.DataChan <- message
// 		//fmt.Println(message)
// 		err = json.Unmarshal(message, &databuf)
// 		if err != nil {
// 			log.Printf("error: %s", err)				
// 			fmt.Println(message)
// 			continue				
// 		}	

// 		// dQ.DataQueue.Push(message)		
// 		fmt.Println("code: ", databuf.Code, " data: ", databuf.Data)
// 		dQ.DataChan <- message		
// 	}
// }

