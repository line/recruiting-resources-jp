package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"todo-example/server/pkg/common"
	"todo-example/server/pkg/config"
)

// CheckSchema validates user request based on api schema.
func CheckSchema() Decorator {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := r.Context().Value("logger").(*common.Log)
			logger.Info(log.Fields{}, "Start API schema filter")
			/* apiSpec is a "go-openapi/spec/swagger.go *Swagger" object,
			   which represents the swagger api spec.
			   type SwaggerProps struct {
					...
					BasePath	string
					Paths       *Paths
					...
			   }
			*/
			conf := config.GetConf()
			apiSpec := conf.Schema.Spec()
			logger.Debug(log.Fields{}, "Start to validate user request")
			err := validateRequest(apiSpec, r, func(path, method string, op *spec.Operation) error {
				return validateParams(op.Parameters, r)
			})
			if err == nil {
				logger.Info(log.Fields{}, "Passed API schema filter")
				h.ServeHTTP(w, r)
			} else {
				logger.Info(log.Fields{"error": err.Error()}, "Failed API schema filter")
				common.RespondErr(w, 400, err.Error())
			}
		})
	}
}

type validateFunc func(path, method string, op *spec.Operation) error

// validate validates every path of api spec.
// it calls a validate function to check parameters of each path.
func validateRequest(spec *spec.Swagger, r *http.Request, fn validateFunc) error {
	/* each path is a "go-openapi/spec/path_item.go PathItem" object,
	   which represents a path in swagger spec.
	   type PathItem struct {
		   Get  *Operation
		   ...
		   Parameters []Parameter
	   }
	*/
	logger := r.Context().Value("logger").(*common.Log)
	route, _ := mux.CurrentRoute(r).GetPathTemplate()
	for path, operations := range spec.Paths.Paths {
		for method, operation := range getOperationsWithMethod(&operations) {
			if operation == nil {
				continue
			}
			if method == r.Method && path == route {
				logger.Debug(log.Fields{"route": route}, "matched api schema path")
				return fn(path, method, operation)
			}
		}
	}
	return errors.New("No such API entry!")

}

// it validates parameters of user request.
func validateParams(params []spec.Parameter, r *http.Request) error {
	logger := r.Context().Value("logger").(*common.Log)
	for _, param := range params {
		/* param is a "go-openapi/spec/parameter.go Parameter" object,
		   which represents a parameter inside of a method in swagger api spec.
		   type Parameter struct {
			   ...
			   In		string
			   Schema   *Schema
			   ...
		   }
		*/
		switch param.In {
		case "body":
			logger.Debug(log.Fields{"paramIn": "body"}, "Start to check request body")
			var requestBody interface{}
			defer r.Body.Close()
			reader := common.GetRequestBody(r)
			err := json.NewDecoder(reader).Decode(&requestBody)
			if err != nil {
				return err
			}
			logger.Debug(log.Fields{"body": requestBody}, "Start to validate body")
			return validate.AgainstSchema(param.Schema, requestBody, strfmt.Default)
		case "query":
			queryParams := r.URL.Query()
			logger.Debug(log.Fields{"paramName": param.Name, "paramIn": "query", "query": queryParams}, "Start to check query parameters")
			_, ok := queryParams[param.Name]
			if !ok {
				if param.Required {

					return errors.New(param.Name + " is required!")
				}
				if param.Default != nil {
					value := fmt.Sprintf("%v", param.Default)
					logger.Debug(log.Fields{"param": param.Name, "value": value}, "using default parameter")
					queryParams.Set(param.Name, value)
					r.URL.RawQuery = queryParams.Encode()
				}
			} else {
				logger.Debug(log.Fields{"paramType": param.Type}, "Start to convert param type")
				value, err := convert(queryParams[param.Name], param.Type)
				if err != nil {
					return errors.New(param.Name + " " + err.Error())
				}
				logger.Debug(log.Fields{}, "Start to validate query param")
				if result := validate.NewParamValidator(&param, strfmt.Default).Validate(value); result != nil {
					msg := ""
					for _, e := range result.Errors {
						msg = msg + e.Error()
					}
					return errors.New(msg)
				}
			}
		}
	}
	return nil
}

/* getOperations returns operations in the path with http method.
   operation is a "go-openapi/spec/operation.go Operation" object.
   type Operation struct {
	   ...
	   Schemes      []string
	   Parameters   []Parameter
	   ...
   }
*/
func getOperationsWithMethod(props *spec.PathItem) map[string]*spec.Operation {
	ops := map[string]*spec.Operation{
		"DELETE":  props.Delete,
		"GET":     props.Get,
		"HEAD":    props.Head,
		"OPTIONS": props.Options,
		"PATCH":   props.Patch,
		"POST":    props.Post,
		"PUT":     props.Put,
	}
	return ops
}

func convert(values []string, valueType string) (interface{}, error) {
	num := len(values)
	if valueType != "array" && num != 1 {
		return values, errors.New("Multiple values detected, should be one!")
	}
	value := values[0]
	switch valueType {
	case "integer":
		i, err := strconv.Atoi(value)
		if err != nil {
			return value, errors.New(value + " should be integer type!")
		} else {
			return i, nil
		}
	// TODO process other types
	default:
		return value, nil
	}
}
