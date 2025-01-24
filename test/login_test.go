package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostLoginAndVerify(t *testing.T) {
    // Вызываем обработчик с тестовым запросом и записью ответа
	url := "http://localhost:5555/Login"; // Замените на URL вашего сервера
	data := []byte(`{"phone_number": "8-999-999-99-99", "password": "password", "id_api": 1}`); // Данные для POST-запроса

	client := &http.Client{};
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data));
	assert.Equal(t, err, nil);
	defer resp.Body.Close();

	assert.Equal(t, resp.StatusCode, http.StatusOK);
	body, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil);
	type Output struct {
		Success bool `json:"success"`
		Token string `json:"token"`
		IdUser int `json:"id_user"`;
	};

	var out Output;
	err = json.Unmarshal(body, &out);
	assert.Equal(t, err, nil);

	url = "http://localhost:5555/Verify"; // Замените на URL вашего сервера
	data = []byte(fmt.Sprintf(`{"token": "%s", "id_api": %d}`, out.Token, 1)); // Данные для POST-запроса

	resp2, err := client.Post(url, "application/json", bytes.NewBuffer(data));
	assert.Equal(t, err, nil);
	defer resp2.Body.Close();

	assert.Equal(t, resp2.StatusCode, http.StatusOK);
}