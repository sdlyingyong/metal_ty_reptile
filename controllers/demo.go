package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"time"
)

//this is some learn code for go example
//see https://gobyexample.com
type DemoController struct {
	beego.Controller
}

func (c *DemoController) ChanDemo() {

	//goroutine()

	//demoChan()
	//chan_message()
	//chan_message_buffer()
	//chan_sync()
	//chan_defections()
	//chan_select()
	//chan_select_timeout()
	//chan_non_block()
	//chan_close()
	//chan_range()


}

//print three times
func f(from string) {
	for i := 0; i < 3; i++ {
		logs.Info(from + ":" + strconv.Itoa(i))
	}
}

func goroutine() {
	f("direct")

	//another executioner
	go f("goroutine")

	//another executioner
	go func(msg string) {
		logs.Info(msg)
	}("going")

	//wait
	time.Sleep(time.Second)
	logs.Info("done")

}

func chan_range() {

	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	for value := range queue {
		logs.Info(value)
	}

}

func chan_close() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	//open another worker
	//it receive job and if no more
	//it send true in done and close
	go func() {
		for {
			j, more := <-jobs
			if more {
				logs.Info("received job", j)
			} else {
				logs.Info("received all jobs")
				done <- true
				return
			}
		}
	}()

	//send 3 job in jobs channel to worker
	//and clone worker channel
	//wait done return with <- chan
	for j := 1; j <= 3; j++ {
		jobs <- j
		logs.Info("sent job" + strconv.Itoa(j))
	}
	close(jobs)
	logs.Info("sent all jobs")
	<-done
}

func chan_non_block() {
	messages := make(chan string)
	signals := make(chan bool)

	//set default, if no channel receive data ,
	//run default and go next line
	select {
	case msg := <-messages:
		logs.Info("received message " + msg)
	default:
		logs.Info("no message received")
	}

	//set default
	//if no receive data from channel
	//use default and go to next line
	msg := "hi"
	select {
	case messages <- msg:
		logs.Info("sent message" + msg)
	default:
		logs.Info("no message sent")
	}

	//set default let not wait here
	//if no receive any one
	//use default and go to next line
	select {
	case msg := <-messages:
		logs.Info("received message" + msg)
	case sig := <-signals:
		logs.Info("received signal " + strconv.FormatBool(sig))
	default:
		logs.Info("no activity")
	}

}

func chan_select_timeout() {
	c1 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c1 <- "result 1"
	}()

	select {
	case res := <-c1:
		logs.Info(res)
	case <-time.After(1 * time.Second):
		logs.Info("time out 1")
	}

	c2 := make(chan string, 1)
	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "result 2"
	}()

	select {
	case res := <-c2:
		logs.Info(res)
	case <-time.After(3 * time.Second):
		logs.Info("time out 2")

	}

	//wait if one of case return
	select {
	//get data from channel default
	case <-c2:
		logs.Info("i get return")
	//if time out of setting,it's error
	case <-time.After(1 * time.Second):
		logs.Error("to loog time err")
	}

}

func chan_select() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		time.Sleep(time.Second)
		c1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		c2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		//execute two task in same time
		//return time is large one , not 1 + 2
		case msg1 := <-c1:
			logs.Info("received " + msg1)
		case msg2 := <-c2:
			logs.Info("received " + msg2)
		}
	}

}

func chan_derections() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	logs.Info(<-pongs)
}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	logs.Info("get message from pings : " + msg)
	pongs <- msg
	logs.Info("send message in pongs : " + msg)
}

func ping(pings chan<- string, msg string) {
	pings <- msg
	logs.Info("send msg in pings :" + msg)
}

func chan_sync() {
	done := make(chan bool, 1)
	go worker(done)
	<-done
}

func worker(done chan bool) {
	logs.Info("working...")
	time.Sleep(time.Second)
	logs.Info("done")

	done <- true
}

func chan_message_buffer() {
	messages := make(chan string, 2)

	messages <- "buffered"
	messages <- "channel"

	logs.Info(<-messages)
	logs.Info(<-messages)
}

func chan_message() {
	messages := make(chan string)

	go func() {
		messages <- "ping"
	}()

	msg := <-messages
	//fmt.Println(msg)
	logs.Info(msg)

}

func demoChan() {
	demoChan := make(chan int)

	go func() {
		demoChan <- 1
	}()

	value := <-demoChan
	logs.Info(value)
}
