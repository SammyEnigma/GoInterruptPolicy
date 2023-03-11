package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

var (
	modntdll       = syscall.NewLazyDLL("ntdll")
	procNtQueryKey = modntdll.NewProc("NtQueryKey")
)

// The KeyInformationClass constants have been derived from the KEY_INFORMATION_CLASS enum definition.
type KeyInformationClass uint32

const (
	KeyBasicInformation KeyInformationClass = iota
	KeyNodeInformation
	KeyFullInformation
	KeyNameInformation
	KeyCachedInformation
	KeyFlagsInformation
	KeyVirtualizationInformation
	KeyHandleTagsInformation
	MaxKeyInfoClass
)

const STATUS_BUFFER_TOO_SMALL = 0xC0000023

// OUT-parameter: KeyInformation, ResultLength.
// *OPT-parameter: KeyInformation.
func NtQueryKey(
	keyHandle uintptr,
	keyInformationClass KeyInformationClass,
	keyInformation *byte,
	length uint32,
	resultLength *uint32,
) int {
	r0, _, _ := procNtQueryKey.Call(keyHandle,
		uintptr(keyInformationClass),
		uintptr(unsafe.Pointer(keyInformation)),
		uintptr(length),
		uintptr(unsafe.Pointer(resultLength)))
	return int(r0)
}

// GetRegistryLocation(uintptr(device.reg))

func GetRegistryLocation(regHandle uintptr) (string, error) {
	var size uint32 = 0
	result := NtQueryKey(regHandle, KeyNameInformation, nil, 0, &size)
	if result == STATUS_BUFFER_TOO_SMALL {
		buf := make([]byte, size)
		if result := NtQueryKey(regHandle, KeyNameInformation, &buf[0], size, &size); result == 0 {
			regPath, err := DecodeUTF16(buf)
			if err != nil {
				log.Println(err)
			}

			tempRegPath := replaceRegistryMachine(regPath)
			tempRegPath = generalizeControlSet(tempRegPath)

			return `HKEY_LOCAL_MACHINE\` + tempRegPath, nil
		}
	} else {
		return "", fmt.Errorf("error: 0x%X", result)
	}
	return "", nil
}

func DecodeUTF16(b []byte) (string, error) {
	if len(b)%2 != 0 {
		return "", fmt.Errorf("must have even length byte slice")
	}

	u16s := make([]uint16, 1)
	ret := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}
