package handlers

import (
	"encoding/json"
	. "github.com/donfrigo/img-manipulation/server/helpers"
	"github.com/donfrigo/img-manipulation/server/media-manipulation"
	"github.com/donfrigo/img-manipulation/server/socket"
	"io"
	"net/http"
	"os"
	"path"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// using websockets to inform user if video has been processed successfully or an error occurred during conversion
	socketId := r.Header.Get("socketId")

	// alternatively use callback
	/*
	callBackUrl := r.URL.Query().Get("callback")

	if (callBackUrl == "") {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Callback url must be provided")
		return
	}
	*/

	// settings
	folder := GetCurrentFolder()
	outputFolder := "output"
	tempFolder := GenerateRandomFolderName()
	filename := "video.mp4"
	// limit filesize
	filesizeLimit := int64(40 * 1024 *1024) // 40 Mb
	maximumNumberOfFrames := float64(500)	// to be extracted from video
	defaultFps := float64(24)

	outputFolderPath := path.Join(folder, outputFolder)
	tempFolderPath := path.Join(outputFolderPath, tempFolder)
	outputFilePath := path.Join(tempFolderPath, filename)

	// create output directory if it doesn't exist
	os.MkdirAll(path.Join(tempFolderPath, "frames"), os.ModePerm)

	// delete temporary folder
	defer os.RemoveAll(tempFolderPath)

	// create file
	file, err := os.Create(outputFilePath)
	CheckErr(err)

	// check file size based on content length
	if r.ContentLength > filesizeLimit {
		w.WriteHeader(413)
		json.NewEncoder(w).Encode("File is too large!")
		return
	}

	// check file size based on read data
	r.Body = http.MaxBytesReader(w, r.Body, filesizeLimit)

	_, err = io.Copy(file, r.Body)
	if err != nil {
		w.WriteHeader(413)
		json.NewEncoder(w).Encode("File is too large!")
		return
	}

	// calculate hash of video
	fileHash := GetFileHash(outputFilePath)

	// new path
	fileHashPath := path.Join(outputFolderPath, fileHash)

	// rename temporary folder to fileHashPath
	err = os.Rename(tempFolderPath, fileHashPath)
	if err != nil {
		// if file has already been uploaded, return
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("File has already been uploaded")
		return
	}

	// calculate fps
	fps := media_manipulation.CalculateFps(maximumNumberOfFrames, defaultFps, path.Join(fileHashPath,"video.mp4"))

	// start video conversion in the background
	go func() {

		// split video to frames
		err = media_manipulation.SplitVideoToFrames(fileHashPath, fps)
		if err != nil {
			// error occurred while splitting the frames
			if len(socketId) != 0 {
				socket.Broadcast(socketId, "error", "Error occurred during conversion")
			}
			return
		}

		// start manipulation of image
		err = media_manipulation.ProcessImages(fileHashPath)
		if err != nil {
			// error occurred while splitting the frames
			if len(socketId) != 0 {
				socket.Broadcast(socketId, "error", "Error occurred during conversion")
			}
			return
		}

		// conversion finished
		// send finished message over socket.io
		if len(socketId) != 0 {
			socket.Broadcast(socketId, "finished", "Conversion finished")
		}

		// or alternatively send GET request to callback url
		//RetryHttpGet(3, time. Second, callBackUrl)
	}()

	// Close connection when upload has finished
	json.NewEncoder(w).Encode("Upload Complete")

}
