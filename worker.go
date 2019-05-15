package gou

import (
	"github.com/lunny/log"
	"os"
	"strconv"
)

// https://log.zvz.im/2018/02/28/handling-million-requests-with-golang/?hmsr=toutiao.io&utm_medium=toutiao.io&utm_source=toutiao.io
// https://medium.com/smsjunk/handling-1-million-requests-per-minute-with-golang-f70ac505fcaa
// http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

var (
	MaxWorker, _ = strconv.Atoi(os.Getenv("MAX_WORKERS"))
	MaxQueue, _  = strconv.Atoi(os.Getenv("MAX_QUEUE"))
	// A buffered channel that we can send work requests on.
	JobQueue = make(chan Job, MaxQueue)
)

// Job represents the job to be run
type Job interface {
	Do() error
}

func Run() {
	// starting n number of workers
	for i := 0; i < MaxWorker; i++ {
		worker := NewWorker()
		go worker.Start()
	}
}

// Worker represents the worker that executes the job
type Worker struct {
	quit chan bool
}

func NewWorker() Worker {
	return Worker{quit: make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	for {
		select {
		case job := <-JobQueue:
			// we have received a work request.
			if err := job.Do(); err != nil {
				log.Errorf("Error uploading to S3: %s", err.Error())
			}
		case <-w.quit:
			// we have received a signal to stop
			return
		}
	}
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Payload struct {
	// [redacted]
}

/*
func (p *Payload) Do() error {
 the storageFolder method ensures that there are no name collision in
 case we get same timestamp in the key name
storage_path := fmt.Sprintf("%v/%v", p.storageFolder, time.Now().UnixNano())
bucket := S3Bucket
b := new(bytes.Buffer)
encodeErr := json.NewEncoder(b).Encode(payload)
if encodeErr != nil {
	return encodeErr
}
// Everything we post to the S3 bucket should be marked 'private'
var acl = s3.Private
var contentType = "application/octet-stream"
return bucket.PutReader(storage_path, b, int64(b.Len()), contentType, acl, s3.Options{})

	return nil
}

type PayloadCollection struct {
	WindowsVersion string    `json:"version"`
	Token          string    `json:"token"`
	Payloads       []Payload `json:"data"`
}

const MaxLength = 10 * 1024 * 1024 // 10M

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Read the body into a string for json decoding
	var content = &PayloadCollection{}
	err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Go through each payload and queue items individually to be posted to S3
	for _, payload := range content.Payloads {
		JobQueue <- &payload
	}
	w.WriteHeader(http.StatusOK)
}
*/
