package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/facebookgo/freeport"
	"github.com/gorilla/mux"
	"github.com/seantcanavan/anon-eth-net/config"
	"github.com/seantcanavan/anon-eth-net/loader"
	"github.com/seantcanavan/anon-eth-net/logger"
	"github.com/seantcanavan/anon-eth-net/profiler"
	"github.com/seantcanavan/anon-eth-net/reporter"
	"github.com/seantcanavan/anon-eth-net/utils"
)

// The acceptable amount of time between the incoming timestamp and the local timestamp in seconds
// Microsoft recommends a maximum of 5 minutes: https://technet.microsoft.com/en-us/library/jj852172(v=ws.11).aspx
const TIMESTAMP_DELTA_SECONDS = 300

// The key to the query parameter for the incoming timestamp value
const TIMESTAMP = "timestamp"

// The key to the query parameter for the reboot delay value
const REBOOT_DELAY = "delay"

// The key to the query parameter for the file type to execute for execute handler
const FILE_TYPE = "filetype"

// The key to the query parameter for the asset file name to perform CRUD operations on over REST
const ASSET_NAME = "assetname"

// The subject of the email to send out after a successfully REST port has been negotiated
const REST_EMAIL_SUBJECT = "REST Service Successfully Started"

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

// The REST path name which calls the asset handler
const ASSET_REST_PATH = "asset"

// The subject of the email to send out when the REST package is finished executing remote code via the loader package
const REST_LOADER_SUBJECT = "Rest Execute Handler Results"

// RestHandler contains all the functionality to interact with this remote
// machine via REST calls. All calls right now require a timestamp that is
// required to be within an acceptable delta to the running machine's timestamp.
// This is designed to prevent replay attacks against the remote host.
// Eventually encryption will be added to authenticate the remote user to
// prevent remote code execution.
type RestHandler struct {
	rtr       *mux.Router
	Port      string
	Endpoints map[string]string
}

