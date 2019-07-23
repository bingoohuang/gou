package str

import "strconv"

func ParseFloat32E(s string) (float32, error) { f, e := strconv.ParseFloat(s, 32); return float32(f), e }
func ParseFloat64E(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

func ParseIntE(s string) (int, error)     { p, e := ParseInt64E(s); return int(p), e }
func ParseInt8E(s string) (int8, error)   { p, e := strconv.ParseInt(s, 0, 8); return int8(p), e }
func ParseInt16E(s string) (int16, error) { p, e := strconv.ParseInt(s, 0, 16); return int16(p), e }
func ParseInt32E(s string) (int32, error) { p, e := strconv.ParseInt(s, 0, 32); return int32(p), e }
func ParseInt64E(s string) (int64, error) { return strconv.ParseInt(s, 0, 64) }

func ParseUintE(s string) (uint, error)     { p, e := ParseUint64E(s); return uint(p), e }
func ParseUint8E(s string) (uint8, error)   { p, e := strconv.ParseUint(s, 0, 8); return uint8(p), e }
func ParseUint16E(s string) (uint16, error) { p, e := strconv.ParseUint(s, 0, 16); return uint16(p), e }
func ParseUint32E(s string) (uint32, error) { p, e := strconv.ParseUint(s, 0, 32); return uint32(p), e }
func ParseUint64E(s string) (uint64, error) { return strconv.ParseUint(s, 0, 64) }

func ParseInt8(s string) int8   { i, _ := ParseInt8E(s); return i }
func ParseInt16(s string) int16 { i, _ := ParseInt16E(s); return i }
func ParseInt32(s string) int32 { i, _ := ParseInt32E(s); return i }
func ParseInt64(s string) int64 { i, _ := ParseInt64E(s); return i }

func ParseUint8(s string) uint8   { i, _ := ParseUint8E(s); return i }
func ParseUint16(s string) uint16 { i, _ := ParseUint16E(s); return i }
func ParseUint32(s string) uint32 { i, _ := ParseUint32E(s); return i }
func ParseUint64(s string) uint64 { i, _ := ParseUint64E(s); return i }

func ParseFloat32(s string) float32 { f, _ := ParseFloat32E(s); return f }
func ParseFloat64(s string) float64 { f, _ := ParseFloat64E(s); return f }

func ParseInt(s string) int   { f, _ := ParseIntE(s); return f }
func ParseUint(s string) uint { f, _ := ParseUintE(s); return f }
