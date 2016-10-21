package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seantcanavan/logger"
)

const TIMESTAMP_PARAM = "timestamp"

type RestHandler struct {
	router *mux.Router
	logger *logger.SeanLogger
}

func NewRestHandler(seanLogger *logger.SeanLogger) *RestHandler {
	r := RestHandler{}
	r.logger = seanLogger
	r.router = mux.NewRouter()
	r.router.HandleFunc("/", CheckinHandler)
	r.router.HandleFunc("/execute/{"+TIMESTAMP_PARAM+"}/{fileurl}", r.ExecuteHandler)
	r.router.HandleFunc("/reboot/{"+TIMESTAMP_PARAM+"}/{delay}", r.RebootHandler)
	r.router.HandleFunc("/sendlogs/{"+TIMESTAMP_PARAM+"}/{emailaddress}", r.LogHandler)
	r.router.HandleFunc("/forceupdate/{"+TIMESTAMP_PARAM+"}/{remoteurl}", r.UpdateHandler)
	r.router.HandleFunc("/updateconfig/{"+TIMESTAMP_PARAM+"}", r.ConfigHandler)

	return &r
}

func Setup() {

}

func CheckinHandler(writer http.ResponseWriter, request *http.Request) {

}

func (rh *RestHandler) ExecuteHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("ExecuteHandler started")
	defer rh.logger.LogMessage("ExecuteHandler finished")

	queryParams := mux.Vars(request)
	if err := rh.verifyTimeStamp(queryParams[TIMESTAMP_PARAM]); err != nil {

	}
}

func (rh *RestHandler) RebootHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("RebootHandler started")
	defer rh.logger.LogMessage("RebootHandler finished")

	queryParams := mux.Vars(request)
	if err := rh.verifyTimeStamp(queryParams[TIMESTAMP_PARAM]); err != nil {

	}
}

func (rh *RestHandler) LogHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("LogHandler started")
	defer rh.logger.LogMessage("LogHandler finished")

	queryParams := mux.Vars(request)
	if err := rh.verifyTimeStamp(queryParams[TIMESTAMP_PARAM]); err != nil {

	}
}

func (rh *RestHandler) UpdateHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("UpdateHandler started")
	defer rh.logger.LogMessage("UpdateHandler finished")

	queryParams := mux.Vars(request)
	if err := rh.verifyTimeStamp(queryParams[TIMESTAMP_PARAM]); err != nil {

	}
}

func (rh *RestHandler) ConfigHandler(writer http.ResponseWriter, request *http.Request) {
	rh.logger.LogMessage("ConfigHandler started")
	defer rh.logger.LogMessage("ConfigHandler finished")

	queryParams := mux.Vars(request)
	if err := rh.verifyTimeStamp(queryParams[TIMESTAMP_PARAM]); err != nil {

	}
}

func  (rh *RestHandler) verifyTimeStamp(incomingTimeStamp string) error {
	return nil
}
