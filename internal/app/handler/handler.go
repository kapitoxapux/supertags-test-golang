package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/idoubi/goz"

	"supertags/internal/app/config"
	"supertags/internal/app/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {

	return &Handler{
		service: service,
	}
}

type ArangoEvent struct {
	MESSAGES Message    `json:"MESSAGES"`
	DATA     ArangoData `json:"DATA"`
	STATUS   string     `json:"STATUS"`
}

type Message struct {
	Error   string         `json:"error"`
	Warning string         `json:"warning"`
	Info    map[int]string `json:"info"`
}

type ArangoData struct {
	Page        uint32        `json:"page"`
	Pages_count uint32        `json:"pages_count"`
	Rows_count  uint32        `json:"rows_count"`
	Rows        []service.Row `json:"rows"`
}

type RawFilter struct {
	Filter FilterFieldStruct `json:"filter"`
	Sort   SortFields        `json:"sort"`
	Limit  uint32            `json:"limit"`
}

type FilterFieldStruct struct {
	Field FieldStruct `json:"field"`
}

type FieldStruct struct {
	Key    string   `json:"key"`
	Sign   string   `json:"sign"`
	Values []string `json:"values"`
}

type SortFields struct {
	Fields    []string `json:"fields"`
	Direction string   `json:"direction"`
}

func CustomMiddleware(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// log.Println("in here")
			// 			cookie, err := r.Cookie("SESSION")
			// 			if cookie == nil {
			// 				http.Error(w, err.Error(), http.StatusUnauthorized)

			// 				return
			// 			}

			h.ServeHTTP(w, r)
		},
	)
}

func ConnectionDBCheck() (int, string) {
	db, err := sql.Open("mysql", config.GetDB())
	if err != nil {

		return 500, err.Error()
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {

		return 500, err.Error()
	}

	return 200, "Connected!"
}

func SetUserCookie(req *http.Request, data string) *http.Cookie {
	expiration := time.Now().Add(6000 * time.Second)

	return &http.Cookie{
		Name:    "SESSION",
		Value:   data,
		Path:    "/",
		Expires: expiration,
	}
}

func (h *Handler) GetAuthentication(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed for this route", http.StatusMethodNotAllowed)

		return
	}

	cli := goz.NewClient()

	resp, err := cli.Get("https://development.kpi-drive.ru/_api/auth/login", goz.Options{
		Query: map[string]interface{}{
			"login":    "admin",
			"password": "admin",
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	resp_cookie := resp.GetHeader("Set-Cookie")

	cookies := strings.Split(resp_cookie[0], ";")
	final_row_cookie := strings.Split(cookies[0], "=")
	http.SetCookie(res, SetUserCookie(req, final_row_cookie[1]))

	config.SetCookieConfig(final_row_cookie[1])

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}

func (h *Handler) GetArangoData(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed for this route!", http.StatusMethodNotAllowed)

		return
	}

	fmt.Println("Получение данных из ArangoDB:")
	fmt.Print("\n")
	raw := RawFilter{
		Filter: FilterFieldStruct{
			Field: FieldStruct{
				Key:    "type",
				Sign:   "LIKE",
				Values: []string{"MATRIX_REQUEST"},
			},
		},
		Sort: SortFields{
			Fields:    []string{"time"},
			Direction: "DESC",
		},
		Limit: 10,
	}

	body, _ := json.Marshal(raw)

	r, err := http.NewRequest(http.MethodGet, "https://development.kpi-drive.ru/_api/events", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	r.Header.Set("Content-Type", "application/json")
	r.AddCookie(req.Cookies()[0])

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(body))
		fmt.Print("\n")

		argo := new(ArangoEvent)
		json.Unmarshal(body, &argo)

		fmt.Println("Обработка значений ArangoDB и запись в MySQL")
		fmt.Print("\n")

		results := []string{}
		var content_type string

		for i, row := range argo.DATA.Rows {
			t, _ := time.Parse(time.RFC3339, row.Time)

			fact := new(service.Fact)
			supertags := []service.Supertag{}
			comments := []service.Row{}

			supertag := new(service.Supertag)
			tag := new(service.Tag)

			tag.Id = 2
			tag.Name = "Клиент"
			tag.Key = "client"
			tag.ValuesSource = 0

			supertag.Tag = *tag
			supertag.Value = row.Authors.UserName

			supertags = append(supertags, *supertag)
			comments = append(comments, row)

			fact.PeriodStart = "2023-09-01"
			fact.PeriodEnd = "2023-09-30"
			fact.PeriodKey = "month"
			fact.IndicatorToMoId = 315914
			fact.IndicatorToFactId = 0
			fact.Value = 1
			fact.FactTime = t.Format("2006-01-02")
			fact.IsPlan = 0
			fact.Supertags = supertags
			fact.AuthUserId = 40
			fact.Comments = comments

			fact_in_byte, _ := json.Marshal(fact)
			fmt.Println(fmt.Sprintf("%v", i+1) + " row to store:")
			fmt.Print("\n")
			fmt.Println(string(fact_in_byte))
			fmt.Print("\n")
			cli := &http.Client{
				Timeout: time.Second * 10,
			}

			st, _ := json.Marshal(fact.Supertags)
			comment, _ := json.Marshal(fact.Comments)

			form_data := map[string]string{
				"period_start":            fact.PeriodStart,
				"period_end":              fact.PeriodEnd,
				"period_key":              fact.PeriodKey,
				"indicator_to_mo_id":      fmt.Sprintf("%v", fact.IndicatorToMoId),
				"indicator_to_mo_fact_id": fmt.Sprintf("%v", fact.IndicatorToFactId),
				"value":                   fmt.Sprintf("%v", fact.Value),
				"fact_time":               fact.FactTime,
				"is_plan":                 fmt.Sprintf("%v", fact.IsPlan),
				"supertags":               string(st),
				"auth_user_id":            fmt.Sprintf("%v", fact.AuthUserId),
				"comment":                 string(comment),
			}

			ct, form, err := createForm(form_data)
			content_type = ct
			if err != nil {
				panic(err)
			}

			w, e := http.NewRequest("POST", "https://development.kpi-drive.ru/_api/facts/save_fact", form)
			if e != nil {
				panic(e)
			}

			w.Header.Set("Content-Type", ct)
			w.AddCookie(req.Cookies()[0])

			rsp, _ := cli.Do(w)
			if rsp.StatusCode != http.StatusOK {
				log.Printf("Request failed with response code: %d", rsp.StatusCode)
			}

			form_body, _ := io.ReadAll(rsp.Body)

			results = append(results, string(form_body))
			fmt.Println("response from MySQL", string(form_body))
			fmt.Print("\n")

			// h.channel.InputChannel <- fact
		}

		res.Header().Set("Content-Type", content_type)
		res.Header().Add("Accept", "application/json")

		stringByte := strings.Join(results, "\n")

		res.Write([]byte(stringByte))

	} else {
		fmt.Println("Get failed with error: ", resp.Status)
	}
}

func createForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)

	defer mp.Close()

	for key, val := range form {
		mp.WriteField(key, val)
	}

	return mp.FormDataContentType(), body, nil
}
