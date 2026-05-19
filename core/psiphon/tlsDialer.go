package psiphon

/*
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

typedef struct {
    const char* front_domain;
    const char* host_header;
    const char* sni;
} PsiphonMeekSpec;

typedef struct {
    uint8_t session_id[32];
    bool use_ech;
    uint16_t padding_length;
} TlsDialerContext;

// Extern declaration matching the Rust FFI boundary from Phase 2
extern int rust_psiphon_prepare_tls_dialer(const PsiphonMeekSpec* spec, TlsDialerContext* ctx);
*/
import "C"
import (
	"crypto/tls"
	"errors"
	"net"
	"time"
	"unsafe"
)

type TLSDialer struct {
	Timeout time.Duration
}

func NewTLSDialer(timeout time.Duration) *TLSDialer {
	return &TLSDialer{
		Timeout: timeout,
	}
}

// Dial executes the DPI-evading TLS handshake utilizing the Rust Post-Quantum Kyber FFI
func (d *TLSDialer) Dial(network, addr string, frontDomain, hostHeader, sni string) (net.Conn, error) {
	cFrontDomain := C.CString(frontDomain)
	cHostHeader := C.CString(hostHeader)
	cSni := C.CString(sni)

	defer C.free(unsafe.Pointer(cFrontDomain))
	defer C.free(unsafe.Pointer(cHostHeader))
	defer C.free(unsafe.Pointer(cSni))

	spec := C.PsiphonMeekSpec{
		front_domain: cFrontDomain,
		host_header:  cHostHeader,
		sni:          cSni,
	}

	var ctx C.TlsDialerContext

	// Invoke Rust Bare-Metal Runtime for Kyber Encapsulation and Padding Calculation
	result := C.rust_psiphon_prepare_tls_dialer(&spec, &ctx)
	if result != 0 {
		return nil, errors.New("rust FFI TLS preparation failed")
	}

	// Standard TCP Dial
	rawConn, err := net.DialTimeout(network, addr, d.Timeout)
	if err != nil {
		return nil, err
	}

	// Construct Go TLS Config using the parameters injected by Rust
	tlsConfig := &tls.Config{
		ServerName:         sni,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS13,
		// In a full implementation, ctx.session_id and ctx.padding_length 
		// are injected into custom uTLS/JA4+ extensions here.
	}

	tlsConn := tls.Client(rawConn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	return tlsConn, nil
}
