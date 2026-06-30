// Copyright (c) 2026 TeamDev
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package win

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
)

const (
	iconDirSize   = 6
	iconEntrySize = 16
)

type IconEntry struct {
	WidthByte  byte
	HeightByte byte
	ColorCount byte
	Reserved   byte

	PlanesOrHotspotX uint16
	BitCountOrHotY   uint16
	BytesInRes       uint32
	ImageOffset      uint32

	Data []byte
	Orig int
}

func (e IconEntry) Width() int {
	if e.WidthByte == 0 {
		return 256
	}
	return int(e.WidthByte)
}

func (e IconEntry) Height() int {
	if e.HeightByte == 0 {
		return 256
	}
	return int(e.HeightByte)
}

func sortICO(data []byte, descending bool) ([]byte, error) {
	if len(data) < iconDirSize {
		return nil, errors.New("file too small for ICO header")
	}

	reserved := binary.LittleEndian.Uint16(data[0:2])
	iconType := binary.LittleEndian.Uint16(data[2:4])
	count := binary.LittleEndian.Uint16(data[4:6])

	if reserved != 0 {
		return nil, fmt.Errorf("invalid ICO reserved field: %d", reserved)
	}
	if iconType != 1 && iconType != 2 {
		return nil, fmt.Errorf("unsupported icon type %d; expected 1 for ICO or 2 for CUR", iconType)
	}
	if count == 0 {
		return nil, errors.New("ICO contains zero images")
	}

	dirEnd := iconDirSize + int(count)*iconEntrySize
	if len(data) < dirEnd {
		return nil, errors.New("file too small for ICO directory entries")
	}

	entries := make([]IconEntry, 0, count)

	for i := 0; i < int(count); i++ {
		p := iconDirSize + i*iconEntrySize
		raw := data[p : p+iconEntrySize]

		e := IconEntry{
			WidthByte:        raw[0],
			HeightByte:       raw[1],
			ColorCount:       raw[2],
			Reserved:         raw[3],
			PlanesOrHotspotX: binary.LittleEndian.Uint16(raw[4:6]),
			BitCountOrHotY:   binary.LittleEndian.Uint16(raw[6:8]),
			BytesInRes:       binary.LittleEndian.Uint32(raw[8:12]),
			ImageOffset:      binary.LittleEndian.Uint32(raw[12:16]),
			Orig:             i,
		}

		start := int(e.ImageOffset)
		end := start + int(e.BytesInRes)

		if start < 0 || end < start || end > len(data) {
			return nil, fmt.Errorf(
				"entry %d has invalid image range: offset=%d size=%d fileSize=%d",
				i,
				e.ImageOffset,
				e.BytesInRes,
				len(data),
			)
		}

		e.Data = append([]byte(nil), data[start:end]...)
		entries = append(entries, e)
	}

	sort.SliceStable(entries, func(i, j int) bool {
		a := entries[i]
		b := entries[j]

		// Primary sort: larger/smaller width.
		if a.Width() != b.Width() {
			if descending {
				return a.Width() > b.Width()
			}
			return a.Width() < b.Width()
		}

		// Secondary sort: larger/smaller height.
		if a.Height() != b.Height() {
			if descending {
				return a.Height() > b.Height()
			}
			return a.Height() < b.Height()
		}

		// Tertiary sort: prefer higher bit depth.
		if a.BitCountOrHotY != b.BitCountOrHotY {
			if descending {
				return a.BitCountOrHotY > b.BitCountOrHotY
			}
			return a.BitCountOrHotY < b.BitCountOrHotY
		}

		// Preserve original order for exact ties.
		return a.Orig < b.Orig
	})

	out := make([]byte, 0, len(data))

	// ICONDIR header.
	out = append(out, data[:iconDirSize]...)

	// Reserve space for sorted ICONDIRENTRY records.
	entryStart := len(out)
	out = append(out, make([]byte, int(count)*iconEntrySize)...)

	// Write image data in the same order as the sorted directory.
	currentOffset := iconDirSize + int(count)*iconEntrySize

	for i := range entries {
		entries[i].ImageOffset = uint32(currentOffset)
		entries[i].BytesInRes = uint32(len(entries[i].Data))

		currentOffset += len(entries[i].Data)
		out = append(out, entries[i].Data...)
	}

	// Write sorted directory entries with corrected offsets.
	for i, e := range entries {
		p := entryStart + i*iconEntrySize

		out[p+0] = e.WidthByte
		out[p+1] = e.HeightByte
		out[p+2] = e.ColorCount
		out[p+3] = e.Reserved

		binary.LittleEndian.PutUint16(out[p+4:p+6], e.PlanesOrHotspotX)
		binary.LittleEndian.PutUint16(out[p+6:p+8], e.BitCountOrHotY)
		binary.LittleEndian.PutUint32(out[p+8:p+12], e.BytesInRes)
		binary.LittleEndian.PutUint32(out[p+12:p+16], e.ImageOffset)
	}

	return out, nil
}
