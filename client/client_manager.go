package client

import (
	"encoding/json"
	"fmt"
	"log"
	"ontunerestserver/data"
	"ontunerestserver/types"
	"sync"
	"time"
)

var cMapMutex = new(sync.RWMutex)

type ClientRoom struct {
	clientMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
	broadcast chan []byte

	HostLastQueue *data.DataQueue
	BasicQueue    *data.DataQueue
	IOQueue       *data.DataQueue
}

func NewClientRoom() *ClientRoom {
	return &ClientRoom{
		clientMap: make(map[*Client]bool, 5),
		broadcast: make(chan []byte),

		HostLastQueue: data.NewDataQueue(),
		BasicQueue:    data.NewDataQueue(),
		IOQueue:       data.NewDataQueue(),
	}
}

func (w *ClientRoom) dataPump(dQ *types.Queue) { //code uint32,

	//우선 받은것 그대로 보내자.
	// buf := types.RealData{}
	// buf.Code = code

	// if (len(w.clientMap) == 0) { //clientMap 이 0 일때는 realtimedata 도 TCP 로 안받아도 된다.
	// 	return
	// }

	for !dQ.Empty() {
		data := dQ.Pop().([]byte)
		w.broadcast <- data
		// cMapMutex.RLock()
		// for client := range w.clientMap {
		// 	client. <- data
		// }
		// cMapMutex.RUnlock()
	}
	time.Sleep(100 * time.Millisecond) // 0.1s
	// count := dQ.Count()
	// for i := 0; i < count; i++ {
	// 	//buf.Data += dQ.Pop()
	// 	data := dQ.Pop().([]byte)
	// 	for client := range w.clientMap {
	// 		// client.send <- dQ.Pop().([]byte)
	// 		// buf := dQ.Pop().(types.RealData);
	// 		//buf := dQ.Pop().(string);
	// 		// message, err := json.Marshal(buf)
	// 		// if err != nil {
	// 		// 	log.Printf("error: %s", err)
	// 		// }
	// 		client.send <- data
	// 	}
	// }
}

func (w *ClientRoom) RunDatapump() {
	//dataTicker := time.NewTicker(1 * time.Second) //1초 Ticker
	// dataTicker := time.NewTicker(1 * time.Millisecond)
	w.HostLastQueue.SetResponse()
	// w.BasicQueue.SetResponse()
	// w.IOQueue.SetResponse()
	for {
		select {

		case hostlastQ := <-w.HostLastQueue.ResponseChan:
			//fmt.Println("consum hostlast responseChan")
			w.dataPump(hostlastQ)
			// w.HostLastQueue.IsSet = true;
			w.HostLastQueue.SetResponse()

			// case basicQ := <- w.BasicQueue.ResponseChan:
			// 	//fmt.Println("consum basic responseChan")
			// 	w.dataPump(basicQ)
			// 	// w.BasicQueue.IsSet = true;
			// 	w.BasicQueue.SetResponse()

			// case IOQ := <- w.IOQueue.ResponseChan:
			//     //fmt.Println("consum IO responseChan")
			// 	w.dataPump(IOQ)
			// 	w.IOQueue.SetResponse()
			// 	// w.IOQueue.IsSet = true;

			// case <- dataTicker.C:
			// 	w.HostLastQueue.SetResponse()
			// 	w.BasicQueue.SetResponse()
			// 	w.IOQueue.SetResponse()
		}
	}
}

