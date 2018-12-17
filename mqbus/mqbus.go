package mqbus

import (
	"encoding/json"
	"fmt"
	"github.com/jinyidong/go-core/config"
	"github.com/jinyidong/go-core/util/pool"
	"github.com/streadway/amqp"
	"strings"
	"sync"
	"time"
)

type amqpConnection struct {
	Conn *amqp.Connection
}

func (cp *amqpConnection) Dispose() error {
	if cp.Conn != nil {
		return cp.Conn.Close()
	}
	return nil
}

type mqBus struct {
	onceMap       sync.Map
	connectionMap sync.Map
	mu            sync.Mutex
}

var Default = &mqBus{
	onceMap:       sync.Map{},
	connectionMap: sync.Map{},
	mu:            sync.Mutex{},
}

func (mqbus *mqBus) get(connectionString string) (pool.Pool, error) {
	v, ok := mqbus.connectionMap.Load(connectionString)

	if ok {
		return v.(pool.Pool), nil
	} else {
		mqbus.mu.Lock()

		once, _ := mqbus.onceMap.LoadOrStore(connectionString, &sync.Once{})

		var err error

		once.(*sync.Once).Do(func() {
			zkManager := config.NewManager()
			uri := zkManager.Get(connectionString) //"host=10.1.4.131:5672;username=guest;password=guest"
			fmt.Println(uri)
			var connectionPool pool.Pool
			var connection *amqp.Connection
			connectionPool, err = pool.NewChannelPool(3, 3, func() (pool.Object, error) {
				url := parse2Amqp(uri)
				connection, err = amqp.Dial(url)

				if err != nil {
					panic(err)
				}

				if connection == nil {
					panic("connection nil")
				}

				amqp1 := &amqpConnection{Conn: connection}

				receiver := make(chan *amqp.Error)

				amqp1.Conn.NotifyClose(receiver)

				go func(v *amqpConnection) {
					for {
						select {
						case <-receiver:
							{
							Loop:
								for {
									if v != nil {
										v.Dispose()
									}
									url := parse2Amqp(uri)
									connection, err = amqp.Dial(url)

									if err == nil {
										v.Conn = connection
										receiver = make(chan *amqp.Error)
										v.Conn.NotifyClose(receiver)
										fmt.Println("重连成功", &v)
										break Loop
									} else {
										fmt.Println("重连失败,5s后重试", err)

										if connection != nil {
											connection.Close()
										}
										time.Sleep(5 * time.Second)
									}
								}

							}
						default:
							time.Sleep(5 * time.Second)
						}
					}
				}(amqp1)

				return amqp1, nil
			})
			if err != nil {
				panic(err)
			}
			if connectionPool == nil {
				panic("connectionPool nil")
			}
			mqbus.connectionMap.Store(connectionString, connectionPool)
		})
		mqbus.mu.Unlock()

		v, ok := mqbus.connectionMap.Load(connectionString)

		if !ok {
			fmt.Errorf("!ok")
		}
		return v.(pool.Pool), nil
	}
}

func parse2Amqp(amqpConnectionString string) string {

	if amqpConnectionString == "" {
		return ""
	}

	kv := strings.Split(amqpConnectionString, ";")

	var host string
	var username string
	var password string
	var vHost string
	for _, v := range kv {
		vv := strings.Split(v, "=")
		if strings.Contains(v, "host=") {
			host = vv[1]
		}
		if strings.Contains(v, "username=") {
			username = vv[1]
		}
		if strings.Contains(v, "password=") {
			password = vv[1]
		}
		if strings.Contains(v, "virtualHost=") {
			vHost = vv[1]
		}
	}

	if vHost == "" {
		return "amqp://" + username + ":" + password + "@" + host
	} else {
		return "amqp://" + username + ":" + password + "@" + host + "/" + vHost
	}
}

func (mqbus *mqBus) Post(message Message) error {

	connectionPool, _ := mqbus.get(message.ConnectionString())

	pooledObj, err := connectionPool.Get()

	defer pooledObj.Dispose()

	if err != nil {
		return fmt.Errorf("connectionPool.Get: %s", err)
	}

	object := pooledObj.(*pool.PooledObject)

	conn := object.Obj.(*amqpConnection)

	connection := conn.Conn

	channel, err := connection.Channel()

	if err != nil {
		return fmt.Errorf("Channel: %s", err, &conn)
	}

	defer channel.Close()

	//if true {
	//	if err := channel.Confirm(false); err != nil {
	//		return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	//	}
	//
	//	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	//
	//	defer confirmOne(confirms)
	//}

	//var headers = amqp.Table{}
	//
	//headers["x-delay"]=int32(1)

	body, err := json.Marshal(&message)

	if err = channel.Publish(
		message.Exchange(),   // publish to an exchange
		message.RoutingKey(), // routing to 0 or more queues
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			Headers:         amqp.Table{}, //headers,
			ContentType:     "text/plain",
			Type:            message.MessageType(),
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation) {

	if confirmed := <-confirms; confirmed.Ack {
		//log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		//log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
