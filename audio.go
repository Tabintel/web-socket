package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Step 1: Create a new file buffer to store the audio file
	file, err := os.Open("./audio/theTop.m4a") // Replace "path_to_audio_file" with the actual audio file path
	if err != nil {
		fmt.Println("Error opening audio file:", err)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", "audioFile") // Set the file field name as "audio"
	if err != nil {
		fmt.Println("Error writing to buffer:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file to buffer:", err)
		return
	}
	writer.Close()

	// Step 2: Create the HTTP request with the proper endpoint and authorization
	req, err := http.NewRequest("POST", "https://api.gladia.io/audio/text/audio-transcription/", body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-gladia-key", "7370e249-153e-4465-8a8e-b6695d0f2e6a") // Replace "YOUR_API_KEY" with your Gladia API key

	// Step 3: Send the request and handle the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Step 4: Handle the response
	fmt.Println("Response Status:", resp.Status)

	// Handling the JSON response to extract the transcription
	var jsonResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	// Access the 'prediction' field and display the transcribed text
	predictions, ok := jsonResponse["prediction"].([]interface{})
	if !ok {
		fmt.Println("No prediction data found in the response")
		return
	}

	for _, prediction := range predictions {
		predictionMap, ok := prediction.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid prediction format")
			return
		}

		transcription, found := predictionMap["transcription"].(string)
		if !found {
			fmt.Println("No transcription found in the prediction data")
			return
		}

		fmt.Println("Transcription:", transcription)
	}
}
