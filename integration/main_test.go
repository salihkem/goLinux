package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/salihkemaloglu/UnitAndIntegrationTesting-Golang/operations"
)

func TestHttpRequestGetAll(t *testing.T) {
	response, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	if response != nil {
		if len(response) != 3 {
			t.Fatal("Expected value: 3 Received value:", len(response))
		}
	}

}
func TestHttpRequestInsert(t *testing.T) {
	responseGetBefore, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	item := data.Item{
		Name:        "hey",
		Value:       "val",
		Description: "desc",
	}
	bytesRepresentation, err := json.Marshal(item)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	response, err := http.Post("http://localhost:8080/item", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	} else {
		defer response.Body.Close()
		var item data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			fmt.Printf("Json decode error!: %s", err)
		}
	}
	responseGetAfter, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount++
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Insert Fail! Before Insert: %v, After Insert: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func TestHttpRequestUpdate(t *testing.T) {
	responseGetBefore, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	itemGet := responseGetBefore[len(responseGetBefore)-1]
	itemGet.Name = "UpdateName"
	itemGet.Value = "UpdateValue"
	itemGet.Description = "UpdateDesc"
	url := "http://localhost:8080/item/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	} else {
		defer response.Body.Close()
		// contents, err := ioutil.ReadAll(response.Body)
		// if err != nil {
		// 	fmt.Printf("Json decode error!: %s", err)
		// }
		// fmt.Println("The update result is:", string(contents))
	}

	responseGetAfter, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	itemUpdate := responseGetAfter[len(responseGetAfter)-1]
	if itemGet.Name != itemUpdate.Name {
		t.Fatal(fmt.Printf("Update Fail! Before Update: %v, After Update: %v \n", itemGet, itemUpdate))
	}

}

func TestHttpRequestDelete(t *testing.T) {
	responseGetBefore, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	itemGet := responseGetBefore[0]
	url := "http://localhost:8080/item/" + bson.ObjectId(itemGet.ID).Hex()
	bytesRepresentation, err := json.Marshal(itemGet)
	if err != nil {
		fmt.Printf("Json decode error!: %s", err)
	}
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
		os.Exit(1)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request Error!: %s", err)
	} else {
		defer response.Body.Close()
	}

	responseGetAfter, err := GetAll()
	if err != nil {
		t.Fatal("End point does not responde!", err.Error())
	}
	responseGetBeforeCount := len(responseGetBefore)
	responseGetBeforeCount--
	if responseGetBeforeCount != len(responseGetAfter) {
		t.Fatal(fmt.Printf("Delete Fail! Before Delete: %v, After Delete: %v \n", responseGetBeforeCount, len(responseGetAfter)))
	}
}

func GetAll() ([]data.Item, error) {
	response, err := http.Get("http://localhost:8080/item")
	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		var item []data.Item
		if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
			fmt.Printf("Json decode error!: %s", err)
		}
		return item, err
	}
}