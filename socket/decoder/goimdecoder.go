package goimdecoder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Message struct {
	Version   uint32
	Operation uint32
	Seq       uint32
	Body      []byte
}

const MaxBodySize = (1 << 16)

const (
	_packageLen     = 4
	_headerLen      = 2
	_versionLen     = 2
	_operationLen   = 4
	_sequenceIdLen  = 4
	_totalHeaderLen = _packageLen + _headerLen + _versionLen + _operationLen + _sequenceIdLen
	_maxPackSize    = MaxBodySize + _totalHeaderLen
	// offset
	_packageLenOffset = 0
	_headerLenOffset  = _packageLenOffset + _packageLen
	_versionOffset    = _headerLenOffset + _headerLen
	_operationOffset  = _versionOffset + _versionLen
	_sequenceIdOffset = _operationOffset + _operationLen
	_bodyOffset       = _sequenceIdOffset + _sequenceIdLen
)

var (
	ErrParam      = errors.New("input buf len is too short")
	ErrPackageLen = errors.New("server codec package length error")
	ErrHeaderLen  = errors.New("server codec header length error")
	ErrIncomplete = errors.New("data in buffer not enough for one message")
)

func GoimDecoder(reader *io.Reader, buf []byte) (msg *Message, err error) {
	if len(buf) < _totalHeaderLen {
		fmt.Println("buffer len is too short!")
		err = ErrParam
		return nil, err
	}

	var packLen = binary.BigEndian.Uint32(buf[_packageLenOffset:_headerLenOffset])
	var headerLen = binary.BigEndian.Uint16(buf[_headerLenOffset:_versionOffset])

	if packLen > _maxPackSize {
		err = ErrPackageLen
		return nil, err
	}

	if headerLen != _totalHeaderLen {
		err = ErrHeaderLen
		return nil, err
	}

	if packLen > uint32(len(buf)) {
		err = ErrIncomplete
		return nil, err
	}

	msg.Version = uint32(binary.BigEndian.Uint16(buf[_versionOffset:_operationOffset]))
	msg.Operation = binary.BigEndian.Uint32(buf[_operationOffset:_sequenceIdOffset])
	msg.Seq = binary.BigEndian.Uint32(buf[_sequenceIdOffset:])
	msg.Body = buf[_bodyOffset:packLen]

	return msg, nil
}
