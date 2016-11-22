package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/facebookgo/freeport"
	"github.com/gorilla/mux"
	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/profiler"
	"github.com/seantcanavan/reporter"
	"github.com/seantcanavan/utils"
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

// The key to the query parameter for the file type to execute for execute handler
const FILE_TYPE = "filetype"

// The subject of the email to send out after a successfuly REST port has been negotiated
const PORT_EMAIL_SUBJECT = "REST Service Successfully Started"

// The REST path name which calls the execute handler
const EXECUTE_REST_PATH = "execute"

// The REST path name which calls the reboot handler
const REBOOT_REST_PATH = "reboot"

// The REST path name which calls the log handler
const LOG_REST_PATH = "logs"

// The REST path name which calls the update handler
const UPDATE_REST_PATH = "update"

// The REST path name which calls the config handler
const CONFIG_REST_PATH = "config"

// The REST path name which calls the check in handler
const CHECKIN_REST_PATH = "checkin"

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
	rh.rtr.HandleFunc(buildGorillaPath(EXECUTE_REST_PATH, TIMESTAMP, FILE_TYPE), rh.executeHandler)
	rh.rtr.HandleFunc(buildGorillaPath(REBOOT_REST_PATH, TIMESTAMP, REBOOT_DELAY), rh.rebootHandler)
	rh.rtr.HandleFunc(buildGorillaPath(LOG_REST_PATH, TIMESTAMP, RECIPIENT_GMAIL), rh.logHandler)
	rh.rtr.HandleFunc(buildGorillaPath(UPDATE_REST_PATH, TIMESTAMP, REMOTE_ADDRESS), rh.updateHandler)
	rh.rtr.HandleFunc(buildGorillaPath(CONFIG_REST_PATH, TIMESTAMP, REMOTE_ADDRESS), rh.configHandler)
	rh.rtr.HandleFunc(buildGorillaPath(CHECKIN_REST_PATH, TIMESTAMP, RECIPIENT_GMAIL), rh.checkinHandler)

	rh.startupRestServer()
	return &rh, nil
}

