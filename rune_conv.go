package httpc

//#include <stdlib.h>
//#include <unicode/ucnv.h>
//#cgo LDFLAGS: -licuuc
/*
int pointerLess(char *a, char *b) {
  if (a < b) return 1;
  return 0;
}
*/
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

func From(encoding string, in []byte) (ret []rune, err error) {
	var errCode C.UErrorCode
	cEncoding := C.CString(encoding)
	defer C.free(unsafe.Pointer(cEncoding))
	conv := C.ucnv_open(cEncoding, &errCode)
	if errCode != C.U_ZERO_ERROR {
		err = errors.New("to_runes: error when open converter")
		return
	}
	defer C.ucnv_close(conv)
	ret = make([]rune, 0, len(in))
	var target C.UChar32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	source := (*C.char)(unsafe.Pointer(header.Data))
	sourceLimit := (*C.char)(unsafe.Pointer(header.Data + uintptr(header.Len)))
	for C.pointerLess(source, sourceLimit) == C.int(1) {
		target = C.ucnv_getNextUChar(conv, &source, sourceLimit, &errCode)
		if errCode != C.U_ZERO_ERROR {
			err = errors.New("to_runes: error when decoding")
			return
		}
		ret = append(ret, rune(target))
	}
	return
}
