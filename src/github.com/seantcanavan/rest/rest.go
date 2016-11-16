package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/facebookgo/freeport"
	"github.com/gorilla/mux"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/profiler"
	"github.com/seantcanavan/reporter"
)

// The acceptable amount of time between the incoming timestamp and the local timestamp in seconds
// Microsoft recommends a maximum of 5 minutes:https://technet.microsoft.com/en-us/library/jj852172(v=ws.11).aspx
const TIMESTAMP_DELTA_SECONDS = 300

// The key to the query parameter for the incoming timestamp value
const TIMESTAMP = "timestamp"

// The key to the query parameter for the reboot delay value
const REBOOT_DELAY = "delay"

// The key to the query parameter for the remote log email address recipient value
const RECIPIENT_GMAIL = "emailaddress"

// The key to the query parameter for the address where the remote file that is required can be obtained from
const REMOTE_ADDRESS = "remoteupdateurl"

// The subject of the email to send out after a successfuly REST port has been negotiated
const PORT_EMAIL_SUBJECT = "REST Service Successfully Started"

// RestHandler contains all the functionality to interact with this remote
// machine via REST calls. All calls right now require a timestamp that is
// required to be within an acceptable delta to the running machine's timestamp.
// This is designed to prevent replay attacks against the remote host.
// Eventually encryption will be added to authenticate the remote user to
// prevent remote code execution.
type RestHandler struct {
	rtr  *mux.Router
	lgr  *logger.Logger
	Port int
}

// NewRestHandler will return a new RestHandler struct with all of the REST
// endpoints configured. It will also startup the REST server.
func NewRestHandler() (*RestHandler, error) {

	rh := RestHandler{}

	lgr, lgrErr := logger.FromVolatilityValue("rest_package")
	if lgrErr != nil {
		return nil, lgrErr
	}

	rh.lgr = lgr
	rh.rtr = mux.NewRouter()
	rh.rtr.HandleFunc(buildRestPath("execute", TIMESTAMP, REMOTE_ADDRESS), rh.executeHandler)
	rh.rtr.HandleFunc(buildRestPath("reboot", TIMESTAMP, REBOOT_DELAY), rh.rebootHandler)
	rh.rtr.HandleFunc(buildRestPath("sendlogs", TIMESTAMP, RECIPIENT_GMAIL), rh.logHandler)
	rh.rtr.HandleFunc(buildRestPath("forceupdate", TIMESTAMP, REMOTE_ADDRESS), rh.updateHandler)
	rh.rtr.HandleFunc(buildRestPath("updateconfig", TIMESTAMP, REMOTE_ADDRESS), rh.configHandler)
	rh.rtr.HandleFunc(buildRestPath("checkin", TIMESTAMP, RECIPIENT_GMAIL), rh.checkinHandler)

	rh.startupRestServer()
	return &rh, nil
}

func buildRestPath(root string, arguments ...string) string {
	var routeBuf bytes.Buffer
	routeBuf.WriteString("/")
	routeBuf.WriteString(root)

	for _, arg := range arguments {
		routeBuf.WriteString("/")
		routeBuf.WriteString("{")
		routeBuf.WriteString(arg)
		routeBuf.WriteString("}")
	}
	return routeBuf.String()
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

	rh.Port = port
	go http.ListenAndServe(":"+strconv.Itoa(port), rh.rtr)
	rh.lgr.LogMessage("REST server successfully started up on port %v", port)
	reporter.SendPlainEmail(PORT_EMAIL_SUBJECT, []byte(strconv.Itoa(port)))
	return nil
}

