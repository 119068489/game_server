//Package easygo ...
package easygo

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

//电竞服开启
var (
	rePathUpload = regexp.MustCompile(`^/upload$`)
	rePathFiles  = regexp.MustCompile(`^/files/([^/]+)$`)

	errTokenMismatch = errors.New("token mismatched")
	errMissingToken  = errors.New("missing token")

	protectedMethods = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut}
)

type response struct {
	OK bool `json:"ok"`
}

type uploadedResponse struct {
	response
	Path string `json:"path"`
}

func newUploadedResponse(path string) uploadedResponse {
	return uploadedResponse{response: response{OK: true}, Path: path}
}

type errorResponse struct {
	response
	Message string `json:"error"`
}

func newErrorResponse(err error) errorResponse {
	return errorResponse{response: response{OK: false}, Message: err.Error()}
}

func writeError(w http.ResponseWriter, err error) (int, error) {
	body := newErrorResponse(err)
	b, e := json.Marshal(body)
	// if an error is occured on marshaling, write empty value as response.
	if e != nil {
		return w.Write([]byte{})
	}
	return w.Write(b)
}

func writeSuccess(w http.ResponseWriter, path string) (int, error) {
	body := newUploadedResponse(path)
	b, e := json.Marshal(body)
	// if an error is occured on marshaling, write empty value as response.
	if e != nil {
		return w.Write([]byte{})
	}
	return w.Write(b)
}