func buildGorillaPath(root string, arguments ...string) string {
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

func buildRestPath(protocol, host, port, root string, arguments ...string) string {
	var routeBuf bytes.Buffer
	routeBuf.WriteString(protocol)
	routeBuf.WriteString("://")
	routeBuf.WriteString(host)
	routeBuf.WriteString(":")
	routeBuf.WriteString(port)
	routeBuf.WriteString("/")
	routeBuf.WriteString(root)

	for _, arg := range arguments {
		routeBuf.WriteString("/")
		routeBuf.WriteString(arg)
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

// writeResponseAndLog will write the appropriate HTTP status code to the writer
// and also log an appropriate success or failure message to the logger in this
// RestHandler instance.
func (rh *RestHandler) writeResponseAndLog(httpStatusCode int, writer http.ResponseWriter, request *http.Request) {
	var statusBuffer bytes.Buffer

	switch httpStatusCode {
	case http.StatusUnauthorized:
		statusBuffer.WriteString("http.StatusUnauthorized")
	case http.StatusBadRequest:
		statusBuffer.WriteString("http.StatusBadRequest")
	case http.StatusOK:
		statusBuffer.WriteString("http.StatusOK")
	case http.StatusMethodNotAllowed:
		statusBuffer.WriteString("http.StatusMethodNotAllowed")
	default:
		statusBuffer.WriteString(fmt.Sprintf("Unknown HTTP status code: %d", httpStatusCode))
	}

	writer.WriteHeader(httpStatusCode)

	statusBuffer.WriteString(" for writer:")
	statusBuffer.WriteString(fmt.Sprintf("%+v", writer))
	statusBuffer.WriteString("and request:")
	statusBuffer.WriteString(fmt.Sprintf("%+v", &request))
	statusBuffer.WriteString("\n")
	rh.lgr.LogMessage(statusBuffer.String())
	rh.lgr.Flush()
}

// checkinHandler will handle receiving and verifying check-in commands via REST.
// Check-in commands will notify the remote machine that the remote user would
// like the machine to perform a check-in. A check-in will send all pertinent data
// regarding the current operating status of this remote machine.
func (rh *RestHandler) checkinHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	recipientEmail := queryParams[RECIPIENT_GMAIL]

	rh.lgr.LogMessage("checkinHandler - remoteTimestamp: %v recipientEmail: %v", remoteTimestamp, recipientEmail)
	defer rh.lgr.LogMessage("checkinHandler finished")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(http.StatusUnauthorized, writer, request)
		return
	}

	err = rh.verifyQueryParams(recipientEmail)
	if err != nil {
		rh.writeResponseAndLog(http.StatusBadRequest, writer, request)
		return
	}

	switch request.Method {
	case "GET":
		archive, err := profiler.SendArchiveProfileAsAttachment()
		if err != nil {
			rh.lgr.LogMessage("checkinHandler failed to email system profile: %v", err.Error())
			rh.writeResponseAndLog(http.StatusInternalServerError, writer, request)
		} else {
			defer os.Remove(archive.Name())
			rh.writeResponseAndLog(http.StatusOK, writer, request)
		}
	default:
		rh.writeResponseAndLog(http.StatusMethodNotAllowed, writer, request)
	}
	return
}

// executeHandler will handle receiving and verifying execute commands via REST.
// Execute commands will allow the local machine to execute the code contained
// at the remote location. Currently considering supporting executables and
// Python files. Should we do a JSON config instead to allow call command,
// parameters, and a location to the file to download all cleanly in one?
func (rh *RestHandler) executeHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	fileType := queryParams[FILE_TYPE]

	rh.lgr.LogMessage("remoteTimestamp: %v fileType: %v", remoteTimestamp, fileType)
	defer rh.lgr.LogMessage("executeHandler finished")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(http.StatusUnauthorized, writer, request)
		return
	}

	err = rh.verifyQueryParams(fileType)
	if err != nil {
		rh.writeResponseAndLog(http.StatusBadRequest, writer, request)
		return
	}

	switch request.Method {
	case "POST":
		switch fileType {
		case "python":
			// save the bytes to a local file and execute the python interpreter
			rh.lgr.LogMessage("executeHandler is executing remote python file")
			rh.writeResponseAndLog(http.StatusOK, writer, request)
		case "binary":
			// save the bytes to a local file and execute the binary
			rh.lgr.LogMessage("executeHandler is executing remote binary file")
			rh.writeResponseAndLog(http.StatusOK, writer, request)
		case "script":
			// save the bytes to a local file and execute the script
			rh.lgr.LogMessage("executeHandler is executing remote script file")
			rh.writeResponseAndLog(http.StatusOK, writer, request)
		default:
			rh.writeResponseAndLog(http.StatusBadRequest, writer, request)
		}
	default:
		rh.writeResponseAndLog(http.StatusMethodNotAllowed, writer, request)
	}
	return
}

