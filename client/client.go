package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ontunerestserver/types"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	//SetWriteDeadline은 기본 네트워크 연결에 대한 쓰기 기한을 설정합니다.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	//SetReadDeadline은 기본 네트워크 연결에 대한 읽기 기한을 설정합니다.
	//만약 ping message 보내지 않으면 120 초 후에 connection 은 닫힌다.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	//Ping Message 처리를 위한 .. pingPreid
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	///SetReadLimit는 피어에서 읽은 메시지의 최대 크기(바이트)를 설정합니다.
	// maxMessageSize = 512
	//maxMessageSize = 65536
	//maxMessageSize = 1024
	maxMessageSize = 65536
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	id      uint64
	Room    *ClientRoom
	conn    *websocket.Conn
	send    chan []byte
	dataKey types.Bitmask
}

func NewClient(i uint64, w *ClientRoom, c *websocket.Conn) (client *Client) {
	client = &Client{
		id:   i,
		Room: w,
		conn: c,
		send: make(chan []byte, 65536), //256
	}
	go client.readPump()
	go client.writePump()
	return client
}

func (c *Client) readPump() {
	defer func() {
		c.Room.ChanLeave <- c //readPump 종료시 ChanLeave 를 호출하는구먼.. world 에서 Client 처리
		c.conn.Close()
	}()

	//SetReadLimit는 피어에서 읽은 메시지의 최대 크기(바이트)를 설정합니다. 메시지가 제한을 초과하면 연결은 피어에게 닫기 메시지를 보내고 응용 프로그램에 ErrReadLimit를 반환합니다.
	//c.conn.SetReadLimit(maxMessageSize)

	//SetReadDeadline은 기본 네트워크 연결에 대한 읽기 기한을 설정합니다. 읽기 시간이 초과되면 websocket 연결 상태가 손상되고 이후의 모든 읽기는 오류를 반환합니다. t의 0 값은 읽기가 시간 초과되지 않음을 의미합니다.

	//c.conn.SetReadDeadline(time.Now().Add(pongWait))

	//SetPongHandler는 피어에서 받은 pong 메시지에 대한 핸들러를 설정합니다.
	//h에 대한 appData 인수는 PONG 메시지 애플리케이션 데이터입니다. 기본 퐁 핸들러는 아무 작업도 수행하지 않습니다.
	//처리기 함수는 NextReader, ReadMessage 및 메시지 판독기 Read 메서드에서 호출됩니다. 응용 프로그램은 위의 제어 메시지 섹션에 설명된 대로 pong 메시지를 처리하기 위해 연결을 읽어야 합니다.

	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		codebuf := types.DataCode{}
		err = json.Unmarshal(message, &codebuf)
		if err != nil {
			log.Printf("error: %s", err)
			c.send <- message
			continue
		}

		fmt.Println("code: ", codebuf.Code)
		if codebuf.Code == types.DATAKEY_CODE {

			buf := types.DataKey{}
			err = json.Unmarshal(message, &buf)
			if err != nil {
				log.Printf("error: %s", err)
				c.send <- message
				continue
			}
			fmt.Println("code: ", buf.Code, " key: ", buf.Key)

			if buf.Code == types.DATAKEY_CODE {

				c.dataKey = buf.Key
				var realdata string = "client[" + strconv.Itoa(int(c.id)) + "]"
				realdata += " datakey[" + strconv.Itoa(int(c.dataKey)) + "] "

				if c.dataKey.IsSet(types.HOST_KEY) {
					realdata += ",host"
				}
				if c.dataKey.IsSet(types.LASTPERF_KEY) {
					realdata += ",last"
				}
				if c.dataKey.IsSet(types.BASIC_KEY) {
					realdata += ",basic"
				}
				if c.dataKey.IsSet(types.CPU_KEY) {
					realdata += ",cpu"
				}
				if c.dataKey.IsSet(types.MEM_KEY) {
					realdata += ",mem"
				}
				if c.dataKey.IsSet(types.NET_KEY) {
					realdata += ",net"
				}
				if c.dataKey.IsSet(types.DISK_KEY) {
					realdata += ",disk"
				}
				message = []byte(realdata)
				// message, err = json.Marshal(buf)
				// if err != nil {
				// 	log.Printf("error: %s", err)
				// }
			}
			c.send <- message
		}

		// else if (codebuf.Code == types.HOST_CODE) {
		// 	//....
		// } else if (codebuf.Code == types.LASTPERF_CODE) {
		// 	//.....
		// } else if (codebuf.Code == types.BASIC_CODE) {
		// 	//.....
		// } else if (codebuf.Code == types.CPU_CODE) {
		// 	//.....
		// } else if (codebuf.Code == types.MEM_CODE) {
		// 	//.....
		// } else if (codebuf.Code == types.NET_CODE) {
		// 	//.....
		// } else if (codebuf.Code == types.DISK_CODE) {
		// 	//.....
		// }

		// c.send <- message
		// buf := lib.DataKey{}
		// err = json.Unmarshal(message, &buf)
		// if err != nil {
		// 	log.Printf("error: %s", err)
		// 	c.send <- message
		// 	continue
		// }
		// fmt.Println("code: ", buf.Code, " key: ", buf.Key)

		// if (buf.Code == lib.DATAKEY_CODE) {
		// 	c.dataKey = buf.Key;

		// 	var realdata string = strconv.Itoa(int(c.id));

		// 	realdata += "("+ strconv.Itoa(int(c.dataKey))+")"

		// 	if c.dataKey.IsSet(lib.HOST_KEY) {
		// 		realdata += ",host"
		// 	}
		// 	if c.dataKey.IsSet(lib.LASTPERF_KEY) {
		// 		realdata += ",last"
		// 	}
		// 	if c.dataKey.IsSet(lib.BASIC_KEY) {
		// 		realdata += ",basic"
		// 	}
		// 	if c.dataKey.IsSet(lib.CPU_KEY) {
		// 		realdata += ",cpu"
		// 	}
		// 	if c.dataKey.IsSet(lib.MEM_KEY) {
		// 		realdata += ",mem"
		// 	}
		// 	if c.dataKey.IsSet(lib.NET_KEY) {
		// 		realdata += ",net"
		// 	}
		// 	if c.dataKey.IsSet(lib.DISK_KEY) {
		// 		realdata += ",disk"
		// 	}
		// 	message = []byte(realdata)
		// 	// message, err = json.Marshal(buf)
		// 	// if err != nil {
		// 	// 	log.Printf("error: %s", err)
		// 	// }
		// }
		// c.send <- message
		// c.room.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) //Ping Message 처리를 위한 .. pingPreid
	defer func() {
		ticker.Stop()
		c.conn.Close() //Close는 닫기 메시지를 보내거나 기다리지 않고 기본 네트워크 연결을 닫습니다.
	}()
	for {
		select {
		case message, ok := <-c.send: //world.run 에서 1초 ticker 에 의해 send
			//SetWriteDeadline은 기본 네트워크 연결에 대한 쓰기 기한을 설정합니다. 쓰기 시간이 초과되면 websocket 상태가 손상되고 이후의 모든 쓰기는 오류를 반환합니다. t 값이 0이면 쓰기가 시간 초과되지 않음을 의미합니다.
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				fmt.Println(">>>>>>>>>>>>> send not ok")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) //에러시. CloseMessage 8
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message) //메시지 보낸다.
			//c.conn.WriteMessage(websocket.BinaryMessage, message) //메시지 보낸다.

		case <-ticker.C: //ping message 를 보내지 않을 경우 1분후 종료
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil { return } //PingMessage 로 HandShake?
		}
	}
}