func getSize(content io.Seeker) (int64, error) {
	size, err := content.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	_, err = content.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// Server represents a simple-upload server.
type Server struct {
	DocumentRoot string
	// MaxUploadSize limits the size of the uploaded content, specified with "byte".
	MaxUploadSize int64
	SecureToken   string
	EnableCORS    bool
}

// NewServer creates a new simple-upload server.
func NewServer(documentRoot string, maxUploadSize int64, token string, enableCORS bool) Server {
	return Server{
		DocumentRoot:  documentRoot,
		MaxUploadSize: maxUploadSize,
		SecureToken:   token,
		EnableCORS:    enableCORS,
	}
}

//=====================================================================================api
func (s Server) handleGet(w http.ResponseWriter, r *http.Request) {
	if !rePathFiles.MatchString(r.URL.Path) {
		w.WriteHeader(http.StatusNotFound)
		writeError(w, fmt.Errorf("\"%s\" is not found", r.URL.Path))
		return
	}
	if s.EnableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	http.StripPrefix("/files/", http.FileServer(http.Dir(s.DocumentRoot))).ServeHTTP(w, r)
}

func (s Server) HandlePost(w http.ResponseWriter, r *http.Request) string {
	srcFile, info, err := r.FormFile("file")
	if err != nil {
		logs.Error("failed to acquire the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return ""
	}
	defer srcFile.Close()

	// logs.Debug(info)
	size, err := getSize(srcFile)
	if err != nil {
		logs.Error("failed to get the size of the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return ""
	}
	if size > s.MaxUploadSize {
		logs.Error("size", size, "file size exceeded")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		writeError(w, errors.New("uploaded file size exceeds the limit"))
		return ""
	}

	body, err := ioutil.ReadAll(srcFile)
	if err != nil {
		logs.Error("failed to read the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return ""
	}
	filename := info.Filename
	if filename == "" {
		filename = fmt.Sprintf("%x", sha1.Sum(body))
	}

	if s.DocumentRoot != "" {
		os.Mkdir("./"+s.DocumentRoot, os.ModePerm)
	}

	dstPath := path.Join(s.DocumentRoot, filename)

	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logs.Error("path", dstPath, "failed to read the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return ""
	}
	defer dstFile.Close()
	if written, err := dstFile.Write(body); err != nil {
		logs.Error("path", dstPath, "failed to write the content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return ""
	} else if int64(written) != size {
		logs.Error("size", size, "written", written, "failed to read the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, fmt.Errorf("the size of uploaded content is %d, but %d bytes written", size, written))
	}
	// logs.Info("path", dstPath, "url", uploadedURL, "size", size, "file uploaded by POST")
	if s.EnableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.WriteHeader(http.StatusOK)
	writeSuccess(w, dstPath)

	return dstPath
}

func (s Server) handlePut(w http.ResponseWriter, r *http.Request) {
	matches := rePathFiles.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		logs.Error("path", r.URL.Path, "invalid path")
		w.WriteHeader(http.StatusNotFound)
		writeError(w, fmt.Errorf("\"%s\" is not found", r.URL.Path))
		return
	}
	targetPath := path.Join(s.DocumentRoot, matches[1])

	// We have to create a new temporary file in the same device to avoid "invalid cross-device link" on renaming.
	// Here is the easiest solution: create it in the same directory.
	tempFile, err := ioutil.TempFile(s.DocumentRoot, "upload_")
	if err != nil {
		logs.Error("failed to create a temporary file")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}
	defer r.Body.Close()
	srcFile, info, err := r.FormFile("file")
	if err != nil {
		logs.Error("path", targetPath, "failed to acquire the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}
	defer srcFile.Close()
	// dump headers for the file
	logs.Debug(info.Header)

	size, err := getSize(srcFile)
	if err != nil {
		logs.Error("path", targetPath, "failed to get the size of the uploaded content")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}
	if size > s.MaxUploadSize {
		logs.Error("path", targetPath, "size", size, "file size exceeded")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		writeError(w, errors.New("uploaded file size exceeds the limit"))
		return
	}

	n, err := io.Copy(tempFile, srcFile)
	if err != nil {
		logs.Error("path", tempFile.Name(), "size", size, "failed to write body to the file")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}
	// excplicitly close file to flush, then rename from temp name to actual name in atomic file
	// operation if on linux or other unix-like OS (windows hosts should look into https://github.com/natefinch/atomic
	// package for atomic file write operations)
	tempFile.Close()
	if err := os.Rename(tempFile.Name(), targetPath); err != nil {
		os.Remove(tempFile.Name())
		logs.Error("path", targetPath, "failed to rename temp file to final filename for upload")
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, err)
		return
	}
	logs.Error("path", r.URL.Path, "size", n, "file uploaded by PUT")
	if s.EnableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.WriteHeader(http.StatusOK)
	writeSuccess(w, r.URL.Path)
}

func (s Server) handleOptions(w http.ResponseWriter, r *http.Request) {
	var allowedMethods []string
	if rePathFiles.MatchString(r.URL.Path) {
		allowedMethods = []string{http.MethodPut, http.MethodGet, http.MethodHead}
	} else if rePathUpload.MatchString(r.URL.Path) {
		allowedMethods = []string{http.MethodPost}
	} else {
		w.WriteHeader(http.StatusNotFound)
		writeError(w, errors.New("not found"))
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
	w.WriteHeader(http.StatusNoContent)
}

func (s Server) checkToken(r *http.Request) error {
	// first, try to get the token from the query strings
	token := r.URL.Query().Get("token")
	// if token is not found, check the form parameter.
	if token == "" {
		token = r.FormValue("token")
	}
	if token == "" {
		return errMissingToken
	}
	if token != s.SecureToken {
		return errTokenMismatch
	}
	return nil
}

func isAuthenticationRequired(r *http.Request) bool {
	for _, m := range protectedMethods {
		if m == r.Method {
			return true
		}
	}
	return false
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.checkToken(r); isAuthenticationRequired(r) && err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		s.handleGet(w, r)
	case http.MethodPost:
		s.HandlePost(w, r)
	case http.MethodPut:
		s.handlePut(w, r)
	case http.MethodOptions:
		s.handleOptions(w, r)
	default:
		w.Header().Add("Allow", "GET,HEAD,POST,PUT")
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeError(w, fmt.Errorf("method \"%s\" is not allowed", r.Method))
	}
}

func run(args []string) int {
	bindAddress := flag.String("ip", "0.0.0.0", "IP address to bind")
	listenPort := flag.Int("port", 8001, "port number to listen on")
	tlsListenPort := flag.Int("tlsport", 25443, "port number to listen on with TLS")
	// 5,242,880 bytes == 5 MiB
	maxUploadSize := flag.Int64("upload_limit", 5242880, "max size of uploaded file (byte)")
	tokenFlag := flag.String("token", "", "specify the security token (it is automatically generated if empty)")
	certFile := flag.String("cert", "", "path to certificate file")
	keyFile := flag.String("key", "", "path to key file")
	corsEnabled := flag.Bool("cors", false, "if true, add ACAO header to support CORS")
	flag.Parse()
	serverRoot := flag.Arg(0)
	if len(serverRoot) == 0 {
		flag.Usage()
		return 2
	}

	token := *tokenFlag
	if token == "" {
		count := 10
		b := make([]byte, count)
		if _, err := rand.Read(b); err != nil {
			logs.Error("could not generate token")
			return 1
		}
		token = fmt.Sprintf("%x", b)
		logs.Error("token", token, "token generated")
	}
	tlsEnabled := *certFile != "" && *keyFile != ""
	server := NewServer(serverRoot, *maxUploadSize, token, *corsEnabled)
	http.Handle("/upload", server)
	http.Handle("/files/", server)

	errors := make(chan error)

	go func() {

		logs.Info("ip", *bindAddress,
			"port", *listenPort,
			"token", token,
			"upload_limit", *maxUploadSize,
			"root", serverRoot,
			"cors", *corsEnabled, "start listening")
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *bindAddress, *listenPort), nil); err != nil {
			errors <- err
		}
	}()

	if tlsEnabled {
		go func() {
			logs.Info("cert", *certFile,
				"key", *keyFile,
				"port", *tlsListenPort, "start listening TLS")
			if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *bindAddress, *tlsListenPort), *certFile, *keyFile, nil); err != nil {
				errors <- err
			}
		}()
	}

	err := <-errors
	logs.Info(err, "closing server")
	return 0
}

//===================================================================================rpc
func UploadFile(filename string, file []byte, path ...string) (pathfile string, filepath string, err error) {
	maxUploadSize := NewInt64(5242880)
	path = append(path, "upload")
	tokenFlag := NewString("")
	corsEnabled := NewBool(false)
	flag.Parse()

	token := *tokenFlag
	if token == "" {
		count := 10
		b := make([]byte, count)
		if _, err := rand.Read(b); err != nil {
			logs.Error("could not generate token")
		}
		token = fmt.Sprintf("%x", b)
	}

	server := NewServer(path[0], *maxUploadSize, token, *corsEnabled)
	pathfile, err = server.rpcUploadFile(filename, file)
	return pathfile, path[0], err
}

func (s Server) rpcUploadFile(filename string, file []byte) (string, error) {

	os.Mkdir("./"+s.DocumentRoot, os.ModePerm)
	dstPath := path.Join(s.DocumentRoot, filename)
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logs.Error("path", dstPath, "failed to read the uploaded content")
		return "", err
	}
	defer dstFile.Close()
	if _, err = dstFile.Write(file); err != nil {
		logs.Error("path", dstPath, "failed to write the content")
		return "", err
	}

	return dstPath, nil
}
