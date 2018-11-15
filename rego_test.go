package rego

import (
	"fmt"
	"testing"
)

var (
	confMap map[string]interface{}
)

func init() {
	confMap = make(map[string]interface{})
	confMap["Host"] = "127.0.0.1:6379"
	confMap["Password"] = "123456"
	confMap["Db"] = int64(1)

}

// test redis connect
func Test_Connect(t *testing.T) {
	conn, err := Connect(confMap)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i <= 200000; i++ {
		//conn.Do("Set", "aaa", i)
		conn.Do("GET", "aaa")
	}
	conn.Close()
	t.Log("success")
}

// test redis sadd string slice
func TestSaddSlice(t *testing.T) {
	conn, _ := Connect(confMap)
	slice := []string{"a", "b"}
	SaddStringSlice(conn, "testSet", slice)
}

// test redis connect pool
func Test_GetConnectionPool(t *testing.T) {
	rc, err := GetConnectionPool(confMap)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i <= 100000; i++ {
		conn := rc.Get()
		//conn.Do("Set", "aaa", i)
		conn.Do("GET", "aaa")
		conn.Close()
	}
	t.Log("success")
}

// test redis connect pool
func Test_GetConnectionPoolRoutine(t *testing.T) {
	rc, err := GetConnectionPool(confMap)
	fmt.Println(err)
	if err != nil {
		t.Error(err)
		return
	}
	c := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j <= 10000; j++ {
				conn := rc.Get()
				//conn.Do("Set", "aaa", i)
				conn.Do("GET", "aaa")
				conn.Close()
			}
			fmt.Println("ok")
			c <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-c
	}
	close(c)
	t.Log("success")
}
