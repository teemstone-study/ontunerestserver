package main

import (
	"bytes"
	"encoding/binary"
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

	go TryConnectedTCP(env, clientRoom) //재연결을 위해 고루틴으로 변경
	// conn, err := net.Dial("tcp", env.GetManagerConnString())
	// if err != nil {
	// 	fmt.Println("Faield to Dial : ", err)
	// }
	// defer conn.Close()
	// go writeHandler(conn)
	// go readHandler(conn, clientRoom.HostLastQueue)

	ngri := negroni.Classic()
	ngri.UseHandler(mux)

	log.Println("Started App")
	err := http.ListenAndServe(":"+env.GetWPort(), ngri)
	if err != nil {
		panic(err)
	}
}

func TryConnectedTCP(env *app.DBInfo, clientRoom *client.ClientRoom) net.Conn {
	var newconn net.Conn
	var err error
	var innerErrorChecker = make(chan bool, 1)
	var changeValue = make(chan uint32, 1)
	for {
		newconn, err = net.Dial("tcp", env.GetManagerConnString())
		if err != nil {
			fmt.Println("Faield to Dial : ", err)
			if newconn != nil {
				newconn.Close()
			}
			time.Sleep(3 * time.Second)
			fmt.Println("retry")
			continue
		} else {
			changeValue <- uint32(254)
			go ReadWriteHandler(newconn, clientRoom.HostLastQueue, changeValue, innerErrorChecker)
		}
		time.Sleep(1 * time.Second)
		<-innerErrorChecker //여기서 ReadWriteHandler 내부에 오류가 발생하면 연결 끊고 새로 진행되도록 처리
		err = newconn.Close()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("tcp Error Restart")

	}
}

// func writeHandler(conn net.Conn, errorChan chan bool) {

// 	dataKey := types.DataKey{}
// 	dataKey.Code = types.DATAKEY_CODE
// 	dataKey.Key = 254

// 	send, err := json.Marshal(dataKey)
// 	if err != nil {
// 		log.Printf("error: %s", err)
// 	}
// 	for {
// 		_, err := conn.Write([]byte(send))
// 		if err != nil {
// 			fmt.Println("Failed to write data : ", err)
// 			break
// 		} else {
// 			fmt.Print("TCP Send data : ")
// 			fmt.Println(dataKey)
// 		}
// 		time.Sleep(6 * time.Second)
// 	}
// 	errorChan <- true
// 	fmt.Println("shutdown TCP Handller")
// }

func readBufferData(messageLength uint32, conn net.Conn) []byte {
	Buffer := new(bytes.Buffer)
	var recv []byte = make([]byte, 4096)
	for messageLength > 0 {

		n, err := conn.Read(recv)
		if nil != err {
			fmt.Println("err : " + err.Error())
		}

		if 0 < n {
			data := recv[:n]
			Buffer.Write(data)
			messageLength -= uint32(n)
		}
	}
	var response []byte
	response = append(response, 0xFF)
	_, err := conn.Write(response)
	if err != nil {
		fmt.Println("errorTCPSEND : " + err.Error())
	}

	strTest := Buffer.Bytes()
	// err3 := json.Unmarshal([]byte(strTest), &realData)
	// if err3 != nil {
	// 	fmt.Println(err3)
	// }
	return strTest
}

// func readHandler(conn net.Conn, dQ *data.DataQueue, errorChan chan bool) {
// 	var header = make([]byte, 4)
// 	for {
// 		n, err := conn.Read(header)
// 		if err != nil {
// 			if io.EOF == err {
// 				break
// 			}
// 			log.Printf("Failed Connection: %v\n", err)
// 			break
// 		}
// 		if 0 < n {
// 			messageLength := binary.LittleEndian.Uint32(header)
// 			realTimeData := readBufferData(messageLength, conn)
// 			dQ.DataChan <- realTimeData
// 		}
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	errorChan <- true
// }

func ReadWriteHandler(conn net.Conn, dQ *data.DataQueue, changeValue chan uint32, errorChan chan bool) {
	var header = make([]byte, 4)
	forCheck := false
	for !forCheck {
		select {
		case bindData := <-changeValue:
			dataKey := types.DataKey{}
			dataKey.Code = types.DATAKEY_CODE
			dataKey.Key = types.Bitmask(bindData)

			send, err := json.Marshal(dataKey)
			if err != nil {
				log.Printf("error: %s", err)
			}
			//send := "Hello2"
			_, err = conn.Write([]byte(send))
			if err != nil {
				fmt.Println("Failed to write data : ", err)
				forCheck = true
				break
			} else {
				fmt.Print("TCP Send data : ")
				fmt.Println(dataKey)
			}

			break
		default:
			n, err := conn.Read(header)
			if err != nil {
				if io.EOF == err {
					break
				}
				log.Printf("Failed Connection: %v\n", err)
				forCheck = true
				break
			}
			if 0 < n {
				messageLength := binary.LittleEndian.Uint32(header)
				realTimeData := readBufferData(messageLength, conn)
				dQ.DataChan <- realTimeData
			}
			time.Sleep(100 * time.Millisecond)
			break
		}
	}
	errorChan <- true
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
