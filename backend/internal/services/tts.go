package services

import (
	"fmt"
	"os"
	"strings"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

func CustomSplit(text string, maxLength int) []string {
	ans := []string{}

	for len(text) > 0 {
		text = strings.TrimSpace(text) // Always clean up leading spaces

		if len(text) <= maxLength {
			ans = append(ans, text)
			return ans
		}

		// 1. Look at our window
		searchArea := text[:maxLength]

		// 2. Try to find a Period
		cutPoint := strings.LastIndex(searchArea, ".")

		// 3. If no period, try to find a Space (so we don't break words)
		if cutPoint == -1 {
			cutPoint = strings.LastIndex(searchArea, " ")
		}

		// 4. Perform the cut
		if cutPoint == -1 {
			// Absolute fallback: No period AND no space? Hard cut at maxLength.
			ans = append(ans, text[:maxLength])
			text = text[maxLength:]
		} else {
			ans = append(ans, text[:cutPoint+1])
			text = text[cutPoint+1:]
		}
	}
	return ans

}

func MergeAudio(files []string, output string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		// Write the binary data to the master file
		out.Write(data)
	}
	return nil
}

func GetFinalAudio(text string, uniqueID int) (string, error) {
	chunks := CustomSplit(text, 200)
	speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}

	var chunkFiles []string

	for i, chunk := range chunks {
		chunkName := fmt.Sprintf("%d_part_%d", uniqueID, i)

		speech.CreateSpeechFile(chunk, chunkName)
		chunkFiles = append(chunkFiles, "audio/"+chunkName+".mp3")
	}
	uniqueId := fmt.Sprintf("%d", uniqueID)
	finalFileName := "audio/" + uniqueId + "_final.mp3"
	err := MergeAudio(chunkFiles, finalFileName)
	// Clean up the temporary chunk files to save disk space
	for _, file := range chunkFiles {
		os.Remove(file)
	}
	if err != nil {
		return "", err
	}
	return finalFileName, nil
}
