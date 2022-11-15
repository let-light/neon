package rtsp

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Url     string
	Version string
	Lines   HeaderLines
	Content []byte
}

func UnmarshalRequest(buf []byte) (*Request, int, error) {
	headerEndOffset := bytes.Index(buf, []byte("\r\n\r\n"))
	if headerEndOffset == -1 {
		return nil, -1, errors.New("incomplete packet")
	}

	endOffset := headerEndOffset + 4

	contentLength := 0

	req := &Request{
		Lines: make(HeaderLines),
	}

	lines := bytes.Split(buf[:headerEndOffset], []byte("\r\n"))
	if len(lines) < 2 {
		return nil, endOffset, errors.New("invalid packet")
	}

	// parse first line
	methodLine := lines[0]
	methodLineParts := bytes.Split(methodLine, []byte(" "))
	if len(methodLineParts) != 3 {
		return nil, endOffset, errors.New("invalid packet")
	}

	req.Method = strings.ToLower(string(methodLineParts[0]))
	req.Url = strings.ToLower(string(methodLineParts[1]))
	req.Version = strings.ToLower(string(methodLineParts[2]))

	// parse other lines
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		idx := bytes.Index(line, []byte(":"))
		if idx == -1 {
			continue
		}

		key := strings.ToLower(string(line[:idx]))
		value := string(line[idx+1:])
		req.Lines[key] = value

		if key == "content-length" {
			var err error
			contentLength, err = strconv.Atoi(value)
			if err != nil {
				return nil, headerEndOffset + 4, err
			}
		}
	}

	req.Content = buf[headerEndOffset+4 : headerEndOffset+4+contentLength]

	endOffset += contentLength

	return req, endOffset, nil
}

func (lines HeaderLines) String() string {
	s := ""
	for k, v := range lines {
		s += k + ": " + v + "\r\n"
	}
	return s
}

func (req *Request) GetLine(key string) string {
	return req.Lines[key]
}

func (req *Request) SetLine(key, value string) {
	req.Lines[key] = value
}

func (req *Request) ToString() string {
	return req.Method + " " + req.Url + " " + req.Version + "\r\n" + req.Lines.String() + "\r\n" + string(req.Content)
}

func (req *Request) CSeq() int {
	cseqLine := req.Lines["cseq"]
	if cseqLine == "" {
		return -1
	}

	cseq, err := strconv.Atoi(cseqLine)
	if err != nil {
		return -1
	}

	return cseq
}

func (req *Request) Session() string {
	return req.Lines["session"]
}

func (req *Request) Option() *OptionsRequest {
	return &OptionsRequest{
		Request: req,
	}
}

func (req *Request) Describe() *DescribeRequest {
	return &DescribeRequest{
		Request: req,
	}
}

func (req *Request) Announce() *AnnounceRequest {
	return &AnnounceRequest{
		Request: req,
	}
}

func (req *Request) Setup() *SetupRequest {
	return &SetupRequest{
		Request: req,
	}
}

func (req *Request) Play() *PlayRequest {
	return &PlayRequest{
		Request: req,
	}
}

func (req *Request) Pause() *PauseRequest {
	return &PauseRequest{
		Request: req,
	}
}

func (req *Request) Teardown() *TeardownRequest {
	return &TeardownRequest{
		Request: req,
	}
}

func (req *Request) GetParameter() *GetParameterRequest {
	return &GetParameterRequest{
		Request: req,
	}
}

func (req *Request) SetParameter() *SetParameterRequest {
	return &SetParameterRequest{
		Request: req,
	}
}

func (req *Request) Record() *RecordRequest {
	return &RecordRequest{
		Request: req,
	}
}

// OptionsRequest is a RTSP OPTIONS request
type OptionsRequest struct {
	*Request
}

func (req *OptionsRequest) Require() string {
	return req.Lines["require"]
}

func (req *OptionsRequest) ProxyRequire() string {
	return req.Lines["proxy-require"]
}

// DescribeRequest is a RTSP DESCRIBE request
type DescribeRequest struct {
	*Request
}

func (req *DescribeRequest) Accept() string {
	return req.Lines["accept"]
}

// AnnounceRequest is a RTSP ANNOUNCE request
type AnnounceRequest struct {
	*Request
}

func (req *AnnounceRequest) ContentType() string {
	return req.Lines["content-type"]
}

// SetupRequest is a RTSP SETUP request
type SetupRequest struct {
	*Request
}

func (req *SetupRequest) TransportString() string {
	return req.Lines["transport"]
}

func (req *SetupRequest) Transport() (*Transport, error) {
	return NewTransport(req.TransportString())
}

// PlayRequest is a RTSP PLAY request
type PlayRequest struct {
	*Request
}

func (req *PlayRequest) Range() string {
	return req.Lines["range"]
}

// PauseRequest is a RTSP PAUSE request
type PauseRequest struct {
	*Request
}

func (req *PauseRequest) Range() string {
	return req.Lines["range"]
}

// TeardownRequest is a RTSP TEARDOWN request
type TeardownRequest struct {
	*Request
}

// GetParameterRequest is a RTSP GET_PARAMETER request
type GetParameterRequest struct {
	*Request
}

func (req *GetParameterRequest) Parameters() []string {
	if len(req.Content) == 0 {
		return nil
	}

	return strings.Split(string(req.Content), "\r\n")
}

// SetParameterRequest is a RTSP SET_PARAMETER request
type SetParameterRequest struct {
	*Request
}

func (req *SetParameterRequest) Parameters() map[string]string {
	if len(req.Content) == 0 {
		return nil
	}

	lines := strings.Split(string(req.Content), "\r\n")
	params := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		params[parts[0]] = parts[1]
	}

	return params
}

// RecordRequest is a RTSP RECORD request
type RecordRequest struct {
	*Request
}