// checkinHandler will handle receiving and verifying check-in commands via REST.
// Check-in commands will notify the remote machine that the remote user would
// like the machine to perform a check-in. A check-in will send all pertinent data
// regarding the current operating status of this remote machine.
func (rh *RestHandler) checkinHandler(writer http.ResponseWriter, request *http.Request) {
	rh.lgr.LogMessage("checkinHandler started")
	defer rh.lgr.LogMessage("checkinHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	recipientEmail := queryParams[RECIPIENT_GMAIL]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(recipientEmail); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back a checkin status to the given email address
				// TODO: utilize the email address from the URL query and send to the profiler
				profiler.SendArchiveProfileAsAttachment()
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

// executeHandler will handle receiving and verifying execute commands via REST.
// Execute commands will allow the local machine to execute the code contained
// at the remote location. Currently considering supporting executables and
// Python files. Should we do a JSON config instead to allow call command,
// parameters, and a location to the file to download all cleanly in one?
func (rh *RestHandler) executeHandler(writer http.ResponseWriter, request *http.Request) {
	rh.lgr.LogMessage("executeHandler started")
	defer rh.lgr.LogMessage("executeHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	remoteFileAddress := queryParams[REMOTE_ADDRESS]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(remoteFileAddress); err == nil {
			switch request.Method {
			case "POST":
				// process POST request - download the remote file and execute it
				// download the body and save it as a JSON configuration file
				// instantiate a new instance of loader and point it to the JSON configuration
				// go!
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

// rebootHandler will handle receiving and verifying reboot commands via REST.
func (rh *RestHandler) rebootHandler(writer http.ResponseWriter, request *http.Request) {
	rh.lgr.LogMessage("rebootHandler started")
	defer rh.lgr.LogMessage("rebootHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	rebootDelay := queryParams[REBOOT_DELAY]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		switch request.Method {
		case "POST":
			intDelay, _ := strconv.Atoi(rebootDelay)
			time.Sleep(time.Duration(intDelay) * time.Second)
			// start a new loader, call utils.SysAssetPath("loader_reboot.json")
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

// logHandler will handle receiving and verifying log retrieval commands? via
// REST.
func (rh *RestHandler) logHandler(writer http.ResponseWriter, request *http.Request) {
	rh.lgr.LogMessage("logHandler started")
	defer rh.lgr.LogMessage("logHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	recipientEmail := queryParams[RECIPIENT_GMAIL]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err = rh.verifyQueryParams(recipientEmail); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back the latest logs to the requester
				// collate all the logs
				// zip them up
				// send them via the reporter
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

// updateHandler will handle receiving and verifying update commands via REST.
// Update commands will allow the remote user to force a local update given a
// specific remote URL - should probably be git for now.
func (rh *RestHandler) updateHandler(writer http.ResponseWriter, request *http.Request) {
	rh.lgr.LogMessage("updateHandler started")
	defer rh.lgr.LogMessage("updateHandler finished")

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

// configHandler will handle receiving and verifying config commands via REST.
// Config commands will allow the remote user to set or get the local config
// file that anon-eth-net uses when started up.
func (rh *RestHandler) configHandler(writer http.ResponseWriter, request *http.Request) {
	rh.lgr.LogMessage("configHandler started")
	defer rh.lgr.LogMessage("configHandler finished")

	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	remoteFileAddress := queryParams[REMOTE_ADDRESS]

	if err := rh.verifyTimeStamp(remoteTimestamp); err == nil {
		if err := rh.verifyQueryParams(remoteFileAddress); err == nil {
			switch request.Method {
			case "GET":
				// process GET request - send back the config file
				// bytes:= ioutil.ReadFile(utils.AssetPath(config.json))
				// writer.Body.Write(bytes)
				writer.WriteHeader(http.StatusOK)
			case "POST":
				// process POST request - get the given config file
				// get the asset path
				// create a new temp file
				// load a config from that file, test for errors
				// if no errors, overwrite the existing config
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

// TimeDiffSeconds returns the difference between the input time and the current
// time in seconds. Returns error if the input time stamp cannot be correctly
// converted to a time instance.
func (rh *RestHandler) TimeDiffSeconds(unixTimeStamp string) (*UnixTimeDiff, error) {
	unixDiff := UnixTimeDiff{}
	otherTime, err := strconv.ParseInt(unixTimeStamp, 10, 64)
	if err != nil {
		return nil, err
	}

	unixDiff.then = otherTime
	unixDiff.now = time.Now().Unix()
	unixDiff.diff = unixDiff.now - unixDiff.then
	unixDiff.rawdiff = unixDiff.diff

	if unixDiff.diff < 0 {
		unixDiff.future = true
		unixDiff.diff = unixDiff.diff * -1
	}

	return &unixDiff, nil
}

// verifyTimeStamp will verify the incoming timestamp from the remote machine is
// within an acceptable delta of the current time. Requires tight
// synchronization of both the local time on the local box and the remote time
// on the remote box.
func (rh *RestHandler) verifyTimeStamp(remoteTimeStamp string) error {
	rh.lgr.LogMessage("verifyTimeStamp called with remoteTimeStamp: %v", remoteTimeStamp)

	// get the difference between then and now in seconds from unix time stamps
	diff, diffErr := rh.TimeDiffSeconds(remoteTimeStamp)
	if diffErr != nil || diff.diff >= TIMESTAMP_DELTA_SECONDS {
		return fmt.Errorf("verifyTimeStamp failed with diff: %v", diff.pprint())
	}

	rh.lgr.LogMessage("verifyTimeStamp succeeded with diff: %v", diff.pprint())
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
			rh.lgr.LogMessage("verifyQueryParams failed with: %v", value)
			return fmt.Errorf("verifyQueryParams failed with: %v", value)
		}
	}
	return nil
}

type UnixTimeDiff struct {
	now     int64
	then    int64
	diff    int64
	rawdiff int64
	future  bool
}

func (utd UnixTimeDiff) pprint() string {
	var prettyBuf bytes.Buffer

	prettyBuf.WriteString("UnixTimeDiff:\n")
	prettyBuf.WriteString(fmt.Sprintf("now: %d", utd.now))
	prettyBuf.WriteString(fmt.Sprintf("then: %d", utd.then))
	prettyBuf.WriteString(fmt.Sprintf("diff: %d", utd.diff))
	prettyBuf.WriteString(fmt.Sprintf("rawdiff: %d", utd.rawdiff))
	prettyBuf.WriteString(fmt.Sprintf("future: %t", utd.future))
	return prettyBuf.String()
}
