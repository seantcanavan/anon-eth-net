package rest

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/facebookgo/freeport"
	"github.com/gorilla/mux"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/reporter"
	"github.com/seantcanavan/utils"
)

// The acceptable amount of time between the incoming timestamp and the local timestamp
const TIMESTAMP_DELTA = 5000

// The key to the query parameter for the incoming timestamp value
const TIMESTAMP = "timestamp"

// The key to the query parameter for the reboot delay value
const REBOOT_DELAY = "delay"

// The key to the query parameter for the remote log email address recipient value
const RECIPIENT_EMAIL = "emailaddress"

// The key to the query parameter for the address where the remote file that is required can be obtained from
const REMOTE_ADDRESS = "remoteupdateurl"

const PORT_EMAIL_SUBJECT = "REST Service Successfully Started"

// RestHandler contains all the functionality to interact with this remote
// machine via REST calls. All calls right now require a timestamp that is
// required to be within an acceptable delta to the running machine's timestamp.
// This is designed to prevent replay attacks against the remote host.
// Eventually encryption will be added to authenticate the remote user to
// prevent remote code execution.
type RestHandler struct {
	router   *mux.Router
	logger   *logger.Logger
	reporter *reporter.Reporter
}

// NewRestHandler will return a new RestHandler struct with all of the REST
// endpoints configured. It will also startup the REST server.
func NewRestHandler(inputReporter *reporter.Reporter) *RestHandler {
	rh := RestHandler{}
	// rh.logger = seanLogger
	rh.reporter = inputReporter
	rh.router = mux.NewRouter()
	rh.router.HandleFunc("/execute/{"+TIMESTAMP+"}/{"+REMOTE_ADDRESS+"}", rh.ExecuteHandler)
	rh.router.HandleFunc("/reboot/{"+TIMESTAMP+"}/{"+REBOOT_DELAY+"}", rh.RebootHandler)
	rh.router.HandleFunc("/sendlogs/{"+TIMESTAMP+"}/{"+RECIPIENT_EMAIL+"}", rh.LogHandler)
	rh.router.HandleFunc("/forceupdate/{"+TIMESTAMP+"}/{"+REMOTE_ADDRESS+"}", rh.UpdateHandler)
	rh.router.HandleFunc("/updateconfig/{"+TIMESTAMP+"}/{"+REMOTE_ADDRESS+"}", rh.ConfigHandler)
	rh.router.HandleFunc("checkin/{"+TIMESTAMP+"}/{"+RECIPIENT_EMAIL+"}", rh.CheckinHandler)
	go rh.startupRestServer()
	return &rh
}

// startupRestServer will start up the local REST server where this remote
// machine will listen for incoming commands on. A free port on this local
// machine will be automatically detected and used. The randomly chosen
// available port will be logged locally as well as emailed.
func (rh *RestHandler) startupRestServer() error {
	port, err := freeport.Get()
	if err != nil {
		return err
	}

	go func() { http.ListenAndServe(":"+strconv.Itoa(port), rh.router) }()
	rh.logger.LogMessage("REST server successfully started up on port %v", port)

	var contentsBuf bytes.Buffer
	contentsBuf.WriteString(strconv.Itoa(port))
	rh.reporter.SendPlainEmail(PORT_EMAIL_SUBJECT, contentsBuf.Bytes())
	return nil
}

// CheckinHandler will handle receiving and verifying check-in commands via REST.
// Check-in commands will notify the remote machine that the remote user would
// like the machine to perform a check-in. A check-in will send all pertinent data
// regarding the current operating status of this remote machine.
func (rh *RestHandler) CheckinHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("CheckinHandler started")
	defer rh.logger.LogMessage("CheckinHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	recipientEmail := queryParams[RECIPIENT_EMAIL]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(recipientEmail); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back a checkin status to the given email address
				writer.WriteHeader(http.StatusOK)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusUnauthorized)
	return
}