// run 시에 ChanEnter, ChanLeave 를 초기화 한다.
func (w *ClientRoom) Run() {
	w.ChanEnter = make(chan *Client)
	w.ChanLeave = make(chan *Client)

	//pingtTicker := time.NewTicker(1 * time.Second) //1초 Ticker
	// pingtTicker := time.NewTicker(5 * time.Millisecond)

	// w.HostLastQueue.StartResponse()
	// w.BasicQueue.StartResponse()
	// w.IOQueue.StartResponse()

	for {
		select {
		case client := <-w.ChanEnter: //클라이언트 접속시 알림 채널
			fmt.Println("connect client(", client.id, ")")
			w.clientMap[client] = true //클라이언트 맵에 True 표시

		case client := <-w.ChanLeave: //클라이언트 접속 끊길시 알림 채널.
			if _, ok := w.clientMap[client]; ok {
				delete(w.clientMap, client) //맵에서 제거
				fmt.Println("disconnect client(", client.id, ")")
				close(client.send) //client 의 send 채널 close
			}

		case message := <-w.broadcast:

			codebuf := types.DataCode{}
			// codebuf := types.RealData{}
			for client := range w.clientMap {
				err := json.Unmarshal(message, &codebuf)
				if err != nil {
					log.Printf("error: %s", err)
					client.send <- message
					continue
				}
				//fmt.Println("Ws Working")
				if client.dataKey.IsSet(types.Bitmask(codebuf.Code)) {
					client.send <- message
				}
			}

			// case hostlastQ := <- w.HostLastQueue.ResponseChan:
			// 	fmt.Println("consum hostlast responseChan")
			//  	w.dataPump(hostlastQ)
			// 	w.HostLastQueue.IsSet = true;
			// 	//w.HostLastQueue.StartResponse()

			// case basicQ := <- w.BasicQueue.ResponseChan:
			// 	fmt.Println("consum basic responseChan")
			// 	w.dataPump(basicQ)
			// 	w.BasicQueue.IsSet = true;
			// 	//w.BasicQueue.StartResponse()

			// case IOQ := <- w.IOQueue.ResponseChan:
			// fmt.Println("consum IO responseChan")
			// 	w.dataPump(IOQ)
			// 	w.IOQueue.IsSet = true;
			//w.IOQueue.StartResponse()

			// case <- dataTicker.C: 	//datatick :=<- dataTicker.C:
			// 	//w.dataProducer(types.HOST_CODE, w.HostLastQueue)
			// 	// w.dataProducer(types.BASIC_CODE, w.BasicQueue)
			// 	// w.dataProducer(types.NET_CODE, w.IOQueue)
			//     //w.dataSet(datatick.String(), types.HOST_CODE, w.HostLastQueue)
			//     w.dataPump(w.HostLastQueue)
			// 	// w.dataPump(w.BasicQueue)
			// 	// w.dataPump(w.IOQueue)

			// case tick := <- pingtTicker.C: //1초마다 실행
			// 	// w.dataProducer(types.HOST_CODE, w.HostLastQueue)
			// 	// w.dataPump(w.HostLastQueue)

			// 	for client := range w.clientMap { //Client Map 전체 for
			// 		client.send <- []byte(tick.String()) //tick 값 Send >> client.go 의 WritePump 에서 처리
			// 		//client.send <- []byte("1,2,3,4,5,5,6,9") //tick 값 Send >> client.go 의 WritePump 에서 처리
			// 	}
			// 	// w.broadcast <- []byte(tick.String()) //tick 값 Send >> client.go 의 WritePump 에서 처리

			// 	// fmt.Println(tick.String())
			// 	// w.HostLastQueue.SetResponse()
			// 	// w.BasicQueue.SetResponse()
			// 	// w.IOQueue.SetResponse()
		}
	}
}

// func (w *ClientRoom) dataSet(tickString string, code uint32, dataQ *types.Queue) {

// 	var header string

// 	if (code == types.HOST_CODE) {
// 		header = "host, last "
// 	} else if (code == types.BASIC_CODE) {
// 		header = "basic, cpu, mem "
// 	} else if (code == types.NET_CODE) {
// 		header = "net, disk "
// 	}

// 	count := 1000
// 	for i := 0; i < count; i++ {

// 	var buf bytes.Buffer
// 	buf.WriteString(header)
// 	buf.WriteString(tickString)

// 	buf.WriteString(strconv.Itoa(i))

// 	realData := types.RealData{}
// 	realData.Code = code
// 	realData.Data = buf.String()

// 	data, err := json.Marshal(realData)
// 	if err != nil {
// 		log.Printf("error: %s", err)
// 	}

// 	// dataQ.Push([]byte(data))
// 	dataQ.Push(data)
// 	}
// }

// func (w *ClientRoom) dataProducer(code uint32, dataQ *types.Queue) {

// 	ticker := time.NewTicker(1 * time.Second) //time.Millisecond
// 	var header string

// 	if (code == types.HOST_CODE) {
// 		header = "host, last "
// 	} else if (code == types.BASIC_CODE) {
// 		header = "basic, cpu, mem "
// 	} else if (code == types.NET_CODE) {
// 		header = "net, disk "
// 	}

// 	for {
// 		select {
// 			case tick := <- ticker.C: //1초마다 실행

// 			  count := 1000
// 			  for i := 0; i < count; i++ {

// 				var buf bytes.Buffer
// 				buf.WriteString(header)
// 				buf.WriteString(tick.String())

// 				buf.WriteString(strconv.Itoa(i))

// 				realData := types.RealData{}
// 	            realData.Code = code
// 				realData.Data = buf.String()

// 				data, err := json.Marshal(realData)
// 				if err != nil {
// 					log.Printf("error: %s", err)
// 				}
//                 // dataQ.Push([]byte(data))
// 				dataQ.Push(data)
// 			  }
// 		}
// 	}
// }
