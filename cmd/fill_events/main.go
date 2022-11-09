package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	SERVER_ADDRESS    = "https://vk-events.ru/"
	EVENTS_CREATE_URI = "event/create"
	UPLOAD_PHOTO_URI  = "image/upload"
	USER_ID_EUGENIY   = "98278046"
	USER_ID_SASHA     = "133937404"
	USER_ID_SERGEY    = "135142986"
)

func main() {
	data, err := ioutil.ReadFile("events.json")
	if err != nil {
		logrus.Fatal("Error ReadFile:", err)
	}

	var events []models.Event
	err = json.Unmarshal(data, &events)
	if err != nil {
		logrus.Fatal("Error unmarshal:", err)
	}

	for _, value := range events {
		time := time.Now().Unix() + int64(60*60*25*rand.Intn(15))

		value.StartsAt = time

		byteData, err := json.Marshal(value)
		if err != nil {
			logrus.Fatal("Error parse marshal json: ", err)
		}

		req, err := http.NewRequest(http.MethodPost, SERVER_ADDRESS+EVENTS_CREATE_URI, bytes.NewBuffer(byteData))
		if err != nil {
			logrus.Fatal("Error NewRequest: ", err)
		}

		randomNumber := rand.Intn(3)
		xUserId := USER_ID_EUGENIY
		if randomNumber == 0 {
			xUserId = USER_ID_SASHA
		} else if randomNumber == 1 {
			xUserId = USER_ID_SERGEY
		}

		req.Header.Set("x-user-id", xUserId)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logrus.Fatal("Error make request: ", err)
		}

		var gettingEvent models.Event
		err = json.NewDecoder(resp.Body).Decode(&gettingEvent)
		if err != nil {
			logrus.Fatal("Error unmarshal answer: ", err)
		}

		b := new(bytes.Buffer)
		w := multipart.NewWriter(b)
		field, err := w.CreateFormFile("photo0", value.Uid)
		if err != nil {
			logrus.Fatal("Error CreateFormFile: ", err)
		}

		fileOpend, err := os.Open("./photo/" + value.Uid)
		if err != nil {
			logrus.Fatal("Error FileOpen: ", err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(fileOpend)

		field.Write(buf.Bytes())
		if err != nil {
			logrus.Fatal("[Write]: ", err)
		}
		w.Close()

		req, err = http.NewRequest(http.MethodPost, SERVER_ADDRESS+UPLOAD_PHOTO_URI+"?uid="+gettingEvent.Uid, b)
		if err != nil {

		}
		req.Header.Set("x-user-id", xUserId)
		req.Header.Set("Content-Type", w.FormDataContentType())
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp.StatusCode)
	}

	fmt.Println("[Success fill events]")
}