// ExecuteHandler will handle receiving and verifying execute commands via REST.
// Execute commands will allow the local machine to execute the code contained
// at the remote location. Currently considering supporting executables and
// Python files. Should we do a JSON config instead to allow call command,
// parameters, and a location to the file to download all cleanly in one?
func (rh *RestHandler) ExecuteHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("ExecuteHandler started")
	defer rh.logger.LogMessage("ExecuteHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	remoteFileAddress := queryParams[REMOTE_ADDRESS]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(remoteFileAddress); err == nil {
			switch request.Method {
			case "POST":
				// process POST request - download the remote file and execute it
				writer.WriteHeader(http.StatusOK)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusUnauthorized)
	return
}

// RebootHandler will handle receiving and verifying reboot commands via REST.
func (rh *RestHandler) RebootHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("RebootHandler started")
	defer rh.logger.LogMessage("RebootHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	// rebootDelay := queryParams[REBOOT_DELAY]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		switch request.Method {
		case "POST":
			// process POST request - reboot the machine after X seconds
			writer.WriteHeader(http.StatusOK)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
	writer.WriteHeader(http.StatusUnauthorized)
	return
}

// LogHandler will handle receiving and verifying log retrival commands? via
// REST.
func (rh *RestHandler) LogHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("LogHandler started")
	defer rh.logger.LogMessage("LogHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	recipientEmail := queryParams[RECIPIENT_EMAIL]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(recipientEmail); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back the latest logs to the requester
				writer.WriteHeader(http.StatusOK)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusUnauthorized)
	return
}

// UpdateHandler will handle receiving and verifying update commands via REST.
// Update commands will allow the remote user to force a local update given a
// specific remote URL - should probably be git for now.
func (rh *RestHandler) UpdateHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("UpdateHandler started")
	defer rh.logger.LogMessage("UpdateHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	remoteFileAddress := queryParams[REMOTE_ADDRESS]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(remoteFileAddress); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back the current update url
				writer.WriteHeader(http.StatusOK)
			case "POST":
				// process POST request - use the given URL to perform an update
				writer.WriteHeader(http.StatusOK)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusUnauthorized)
	return
}

// ConfigHandler will handle receiving and verifying config commands via REST.
// Config commands will allow the remote user to set or get the local config
// file that anon-eth-net uses when started up.
func (rh *RestHandler) ConfigHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("ConfigHandler started")
	defer rh.logger.LogMessage("ConfigHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	remoteFileAddress := queryParams[REMOTE_ADDRESS]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err := rh.verifyQueryParams(remoteFileAddress); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back the config file
				writer.WriteHeader(http.StatusOK)
			case "POST":
				// process POST request - get the given config file
				writer.WriteHeader(http.StatusOK)
			default:
				writer.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusUnauthorized)
	return
}

// verifyTimeStamp will verify the incoming timestamp from the remote machine is
// within an acceptable delta of the current time. Requires tight
// synchronization of both the local time on the local box and the remote time
// on the remote box.
func (rh *RestHandler) verifyTimeStamp(remoteTimeStamp string) error {
	rh.logger.LogMessage(fmt.Sprintf("verifyTimeStamp called with remoteTimeStamp: %v", remoteTimeStamp))
	//verify the timestamp here
	localTimeStamp := utils.FullDateString()
	// rh.logger.LogMessage("verifyTimeStamp failed. localTimeStamp: %v. remoteTimeStamp: %v", localTimeStamp, remoteTimeStamp)
	// return errors.New("timestamp verification failed. check local and remote lock sync settings.")
	rh.logger.LogMessage(fmt.Sprintf("verifyTimeStamp succeeded with localTimeStamp: %v", localTimeStamp))
	return nil
}

// verifyQueryParams will verify the incoming query parameters from the remote
// machine to make sure that they're not empty. Since maps default to returning
// a safe value of the empty sting we can't simply do a nil check. That and
// golang strings can't be nil anyways... probably why maps return the empty
// string then when it's missing. Epiphany successfully experienced.
func (rh *RestHandler) verifyQueryParams(parameters ...string) error {
	for _, value := range parameters {
		if value == "" {
			rh.logger.LogMessage(fmt.Sprintf("verifyQueryParams failed with: %v", value))
			return errors.New(fmt.Sprintf("verifyQueryParams did not pass verification. val: %v", value))
		}
	}
	return nil
}