// NewRestHandler will return a new RestHandler struct with all of the REST
// endpoints configured. It will also startup the REST server.
func NewRestHandler() (*RestHandler, error) {

	rh := RestHandler{}

	rh.Endpoints = make(map[string]string)
	rh.Endpoints[LOG_REST_PATH] = buildGorillaPath(LOG_REST_PATH, TIMESTAMP)
	rh.Endpoints[REBOOT_REST_PATH] = buildGorillaPath(REBOOT_REST_PATH, TIMESTAMP, REBOOT_DELAY)
	rh.Endpoints[UPDATE_REST_PATH] = buildGorillaPath(UPDATE_REST_PATH, TIMESTAMP)
	rh.Endpoints[CHECKIN_REST_PATH] = buildGorillaPath(CHECKIN_REST_PATH, TIMESTAMP)
	rh.Endpoints[EXECUTE_REST_PATH] = buildGorillaPath(EXECUTE_REST_PATH, TIMESTAMP, FILE_TYPE)
	rh.Endpoints[ASSET_REST_PATH] = buildGorillaPath(ASSET_REST_PATH, TIMESTAMP, ASSET_NAME)

	logger.Lgr.LogMessage("Successfully generated REST endpoint map: %+v", rh.Endpoints)

	rh.rtr = mux.NewRouter()
	rh.rtr.HandleFunc(rh.Endpoints[LOG_REST_PATH], rh.logHandler)
	rh.rtr.HandleFunc(rh.Endpoints[REBOOT_REST_PATH], rh.rebootHandler)
	rh.rtr.HandleFunc(rh.Endpoints[UPDATE_REST_PATH], rh.updateHandler)
	rh.rtr.HandleFunc(rh.Endpoints[CHECKIN_REST_PATH], rh.checkinHandler)
	rh.rtr.HandleFunc(rh.Endpoints[EXECUTE_REST_PATH], rh.executeHandler)
	rh.rtr.HandleFunc(rh.Endpoints[ASSET_REST_PATH], rh.assetHandler)

	logger.Lgr.LogMessage("Successfully generated REST gorilla mux router: %+v", rh.rtr)

	logger.Lgr.LogMessage("Started up TLS REST server")
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
func (rh *RestHandler) StartupRestServer() error {
	port, err := freeport.Get()
	if err != nil {
		return err
	}

	rh.Port = strconv.Itoa(port)

	pKeyPath, pKeyPathErr := utils.AssetPath("server.pkey")
	if pKeyPathErr != nil {
		return pKeyPathErr
	}

	logger.Lgr.LogMessage("Successfully located private key asset: %v", pKeyPath)

	certPath, certPathErr := utils.AssetPath("server.cert")
	if certPathErr != nil {
		return certPathErr
	}

	logger.Lgr.LogMessage("Successfully located server cert asset: %v", certPath)

	go http.ListenAndServeTLS(":"+rh.Port, certPath, pKeyPath, rh.rtr)

	logger.Lgr.LogMessage("REST server successfully started up on port %v", port)

	externalIp, extIpErr := utils.ExternalIPAddress()
	if extIpErr != nil {
		logger.Lgr.LogMessage("Failed to retrieve external IP address: %v", extIpErr)
		return reporter.SendPlainEmail(REST_EMAIL_SUBJECT, []byte(strconv.Itoa(port)))
	}

	logger.Lgr.LogMessage("Successfully retrieved external IP: %v", externalIp)

	var baseRestPath bytes.Buffer
	baseRestPath.WriteString("https://")
	baseRestPath.WriteString(externalIp)
	baseRestPath.WriteString(":")
	baseRestPath.WriteString(rh.Port)

	var emailBody bytes.Buffer

	for _, value := range rh.Endpoints {
		emailBody.WriteString(baseRestPath.String())
		emailBody.WriteString(value)
		emailBody.WriteString("\n")
	}

	logger.Lgr.LogMessage("Sending out full REST path specs via email")

	return reporter.SendPlainEmail(REST_EMAIL_SUBJECT, emailBody.Bytes())
}

// writeResponseAndLog will write the appropriate HTTP status code to the writer
// and also log an appropriate success or failure message to the logger in this
// RestHandler instance.
func (rh *RestHandler) writeResponseAndLog(errorMessage string, httpStatusCode int, writer http.ResponseWriter, request *http.Request) {
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

	if errorMessage != "" {
		logger.Lgr.LogMessage(errorMessage)
	}

	logger.Lgr.LogMessage(statusBuffer.String())
}

// checkinHandler will handle receiving and verifying check-in commands via REST.
// Check-in commands will notify the remote machine that the remote user would
// like the machine to perform a check-in. A check-in will send all pertinent data
// regarding the current operating status of this remote machine.
func (rh *RestHandler) checkinHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]

	logger.Lgr.LogMessage("checkinHandler - remoteTimestamp: %v recipientEmail: %v", remoteTimestamp, config.Cfg.CheckInGmailAddress)
	defer logger.Lgr.LogMessage("checkinHandler finished\n")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusUnauthorized, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully validated incoming timestamp")

	switch request.Method {

	case "GET":
		logger.Lgr.LogMessage("Received http.GET request - sending profile out via email")
		archive, err := profiler.SendArchiveProfileAsAttachment()
		if err != nil {
			logger.Lgr.LogMessage("checkinHandler failed to email system profile: %v", err.Error())
			rh.writeResponseAndLog(err.Error(), http.StatusInternalServerError, writer, request)
		} else {
			defer os.Remove(archive.Name())
			logger.Lgr.LogMessage("Successfully emailed out system profile via email")
			rh.writeResponseAndLog("", http.StatusOK, writer, request)
		}
	default:
		logger.Lgr.LogMessage("Received unsupported REST method %v for checkinHandler", request.Method)
		rh.writeResponseAndLog("", http.StatusMethodNotAllowed, writer, request)
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

	logger.Lgr.LogMessage("remoteTimestamp: %v fileType: %v", remoteTimestamp, fileType)
	defer logger.Lgr.LogMessage("executeHandler finished\n")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusUnauthorized, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully validated incoming timestamp")

	err = rh.verifyQueryParams(fileType)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusBadRequest, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully verified query parameters")

	bodyContents, bodyErr := ioutil.ReadAll(request.Body)
	if bodyErr != nil {
		rh.writeResponseAndLog(bodyErr.Error(), http.StatusBadRequest, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully read in request.Body")

	switch request.Method {
	case "POST":
		switch fileType {
		case "python", "binary", "script":
			logger.Lgr.LogMessage("executeHandler is executing remote %v file", fileType)
			// save the bytes to a local file and execute the file in the appropriate manner
			loaderError := rh.executeLoader(fileType, bodyContents)
			if loaderError != nil {
				logger.Lgr.LogMessage("error executing remote code: %v", loaderError.Error())
				rh.writeResponseAndLog(loaderError.Error(), http.StatusBadRequest, writer, request)
				return
			}
			logger.Lgr.LogMessage("Successfully executed remote code")
			rh.writeResponseAndLog("", http.StatusOK, writer, request)
		default:
			logger.Lgr.LogMessage("Received unsupported code type: %v", fileType)
			rh.writeResponseAndLog("", http.StatusBadRequest, writer, request)
		}
	default:
		logger.Lgr.LogMessage("Received unsupported REST method %v for executeHandler", request.Method)
		rh.writeResponseAndLog("", http.StatusMethodNotAllowed, writer, request)
	}
	return
}

func (rh *RestHandler) executeLoader(fileType string, fileContents []byte) error {

	processMap := make(map[string]string)
	fileName := utils.FullDateStringSafe() + ".run"

	logger.Lgr.LogMessage("Successfully created temp file for rest execute loader: %v", fileName)

	writeErr := ioutil.WriteFile(fileName, fileContents, 0777)
	if writeErr != nil {
		return writeErr
	}

	logger.Lgr.LogMessage("Successfully copied bytes from REST body to tmpfile")

	defer os.Remove(fileName)

	absPath, pathErr := filepath.Abs(fileName)
	if pathErr != nil {
		return pathErr
	}

	switch fileType {

	case "python":
		processMap["rest_loader_python"] = "python " + absPath
		logger.Lgr.LogMessage("Generated REST loader python execute command")

	case "binary":
		processMap["rest_loader_binary"] = absPath
		logger.Lgr.LogMessage("Generated REST loader binary execute command")

	case "script":
		processMap["rest_loader_script"] = "/bin/sh " + absPath
		logger.Lgr.LogMessage("Generated REST loader script execute command")

	}

	jsonString, jsonErr := json.Marshal(processMap)
	if jsonErr != nil {
		return jsonErr
	}

	logger.Lgr.LogMessage("Successfully marshaled REST process JSON into a map")

	tmpLoaderFile, tmpLoaderErr := ioutil.TempFile("", "restLoader.json")
	if tmpLoaderErr != nil {
		return tmpLoaderErr
	}

	logger.Lgr.LogMessage("Successfully created tmp rest loader file: %v", tmpLoaderFile.Name())

	defer os.Remove(tmpLoaderFile.Name())

	bufferJsonContents := bytes.NewBufferString(string(jsonString))

	_, copiedErr := io.Copy(tmpLoaderFile, bufferJsonContents)
	if copiedErr != nil {
		return copiedErr
	}

	logger.Lgr.LogMessage("Successfully copied JSON REST loader bytes into new loader process file")

	restLoader, loaderErr := loader.NewLoader(tmpLoaderFile.Name())
	if loaderErr != nil {
		return loaderErr
	}

	logger.Lgr.LogMessage("Successfully instantiated a new Loader instance: %+v", restLoader)

	finishedProcesses := restLoader.StartSynchronous()
	logger.Lgr.LogMessage("Successfully ran code over REST synchronously")

	for _, process := range finishedProcesses {
		reprErr := reporter.SendAttachment(REST_LOADER_SUBJECT, jsonString, process.Lgr.CurrentLogFile())
		if reprErr != nil {
			return reprErr
		}
	}

	logger.Lgr.LogMessage("Successfully sent REST loader results via email")

	return nil
}

// rebootHandler will handle receiving and verifying reboot commands via REST.
func (rh *RestHandler) rebootHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	rebootDelay := queryParams[REBOOT_DELAY]

	logger.Lgr.LogMessage("rebootHandler - remoteTimestamp: %v rebootDelay: %v", remoteTimestamp, rebootDelay)
	defer logger.Lgr.LogMessage("rebootHandler finished\n")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusUnauthorized, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully validated incoming timestamp")

	err = rh.verifyQueryParams(rebootDelay)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusBadRequest, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully verified query parameters")

	switch request.Method {

	case "GET":
		intDelay, intErr := strconv.Atoi(rebootDelay)
		if intErr != nil {
			logger.Lgr.LogMessage("could not convert reboot parameter to an int: %v", intErr.Error())
			rh.writeResponseAndLog(intErr.Error(), http.StatusInternalServerError, writer, request)
			return
		}

		logger.Lgr.LogMessage("sleeping for %d seconds before rebooting", intDelay)

		time.Sleep(time.Duration(intDelay) * time.Second)
		assetPath, assetErr := utils.SysAssetPath("reboot_loader.json")
		if assetErr != nil {
			logger.Lgr.LogMessage("could not successfully locate reboot loader JSON file: %v", assetErr.Error())
			rh.writeResponseAndLog(assetErr.Error(), http.StatusInternalServerError, writer, request)
			return
		}

		logger.Lgr.LogMessage("Successfully loaded reboot_loader asset: %v", assetPath)

		rebootLoader, loaderError := loader.NewLoader(assetPath)
		if loaderError != nil {
			logger.Lgr.LogMessage("could not initialize new reboot loader: %v", loaderError.Error())
			rh.writeResponseAndLog(loaderError.Error(), http.StatusInternalServerError, writer, request)
			return
		}

		logger.Lgr.LogMessage("Successfully instantiated new reboot loader: %+v", rebootLoader)

		rh.writeResponseAndLog("", http.StatusOK, writer, request)
		defer rebootLoader.StartSynchronous()

	default:
		logger.Lgr.LogMessage("Received unsupported REST method %v for rebootHandler", request.Method)
		rh.writeResponseAndLog("", http.StatusMethodNotAllowed, writer, request)
	}

	return
}

// logHandler will handle receiving and verifying log retrieval commands? via
// REST.
func (rh *RestHandler) logHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]

	logger.Lgr.LogMessage("logHandler - remoteTimestamp: %v recipientEmail: %v", remoteTimestamp, config.Cfg.CheckInGmailAddress)
	defer logger.Lgr.LogMessage("logHandler finished\n")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusUnauthorized, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully validated incoming timestamp")

	switch request.Method {
	case "GET":
		logger.Lgr.LogMessage("collating logs and sending to gmail address: %v", config.Cfg.CheckInGmailAddress)
		rh.writeResponseAndLog("", http.StatusOK, writer, request)
	case "DELETE":
		logger.Lgr.LogMessage("deleting all temp files from the local working directory to free up disk space")
		rh.writeResponseAndLog("", http.StatusOK, writer, request)
	default:
		logger.Lgr.LogMessage("Received unsupported REST method %v for logHandler", request.Method)
		rh.writeResponseAndLog("", http.StatusMethodNotAllowed, writer, request)
	}

	return
}

