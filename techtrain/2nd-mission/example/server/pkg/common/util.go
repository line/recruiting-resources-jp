// Package common provides utility functions
package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// HasElement checks if string `a` in stirng slice `s`.
// Returns true if it does.
func HasElement(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// HasAnyElement checks if at least one element of string
// slice `l` in string slice `s`.
// Returns true if it does.
func HasAnyElement(s []string, l []string) bool {
	for _, v := range l {
		if HasElement(s, v) {
			return true
		}
	}
	return false
}

// Unique removes duplicate elements in s.
func Unique(s []string) []string {
	unique := []string{}
	for _, v := range s {
		if !HasElement(unique, v) {
			unique = append(unique, v)
		}
	}
	return unique
}

// GetPathWithoutVersion removes api version from request path.
func GetPathWithoutVersion(path string) string {
	var origin, tmp []string
	origin = strings.Split(path, "/")
	for _, e := range origin {
		if e != "" {
			tmp = append(tmp, e)
		}
	}
	return "/" + strings.Join(tmp[1:], "/")
}

// GetRequestStr returns "method path" of a http request.
func GetRequestStr(r *http.Request) string {
	return r.Method + " " + r.URL.Path
}

// GetRequestBody returns a copy of body.
// It's used for reading the body more than once.
func GetRequestBody(r *http.Request) io.Reader {
	b := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(r.Body, b)
	// TODO make sure no mem leak
	r.Body = ioutil.NopCloser(b)
	return reader
}

// GetReponseBody returns a copy of body.
// It's used for reading the body more than once.
func GetResponseBody(r *http.Response) io.Reader {
	b := bytes.NewBuffer(make([]byte, 0))
	reader := io.TeeReader(r.Body, b)
	// TODO make sure no mem leak
	r.Body = ioutil.NopCloser(b)
	return reader
}

// RespondJSON makes http response with json payload.
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	rep, err := json.Marshal(payload)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "payload": payload}).
			Error("Failed to json marshal")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(rep))
	}
}

// RespondErr makes error response with format {"errMsg": "xxx"}.
func RespondErr(w http.ResponseWriter, status int, msg string) {
	RespondJSON(w, status, map[string]string{"errMsg": msg})
}

// RespondErr makes error response with format {"errMsg": "xxx"}.
func RespondErrf(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	RespondErr(w, status, fmt.Sprintf(msg, args...))
}

type Log struct {
	RequestId string
}

func (l *Log) Configure() {
	formatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FullTimestamp:   true,
	}
	log.SetFormatter(formatter)
}

func (l *Log) Info(fields log.Fields, msg string, args ...interface{}) {
	fields["requestId"] = l.RequestId
	log.WithFields(fields).Infof(msg, args...)
}

func (l *Log) Debug(fields log.Fields, msg string, args ...interface{}) {
	fields["requestId"] = l.RequestId
	log.WithFields(fields).Debugf(msg, args...)
}

func (l *Log) Warn(fields log.Fields, msg string, args ...interface{}) {
	fields["requestId"] = l.RequestId
	log.WithFields(fields).Debugf(msg, args...)
}

func (l *Log) Error(fields log.Fields, msg string, args ...interface{}) {
	fields["requestId"] = l.RequestId
	log.WithFields(fields).Errorf(msg, args...)
}
