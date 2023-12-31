package signaling_client

import "encoding/json"

type CallRequest struct {
	//对端id
	PeerId string `json:"peerId"`
	//请求id
	CallId string `json:"callId"`
	//请求方法
	Method string `json:"method"`
	//请求携带的数据
	Data []byte `json:"data"`
}

func (r *CallRequest) FromBytes(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *CallRequest) ToBytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

type CallResponse struct {
	//请求id
	CallId string `json:"callId"`
	//请求方法
	Method string `json:"method"`
	//响应状态
	Status CallResponseCode `json:"status"`
	//响应携带的数据
	Data []byte `json:"data"`
}

func (r *CallResponse) ToBytes() []byte {
	marshal, _ := json.Marshal(r)
	return marshal
}

func (r *CallResponse) FromBytes(msg []byte) error {
	return json.Unmarshal(msg, r)
}

type CallResponseCode uint32

const (
	// CodeOK is returned on success.
	CodeOK CallResponseCode = 0

	// CodeCanceled indicates the operation was canceled (typically by the caller).
	//
	// The gRPC framework will generate this error code when cancellation
	// is requested.
	CodeCanceled CallResponseCode = 1

	// CodeUnknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	//
	// The gRPC framework will generate this error code in the above two
	// mentioned cases.
	CodeUnknown CallResponseCode = 2

	// CodeInvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	//
	// This error code will not be generated by the gRPC framework.
	CodeInvalidArgument CallResponseCode = 3

	// CodeDeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	//
	// The gRPC framework will generate this error code when the deadline is
	// exceeded.
	CodeDeadlineExceeded CallResponseCode = 4

	// CodeNotFound means some requested entity (e.g., file or directory) was
	// not found.
	//
	// This error code will not be generated by the gRPC framework.
	CodeNotFound CallResponseCode = 5

	// CodeAlreadyExists means an attempt to create an entity failed because one
	// already exists.
	//
	// This error code will not be generated by the gRPC framework.
	CodeAlreadyExists CallResponseCode = 6

	// CodePermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	//
	// This error code will not be generated by the gRPC core framework,
	// but expect authentication middleware to use it.
	CodePermissionDenied CallResponseCode = 7

	// CodeResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	//
	// This error code will be generated by the gRPC framework in
	// out-of-memory and server overload situations, or when a message is
	// larger than the configured maximum size.
	CodeResourceExhausted CallResponseCode = 8

	// CodeFailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	//
	// This error code will not be generated by the gRPC framework.
	CodeFailedPrecondition CallResponseCode = 9

	// CodeAborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	//
	// This error code will not be generated by the gRPC framework.
	CodeAborted CallResponseCode = 10

	// CodeOutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	//
	// This error code will not be generated by the gRPC framework.
	CodeOutOfRange CallResponseCode = 11

	// CodeUnimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	//
	// This error code will be generated by the gRPC framework. Most
	// commonly, you will see this error code when a method implementation
	// is missing on the server. It can also be generated for unknown
	// compression algorithms or a disagreement as to whether an RPC should
	// be streaming.
	CodeUnimplemented CallResponseCode = 12

	// CodeInternal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	//
	// This error code will be generated by the gRPC framework in several
	// internal error conditions.
	CodeInternal CallResponseCode = 13

	// CodeUnavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	//
	// This error code will be generated by the gRPC framework during
	// abrupt shutdown of a server process or network connection.
	CodeUnavailable CallResponseCode = 14

	// CodeDataLoss indicates unrecoverable data loss or corruption.
	//
	// This error code will not be generated by the gRPC framework.
	CodeDataLoss CallResponseCode = 15

	// CodeUnauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	//
	// The gRPC framework will generate this error code when the
	// authentication metadata is invalid or a Credentials callback fails,
	// but also expect authentication middleware to generate it.
	CodeUnauthenticated CallResponseCode = 16

	_maxCode = 17
)

var strToCode = map[string]CallResponseCode{
	`"OK"`: CodeOK,
	`"CANCELLED"`:/* [sic] */ CodeCanceled,
	`"UNKNOWN"`:             CodeUnknown,
	`"INVALID_ARGUMENT"`:    CodeInvalidArgument,
	`"DEADLINE_EXCEEDED"`:   CodeDeadlineExceeded,
	`"NOT_FOUND"`:           CodeNotFound,
	`"ALREADY_EXISTS"`:      CodeAlreadyExists,
	`"PERMISSION_DENIED"`:   CodePermissionDenied,
	`"RESOURCE_EXHAUSTED"`:  CodeResourceExhausted,
	`"FAILED_PRECONDITION"`: CodeFailedPrecondition,
	`"ABORTED"`:             CodeAborted,
	`"OUT_OF_RANGE"`:        CodeOutOfRange,
	`"UNIMPLEMENTED"`:       CodeUnimplemented,
	`"INTERNAL"`:            CodeInternal,
	`"UNAVAILABLE"`:         CodeUnavailable,
	`"DATA_LOSS"`:           CodeDataLoss,
	`"UNAUTHENTICATED"`:     CodeUnauthenticated,
}

var codeToStr = map[CallResponseCode]string{}

func (c CallResponseCode) String() string {
	s, ok := codeToStr[c]
	if !ok {
		return "UNKNOWN"
	}
	return s
}
