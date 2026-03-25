package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8000"



func RandString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func Auth(role string, client *http.Client, t *testing.T) string {
	var email string
	var password string
	if role == "admin"{
		email = "admin@main"
		password = "1234hdsajh"
	}
	if role == "user"{
		email = "user@user"
		password = "123456789"
	}
	loginPayload := map[string]string{
		"email":    email,
		"password": password,
		"role":     role,
	}
	fmt.Printf("%s",email)
	fmt.Printf("%s", password)

	body, _ := json.Marshal(loginPayload)

	_, err := client.Post(baseURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Post(baseURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var loginResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResp)

	token := loginResp["token"].(string)
	return token
}

func TestBookingFlow(t *testing.T) {
	client := &http.Client{}
	token := Auth("admin",client,t)

	roomPayload := map[string]interface{}{
	"name":        RandString(),
	"description": RandString(),
	"capacity":    10,
	}

	body, _ := json.Marshal(roomPayload)

	req, _ := http.NewRequest("POST", baseURL+"/rooms/create", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var roomResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&roomResp)

	room := roomResp["room"].(map[string]interface{})
	roomID := room["ID"].(string)

	schedulePayload := map[string]interface{}{
		"daysOfWeek": []int{1,2,3},
		"startTime":   "10:00",
		"endTime": "18:00",
	}

	body, _ = json.Marshal(schedulePayload)

	req, _ = http.NewRequest("POST", baseURL+"/rooms/"+roomID+"/schedule/create", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var slotResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&slotResp)
}


func TestReserveBooking(t *testing.T) {
	client := &http.Client{}
	user_token := Auth("user", client, t)
	bookingPayload := map[string]interface{}{
		"slotID": "c4d00d2b-91e3-4b89-b850-65378edf3d65",
		"createConferenceLink":true,
	}
	// вставить существующий ID
	
	body, _ := json.Marshal(bookingPayload)
	req, _ := http.NewRequest("POST", baseURL+"/bookings/create", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+user_token)
	req.Header.Set("Content-Type", "application/json")

	_, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCancelBooking(t *testing.T) {
    client := &http.Client{}
    user_token := Auth("user", client, t)
    bookingID := "a043a6af-2693-47a8-b76a-b85c86422775"
	//вставить существующий ID

    req, err := http.NewRequest("POST", baseURL+"/bookings/"+bookingID+"/cancel", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+user_token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Failed to execute request: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        t.Errorf("Expected status 200 or 204, got %d", resp.StatusCode)
    }
}