// updateHandler will handle receiving and verifying update commands via REST.
// Update commands will allow the remote user to force a local update given a
// specific remote URL - should probably be git for now.
func (rh *RestHandler) updateHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]

	logger.Lgr.LogMessage("updateHandler - remoteTimestamp: %v", remoteTimestamp)
	defer logger.Lgr.LogMessage("updateHandler finished\n")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusUnauthorized, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully validated incoming timestamp")

	switch request.Method {
	case "GET":
		logger.Lgr.LogMessage("need to return the current update URL")
		rh.writeResponseAndLog("", http.StatusOK, writer, request)
	case "POST":
		logger.Lgr.LogMessage("need to retrieve the URL that was posted and update config with it")
		rh.writeResponseAndLog("", http.StatusOK, writer, request)
	default:
		logger.Lgr.LogMessage("Received unsupported REST method %v for updateHandler", request.Method)
		rh.writeResponseAndLog("", http.StatusMethodNotAllowed, writer, request)
	}

	return
}

// assetHandler will allow the user to perform basic CRUD operations on files
// within the "assets" folder. The best usage of this endpoint would be to
// update the config file with new data. If the file sent over is config.json
// and the operation is an update or create then the config instance will be
// reinitialized with the new data.
func (rh *RestHandler) assetHandler(writer http.ResponseWriter, request *http.Request) {

	var err error
	queryParams := mux.Vars(request)
	remoteTimestamp := queryParams[TIMESTAMP]
	targetFileName := queryParams[ASSET_NAME]

	logger.Lgr.LogMessage("assetHandler - remoteTimestamp: %v targetFileName: %v", remoteTimestamp, targetFileName)
	defer logger.Lgr.LogMessage("assetHandler finished\n")

	err = rh.verifyTimeStamp(remoteTimestamp)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusUnauthorized, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully validated incoming timestamp")

	err = rh.verifyQueryParams(targetFileName)
	if err != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusBadRequest, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully verified query parameters")

	assetPath, assetErr := utils.AssetPath(targetFileName)
	if assetErr != nil {
		rh.writeResponseAndLog(err.Error(), http.StatusBadRequest, writer, request)
		return
	}

	logger.Lgr.LogMessage("Successfully located asset: %v", assetPath)

	switch request.Method {
	case "GET":
		logger.Lgr.LogMessage("received remote request to retrieve file: %v", targetFileName)
		rh.actionAssetAndReturn("GET", assetPath, writer, request)
	case "POST":
		logger.Lgr.LogMessage("received remote request to create new file: %v", targetFileName)
		rh.actionAssetAndReturn("POST", assetPath, writer, request)
	case "DELETE":
		logger.Lgr.LogMessage("received remote request to delete file: %v", targetFileName)
		rh.actionAssetAndReturn("DELETE", assetPath, writer, request)
	default:
		logger.Lgr.LogMessage("Received unsupported REST method %v for assetHandler", request.Method)
		rh.writeResponseAndLog("", http.StatusMethodNotAllowed, writer, request)
	}

	return
}

