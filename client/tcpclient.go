package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"ontunerestserver/types"
	"time"
)

// func writeHandler(conn net.Conn) {
// 	send := "Hello2"
// 	for {
// 		_, err := conn.Write([]byte(send))
// 		if err != nil {
// 			fmt.Println("Failed to write data : ", err)
// 			break;
// 		}

// 		time.Sleep(1 * time.Second)
// 	}
// }

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

func readHandler(conn net.Conn) {
	recv := make([]byte, 4096)
	databuf := types.RealData{}

	for {
		n, err := conn.Read(recv)
		if err != nil {
			fmt.Println("Failed to Read data : ", err)
			break
		}

		var message = recv[:n]
		// fmt.Println(message)
		err = json.Unmarshal(message, &databuf)
		if err != nil {
			log.Printf("error: %s", err)				
			continue				
		}
		fmt.Println("code: ", databuf.Code, " data: ", databuf.Data)
	}
}


func RunClient() {
	//conn, err := net.Dial("tcp", "192.168.0.58:8088")
	conn, err := net.Dial("tcp", "192.168.0.140:8088")	 
	if err != nil {
		fmt.Println("Faield to Dial : ", err)
	}
	defer conn.Close()

	//go writeHandler(conn) 
	go readHandler(conn) 


	// var (
	//     codecBuffer bytes.Buffer
	//     enc         *gob.Encoder = gob.NewEncoder(&codecBuffer)
	// )

	// for {
	//     enc.Encode(MyMsg{
	//         Header: MyMsgHeader{
	//             MsgType: "ping",
	//             Date:    time.Now().UTC().Format(time.RFC3339),
	//         },
	//         Body: MyMsgBodyPing{
	//             Content: "Hello! I'm alive!",
	//         },
	//     })

	//     conn.Write(codecBuffer.Bytes())
	//     codecBuffer.Reset()
	//     time.Sleep(time.Duration(3) * time.Second)
	// }

	// go func(c net.Conn) {
	// 	send := "Hello"
	// 	for {
	// 		_, err = c.Write([]byte(send))
	// 		if err != nil {
	// 			fmt.Println("Failed to write data : ", err)
	// 			break
	// 		}

	// 		time.Sleep(1 * time.Second)
	// 	}
	// }(conn)

	// go func(c net.Conn) {

	// 	recv := make([]byte, 4096)
	// 	databuf := types.RealData{}

	// 	for {
	// 		n, err := c.Read(recv)
	// 		if err != nil {
	// 			fmt.Println("Failed to Read data : ", err)
	// 			break
	// 		}

	// 		var message = recv[:n]

	// 		fmt.Println(message)

	// 		err = json.Unmarshal(message, &databuf)
	// 		if err != nil {
	// 			log.Printf("error: %s", err)
	// 			continue
	// 		}
	// 		fmt.Println("code: ", databuf.Code, " data: ", databuf.Data)
	// 	}
	// }(conn)

	// fmt.Scanln()
}