// rebootHandler will handle receiving and verifying reboot commands via REST.
func (rh *RestHandler) rebootHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	rebootDelay := queryParams[REBOOT_DELAY]

	rh.lgr.LogMessage("rebootHandler - remoteTimestamp: %v rebootDelay: %v", remoteTimestamp, rebootDelay)
	defer rh.lgr.LogMessage("rebootHandler finished")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(http.StatusUnauthorized, writer, request)
		return
	}

	err = rh.verifyQueryParams(rebootDelay)
	if err != nil {
		rh.writeResponseAndLog(http.StatusBadRequest, writer, request)
		return
	}

	switch request.Method {
	case "GET":
		intDelay, intErr := strconv.Atoi(rebootDelay)
		if intErr != nil {
			rh.lgr.LogMessage("could not convert reboot parameter to an int: %v", intErr.Error())
			rh.writeResponseAndLog(http.StatusInternalServerError, writer, request)
		} else {
			rh.lgr.LogMessage("sleeping for %d seconds before rebooting", intDelay)
			time.Sleep(time.Duration(intDelay) * time.Second)
			assetPath, assetErr := utils.SysAssetPath("loader_reboot.json")
			if assetErr != nil {
				rh.lgr.LogMessage("could not successfully locate reboot loader JSON file: %v", assetErr.Error())
				rh.writeResponseAndLog(http.StatusInternalServerError, writer, request)
			} else {
				rebootLoader, loaderError := loader.NewLoader(assetPath)
				if loaderError != nil {
					rh.lgr.LogMessage("could not initialize new reboot loader: %v", loaderError.Error())
					rh.writeResponseAndLog(http.StatusInternalServerError, writer, request)
				} else {
					rh.writeResponseAndLog(http.StatusOK, writer, request)
					defer rebootLoader.StartSynchronous()
				}
			}

		}
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

// logHandler will handle receiving and verifying log retrieval commands? via
// REST.
func (rh *RestHandler) logHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	recipientEmail := queryParams[RECIPIENT_GMAIL]

	rh.lgr.LogMessage("logHandler - remoteTimestamp: %v recipientEmail: %v", remoteTimestamp, recipientEmail)
	defer rh.lgr.LogMessage("logHandler finished")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(http.StatusUnauthorized, writer, request)
		return
	}

	err = rh.verifyQueryParams(recipientEmail)
	if err != nil {
		rh.writeResponseAndLog(http.StatusBadRequest, writer, request)
		return
	}

	switch request.Method {
	case "GET":
		rh.lgr.LogMessage("collating logs and sending to gmail address: %v", recipientEmail)
		rh.writeResponseAndLog(http.StatusOK, writer, request)
	default:
		rh.writeResponseAndLog(http.StatusMethodNotAllowed, writer, request)
	}
}

// updateHandler will handle receiving and verifying update commands via REST.
// Update commands will allow the remote user to force a local update given a
// specific remote URL - should probably be git for now.
func (rh *RestHandler) updateHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	remoteFileAddress := queryParams[REMOTE_ADDRESS]

	rh.lgr.LogMessage("updateHandler - remoteTimestamp: %v remoteFileAddress: %v", remoteTimestamp, remoteFileAddress)
	defer rh.lgr.LogMessage("updateHandler finished")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(http.StatusUnauthorized, writer, request)
		return
	}

	err = rh.verifyQueryParams(remoteFileAddress)
	if err != nil {
		rh.writeResponseAndLog(http.StatusBadRequest, writer, request)
		return
	}

	switch request.Method {
	case "GET":
		rh.lgr.LogMessage("commencing manual update using remote address: %v", remoteFileAddress)
		rh.writeResponseAndLog(http.StatusOK, writer, request)
	default:
		rh.writeResponseAndLog(http.StatusMethodNotAllowed, writer, request)
	}
	return
}

// configHandler will handle receiving and verifying config commands via REST.
// Config commands will allow the remote user to set or get the local config
// file that anon-eth-net uses when started up.
func (rh *RestHandler) configHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]

	rh.lgr.LogMessage("configHandler - remoteTimestamp: %v", remoteTimestamp)
	defer rh.lgr.LogMessage("configHandler finished")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(http.StatusUnauthorized, writer, request)
		return
	}

	switch request.Method {
	case "GET":
		rh.lgr.LogMessage("received remote request for the current config. returning via REST")
		rh.writeResponseAndLog(http.StatusOK, writer, request)
	case "POST":
		rh.lgr.LogMessage("received remote request to save a new config. downloading via REST")
		rh.writeResponseAndLog(http.StatusOK, writer, request)
	}
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
		return fmt.Errorf("verifyTimeStamp failed with diff: %v", diff.diff)
	}

	rh.lgr.LogMessage("verifyTimeStamp succeeded with diff: %v", diff.diff)
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
	prettyBuf.WriteString(fmt.Sprintf("now: %d\t", utd.now))
	prettyBuf.WriteString(fmt.Sprintf("then: %d\t", utd.then))
	prettyBuf.WriteString(fmt.Sprintf("diff: %d\t", utd.diff))
	prettyBuf.WriteString(fmt.Sprintf("rawdiff: %d\t", utd.rawdiff))
	prettyBuf.WriteString(fmt.Sprintf("future: %t\n", utd.future))
	return prettyBuf.String()
}