func (rh *RestHandler) actionAssetAndReturn(action string, assetPath string, writer http.ResponseWriter, request *http.Request) {

	switch action {

	case "GET":

		fileBytes, readErr := ioutil.ReadFile(assetPath)
		if readErr != nil {
			rh.writeResponseAndLog(fmt.Sprintf("Read error: %v from asset: %v", readErr.Error(), assetPath), http.StatusInternalServerError, writer, request)
			return
		}

		_, writeErr := writer.Write(fileBytes)
		if writeErr != nil {
			rh.writeResponseAndLog(fmt.Sprintf("Write error: %v for asset: %v", writeErr.Error(), assetPath), http.StatusInternalServerError, writer, request)
		} else {
			rh.writeResponseAndLog("", http.StatusOK, writer, request)
		}

	case "POST":

		writeBytes, readErr := ioutil.ReadAll(request.Body)
		if readErr != nil {
			rh.writeResponseAndLog(fmt.Sprintf("Read error: %v for asset: %v", readErr.Error(), assetPath), http.StatusInternalServerError, writer, request)
			return
		}

		defer request.Body.Close()

		writeErr := ioutil.WriteFile(assetPath, writeBytes, 0644)
		if writeErr != nil {
			rh.writeResponseAndLog(fmt.Sprintf("Write error: %v for asset: %v", writeErr.Error(), assetPath), http.StatusInternalServerError, writer, request)
		} else {
			rh.writeResponseAndLog("", http.StatusOK, writer, request)
		}

	case "DELETE":

		deleteErr := os.Remove(assetPath)
		if deleteErr != nil {
			rh.writeResponseAndLog(fmt.Sprintf("Delete error: %v for asset: %v", deleteErr.Error(), assetPath), http.StatusInternalServerError, writer, request)
		} else {
			rh.writeResponseAndLog("", http.StatusOK, writer, request)
		}

	default:
		rh.writeResponseAndLog(fmt.Sprintf("Unsupported asset action: %v", action), http.StatusBadRequest, writer, request)
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

	logger.Lgr.LogMessage("Calculated diff: %+v", unixDiff)

	return &unixDiff, nil
}

// verifyTimeStamp will verify the incoming timestamp from the remote machine is
// within an acceptable delta of the current time. Requires tight
// synchronization of both the local time on the local box and the remote time
// on the remote box.
func (rh *RestHandler) verifyTimeStamp(remoteTimeStamp string) error {

	logger.Lgr.LogMessage("verifyTimeStamp called with remoteTimeStamp: %v", remoteTimeStamp)

	// get the difference between then and now in seconds from unix time stamps
	diff, diffErr := rh.TimeDiffSeconds(remoteTimeStamp)
	if diffErr != nil || diff.diff >= TIMESTAMP_DELTA_SECONDS {
		return fmt.Errorf("verifyTimeStamp failed with diff: %v", diff.diff)
	}

	logger.Lgr.LogMessage("verifyTimeStamp succeeded with diff: %v", diff.diff)
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
			logger.Lgr.LogMessage("verifyQueryParams failed with: %v", value)
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
	prettyBuf.WriteString(fmt.Sprintf("now    : %d\t", utd.now))
	prettyBuf.WriteString(fmt.Sprintf("then   : %d\t", utd.then))
	prettyBuf.WriteString(fmt.Sprintf("diff   : %d\t", utd.diff))
	prettyBuf.WriteString(fmt.Sprintf("rawdiff: %d\t", utd.rawdiff))
	prettyBuf.WriteString(fmt.Sprintf("future : %t\n", utd.future))
	return prettyBuf.String()
}
