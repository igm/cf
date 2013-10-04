package cf

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// Structure used to create a new application.
type NewApp struct {
	SpaceGUID string `json:"space_guid"`
	Name      string `json:"name"`
	Instances int    `json:"instances"`
	Memory    int    `json:"memory"`

	Buildpack *string `json:"buildpack"`
	Command   *string `json:"command"`
	StackGUID *string `json:"stack_guid"`
}

type appMetadata struct {
	Guid string `json:"guid"`
}

type appEntity struct {
	SpaceGUID string  `json:"space_guid"`
	Name      string  `json:"name"`
	Instances int     `json:"instances"`
	Memory    int     `json:"memory"`
	State     string  `json:"state"`
	Buildpack string  `json:"buildpack"`
	Command   string  `json:"command"`
	StackGUID string  `json:"stack_guid"`
	Routes    []Route `json:"routes"`
	Space     Space   `json:"space"`
}

type App struct {
	appMetadata `json:"metadata"`
	appEntity   `json:"entity"`
}

func (a App) String() string { return fmt.Sprintf("%s", a.Name) }

func (target *Target) AppCreate(app *NewApp) (ret *App, err error) {
	body, err := json.Marshal(app)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/v2/apps", target.TargetUrl)
	req, _ := http.NewRequest("POST", url, closeReader{bytes.NewReader(body)})
	req.Header.Set("content-type", "application/json")

	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ret)
	return
}

type Archetype struct {
	Name   string
	Reader interface {
		io.ReadCloser
		io.Seeker
	}
}

// AppPush creates a zip archive from all the provided atchetypes and upload the archive to the server
// with relevant application uuid
// TODO: implement resources diff
func (target *Target) AppPush(appGUID string, archetypes []*Archetype) (err error) {
	// hash := sha256.New()
	// io.Copy(hash, file)
	// file.Seek(0, 0)
	// md := hash.Sum(nil)
	// mdStr := hex.EncodeToString(md)

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	for _, archetype := range archetypes {
		f, err := zw.Create(archetype.Name)
		if err != nil {
			return err
		}
		io.Copy(f, archetype.Reader)
		archetype.Reader.Close()
	}
	err = zw.Close()
	if err != nil {
		return
	}
	target.AppPushArchive(appGUID, buf)
	return
}

func (target *Target) AppPushArchive(appGUID string, reader io.Reader) (err error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	boundary := writer.Boundary()

	part, err := writer.CreateFormField("resources")
	if err != nil {
		return
	}
	_, err = io.Copy(part, bytes.NewBufferString("[]"))
	if err != nil {
		return
	}
	part, err = writer.CreateFormFile("application", "application.zip")
	if err != nil {
		return
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return
	}
	writer.Close()

	url := fmt.Sprintf("%s/v2/apps/%s/bits", target.TargetUrl, appGUID)
	req, _ := http.NewRequest("PUT", url, body)
	req.Header.Set("content-type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))

	_, err = target.sendRequest(req)
	if err != nil {
		return
	}

	return nil
}

func (target *Target) AppStart(appGUID string) (err error) { return target.appState(appGUID, "STARTED") }
func (target *Target) AppStop(appGUID string) (err error)  { return target.appState(appGUID, "STOPPED") }
func (target *Target) appState(appGUID string, state string) (err error) {
	body, err := json.Marshal(struct {
		Console bool   `json:"console"`
		State   string `json:"state"`
	}{true, state})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/v2/apps/%s", target.TargetUrl, appGUID)
	req, _ := http.NewRequest("PUT", url, closeReader{bytes.NewReader(body)})
	req.Header.Set("content-type", "application/json")
	_, err = target.sendRequest(req)
	return
}

func (target *Target) AppDelete(appGUID string) (err error) {
	url := fmt.Sprintf("%s/v2/apps/%s", target.TargetUrl, appGUID)
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err = target.sendRequest(req)
	return
}

func (target *Target) AppsGet() (apps []App, err error) {
	url := fmt.Sprintf("%s/v2/apps?inline-relations-depth=2", target.TargetUrl)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var response struct {
		Apps []App `json:"resources"`
	}
	err = decoder.Decode(&response)
	apps = response.Apps
	return
}

func (target *Target) AppAddRoute(appGUID, routeGUID string) (err error) {
	url := fmt.Sprintf("%s/v2/apps/%s/routes/%s", target.TargetUrl, appGUID, routeGUID)
	req, _ := http.NewRequest("PUT", url, nil)
	_, err = target.sendRequest(req)
	return
}

func (target *Target) AppDeleteRoute(appGUID, routeGUID string) (err error) {
	url := fmt.Sprintf("%s/v2/apps/%s/routes/%s", target.TargetUrl, appGUID, routeGUID)
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err = target.sendRequest(req)
	return
}

func (target *Target) AppRoutesGet(appGUID string) (routes []Route, err error) {
	url := fmt.Sprintf("%s/v2/apps/%s/routes?inline-relations-depth=1", target.TargetUrl, appGUID)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var response struct {
		Routes []Route `json:"resources"`
	}
	err = decoder.Decode(&response)
	routes = response.Routes
	return
}

type Instance struct {
	State       string  `json:"state"`
	Since       float64 `json:"since"`
	DebugIP     string  `json:"debug_ip"`
	DebugPort   string  `json:"debug_port"`
	ConsoleIP   string  `json:"console_ip"`
	ConsolePort string  `json:"console_port"`
}

func (target *Target) AppInstances(appGUID string) (instances map[string]Instance, err error) {
	url := fmt.Sprintf("%s/v2/apps/%s/instances", target.TargetUrl, appGUID)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&instances)
	return
}
