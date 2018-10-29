package media_manipulation

import (
	. "github.com/donfrigo/img-manipulation/server/helpers"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
)

func ConvertToGreyscale(filePath string) error{

	imgfile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	src, _, err := image.Decode (imgfile)

	// store dimensions of image
	width := src.Bounds().Size().X
	height := src.Bounds().Size().Y

	// Create a new grayscale image
	gray := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			gray.Set(x, y, src.At (x, y) )
		}
	}

	imgfile.Close()
	os.Remove(filePath)

	outfile, _ := os.Create(filePath)
	jpeg.Encode(outfile, gray, &jpeg.Options{80})
	outfile.Close()

	return nil
}

func SplitVideoToFrames(fileHashPath string, fps float64) error {

	fpsString := "fps="+FloatToString(fps,2)

	// split video to frames
	_, err := exec.Command("ffmpeg","-i", path.Join(fileHashPath, "video.mp4"), "-f", "image2", "-vf", fpsString, path.Join(fileHashPath, "frames","frames_%03d.jpg")).Output()

	return err
}

func CalculateFps (maximumNumberOfFrames float64, defaultFps float64, filePath string) float64{

	// number of seconds of video
	command := "ffmpeg -i "+filePath+` 2>&1 | grep 'Duration' | cut -d ' ' -f 4 | sed s/,//  | awk '{ split($1, A, ":"); split(A[3], B, "."); print 3600*A[1] + 60*A[2] + B[1] }'`
	buf, err := exec.Command("bash","-c", command).Output()
	CheckErr(err)

	// convert byte array to int
	s := string(buf)
	s = strings.TrimSuffix(s, "\n")
	seconds, err := strconv.ParseFloat(s,64)
	CheckErr(err)

	if defaultFps * seconds <= maximumNumberOfFrames {
		return defaultFps
	} else {
		return maximumNumberOfFrames / seconds
	}

}

func ProcessImages(fileHashPath string) error {
	framesPath :=path.Join(fileHashPath, "frames")

	errorChan := make(chan error)

	frames, err := ioutil.ReadDir(framesPath)
	if err != nil {
		return err
	}

	// create wait group to wait for all threads to finish
	wg := sync.WaitGroup{}

	// store next file
	queue := make(chan string)

	// limit amount of workers
	workerLimit := 4

	// init all workers
	for worker := 0; worker < workerLimit; worker++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for framePath := range queue {
				// convert each file to greyscale
				err = ConvertToGreyscale(framePath)
				if err != nil {
					errorChan <- err
					return
				}
				errorChan <- nil
			}

		}()

	}

	// add to queue
	for _, frame := range frames {

		framePath := path.Join(framesPath, frame.Name())
		queue <- framePath

		// read from error channel
		r := <-errorChan
		if r != nil {
			close(queue)
			close(errorChan)
			return r
		}
	}

	close(queue)

	// wait for all threads
	wg.Wait()

	return nil
}