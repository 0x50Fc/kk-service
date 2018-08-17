package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/gorilla/websocket"
	"github.com/hailongz/kk-service/kk"
	"github.com/hailongz/kk-service/dynamic"
)

func Index(center kk.ICenter) func(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }} // use default options

	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("[" + r.RemoteAddr + "] [OPEN]")

		ch := kk.NewWSChannel(conn)

		defer ch.Close()

		var container kk.IContainer = nil
		var service kk.IService = nil

		c := make(chan bool)

		defer close(c)

		for {

			mType, data, err := conn.ReadMessage()

			if err != nil {
				log.Println("["+r.RemoteAddr+"] [ERROR]", err)
				break
			}

			if mType != websocket.TextMessage {
				log.Println("[" + r.RemoteAddr + "] [ERROR] Message Type Not Is Text")
				break
			}

			var message interface{} = nil

		
			err = json.Unmarshal(data, &message)

			if err != nil {
				log.Println("["+r.RemoteAddr+"] [ERROR]", err)
				break
			}

			type := dynamic.StringValue(dynamic.Get(message,"type"),"");

			if type == "ping" {
				ret := map[string]interface{}
				ret["dtime"] = (time.Now().UnixNano() / 1000000)
				ret["type"] = "pong";

				_,b := json.Marshal(ret)

				err = ch.Send(b)

				if err != nil {
					log.Println("["+r.RemoteAddr+"] [ERROR]", err)
					break
				}
			} else if type == "login" {

				name := dynamic.StringValue(dynamic.Get(message,"name"),"");
				container = center.GetContainer(name, c)
				service = kk.NewService(container, ch, dynamic.IntValue(dynamic.Get(message,"priority"),0), dynamic.StringValue(dynamic.Get(message,"title"),""))

				container.Add(service)
			} else {

				go func() {

					name := dynamic.StringValue(dynamic.Get(message,"name"),"");
					c := make(chan bool)
					defer close(c)

					v := center.GetContainer(name, c)
					s := v.Get(c)

					if s != nil {
						s.Send(data)
					}
					
				}()

			}

		}

		if service != nil && container != nil {
			container.Remove(service)
			service.Exit()
		}

		log.Println("[" + r.RemoteAddr + "] [CLOSE]")

	}
}
