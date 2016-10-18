package rest

import (
    "net/http"

    "github.com/seantcanavan/logger"
)

const TIMESTAMP_PARAM = "timestamp"

func Setup() {
    router := mux.NewRouter()
    router.HandleFunc("/", CheckinHandler)
    router.HandleFunc("/execute/{" + TIMESTAMP_PARAM + "}/{fileurl}", ExecuteHandler)
    router.HandleFunc("/reboot/{" + TIMESTAMP_PARAM + "}/{delay}", RebootHandler}
    router.HandleFunc("/sendlogs/{" + TIMESTAMP_PARAM + "}/{emailaddress}", LogHandler}
    router.HandleFunc("/forceupdate/{" + TIMESTAMP_PARAM + "}/{remoteurl}", UpdateHandler}
    router.HandleFunc("/updateconfig/{" + TIMESTAMP_PARAM + "}", ConfigHandler}
}

func ExecuteHandler(writer http.ResponseWriter, request *http.Request) {
    logger.LogMessage("ExecuteHandler started")
    queryParams := mux.Vars(request)
    if err := verifyTimeStamp(queryParams[TIMESTAMP_PARAM]); err != nil {

    }

    defer logger.LogMessage("ExecuteHandler finished")
}

func RebootHandler(writer http.ResponseWriter, request *http.Request) {
    logger.LogMessage("RebootHandler started")
    verifyTimeStamp(queryParams[TIMESTAMP_PARAM])
    defer logger.LogMessage("RebootHandler finished")
}

func LogHandler(writer http.ResponseWriter, request *http.Request) {
    logger.LogMessage("LogHandler started")
    verifyTimeStamp(queryParams[TIMESTAMP_PARAM])
    defer logger.LogMessage("LogHandler finished")
}

func UpdateHandler(writer http.ResponseWriter, request *http.Request) {
    logger.LogMessage("UpdateHandler started")
    verifyTimeStamp(queryParams[TIMESTAMP_PARAM])
    defer logger.LogMessage("UpdateHandler finished")
}

func ConfigHandler(writer http.ResponseWriter, request *http.Request) {
    logger.LogMessage("ConfigHandler started")
    defer logger.LogMessage("ConfigHandler finished")
}

func verifyTimeStamp(incomingTimeStamp string) error {

}
