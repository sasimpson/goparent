package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/api"
)

var (
	host      string
	port      int
	user      string
	password  string
	token     string
	child     string
	date      string
	startDate string
	endDate   string
	genFlag   bool
	userData  *goparent.User
	loc       time.Location
)

const dateParseString = "2006-01-02"

type sampleError struct {
	Message string
	Origin  string
	Err     error
}

func (e sampleError) Error() string {
	return fmt.Sprintf("%s => %s", e.Origin, e.Message)
}

//NewError is the creator of the new errors
func newError(origin string, err error) error {
	return sampleError{
		Err:     err,
		Origin:  origin,
		Message: err.Error(),
	}
}

func main() {
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.IntVar(&port, "port", 8080, "port on host to connect to")
	flag.StringVar(&user, "user", "", "user to make sample data for")
	flag.StringVar(&password, "pass", "", "password for user")
	flag.StringVar(&token, "token", "", "use included token")
	flag.StringVar(&child, "child", "", "child you would like to make up data for")
	flag.StringVar(&startDate, "start", "", "date to start filling data")
	flag.StringVar(&endDate, "end", "", "date to end filling data")
	flag.StringVar(&date, "date", "", "single day for test data")
	flag.BoolVar(&genFlag, "generate", false, "generate test data")
	flag.Parse()

	fmt.Println(startDate)

	//make sure service is up:
	health := healthCheck()
	if health != nil {
		panic(health)
	}

	if token == "" {
		//get a token
		err := getToken()
		if err != nil {
			panic(err)
		}
	} else {
		err := validateToken()
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("user logged in and verfied:")
	fmt.Println(userData)
	fmt.Printf("token: %s\n", token)

	//generate in pacific timezone.  maybe make this configurable?

	if genFlag {
		switch {
		case date != "":
			// if we got a date flag, we should generate for just that date.
			genDate, err := time.Parse(dateParseString, date)
			if err != nil {
				panic(err)
			}
			fmt.Printf("generating for single date %s\n", genDate)
			err = generateRandomData(genDate)
			break
		case startDate != "" && endDate == "":
			//if we have a start date but no end, just generate from start till now.
			start, err := time.Parse(dateParseString, startDate)
			if err != nil {
				panic(err)
			}
			fmt.Printf("generating from %s\n", start)
			end := time.Now()
			for i := start; i.Before(end); i = i.AddDate(0, 0, 1) {
				err = generateRandomData(i)
				if err != nil {
					panic(err)
				}
			}
			break
		case startDate != "" && endDate != "":
			//we've got both a start and an end, generate between them.
			start, err := time.Parse(dateParseString, startDate)
			if err != nil {
				panic(err)
			}
			end, err := time.Parse(dateParseString, endDate)
			if err != nil {
				panic(err)
			}
			fmt.Printf("generating from %s to %s\n", start, end)
			for i := start; i.Before(end); i = i.AddDate(0, 0, 1) {
				err = generateRandomData(i)
				if err != nil {
					panic(err)
				}
			}
			break
		default:
			panic("no date switch")
		}
	}
}

func generateRandomData(genDate time.Time) error {
	fmt.Println("generating for ", genDate)
	var children []*goparent.Child
	var err error
	if child == "" {
		children, err = getChildren()
		if err != nil {
			return newError("generateRandomData()", err)
		}
	} else {
		c, err := getChild(child)
		if err != nil {
			return newError("generateRandomData()", err)
		}
		children = append(children, c)
	}

	fmt.Println("children seleted for generation: ")
	for _, childGen := range children {
		fmt.Println(*childGen)
		generateRandomDiaper(childGen, genDate)
	}

	return nil
}

func generateRandomDiaper(child *goparent.Child, genDate time.Time) error {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numberOfEntries := r.Intn(7) + 7
	fmt.Printf("\t\t\tNumber of Diaper Entries: %d\n", numberOfEntries)
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	loc, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	genDate = time.Date(genDate.Year(), genDate.Month(), genDate.Day(), 0, 0, 0, 0, loc)
	for x := 0; x < numberOfEntries; x++ {
		randomTime := time.Unix(genDate.Unix()+r.Int63n(86400), 0)
		diaper := goparent.Waste{
			TimeStamp: randomTime,
			ChildID:   child.ID,
			Type:      r.Intn(3) + 1,
		}
		wasteRequest := api.WasteRequest{
			WasteData: diaper,
		}

		js, err := json.Marshal(&wasteRequest)
		if err != nil {
			return newError("generateRandomDiaper()", err)
		}
		var wasteResponse api.WasteRequest
		err = makeRequest(http.MethodPost, "waste", bytes.NewReader(js), false, &wasteResponse)
		if err != nil {
			return newError("generateRandomDiaper()", err)
		}
	}
	return nil
}

func healthCheck() error {
	err := makeRequest(http.MethodGet, "info", nil, false, nil)
	if err != nil {
		return newError("healthCheck()", err)
	}
	return nil
}

func getToken() error {
	formData := url.Values{}
	formData.Add("username", user)
	formData.Add("password", password)
	var authData api.UserAuthResponse
	err := makeRequest(http.MethodPost, "user/login", strings.NewReader(formData.Encode()), true, &authData)
	if err != nil {
		return newError("getToken()", err)
	}
	token = authData.Token
	userData = authData.UserData
	return nil
}

func validateToken() error {
	var userResponse api.UserResponse
	err := makeRequest(http.MethodGet, "user/", nil, false, &userResponse)
	if err != nil {
		return newError("validateToken()", err)
	}
	userData = userResponse.UserData
	return nil
}

func getChildren() ([]*goparent.Child, error) {
	var childrenResp api.ChildrenResponse
	err := makeRequest(http.MethodGet, "children", nil, false, &childrenResp)
	if err != nil {
		return nil, newError("getChildren()", err)
	}
	return childrenResp.Children, nil
}

func getChild(id string) (*goparent.Child, error) {
	var child goparent.Child
	err := makeRequest(http.MethodGet, fmt.Sprintf("children/%s", id), nil, false, &child)
	if err != nil {
		return nil, newError("getChild()", err)
	}
	return &child, nil
}

func makeRequest(method string, path string, body io.Reader, form bool, thing interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d/api/%s", host, port, path), body)
	if err != nil {
		return newError("makeRequest()", err)
	}
	if form {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	resp, err := client.Do(req)
	if err != nil {
		return newError("makeRequest()", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/200 != 1 {
		return newError("makeRequest()", errors.New(resp.Status))
	}
	if thing != nil {
		err = json.NewDecoder(resp.Body).Decode(thing)
		if err != nil {
			return newError("makeRequest()", err)
		}
	}
	return nil
}
