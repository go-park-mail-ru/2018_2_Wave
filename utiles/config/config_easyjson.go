// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package config

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6615c02eDecodeWaveResourcesConfig(in *jlexer.Lexer, out *ServerConfiguration) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "host":
			out.Host = string(in.String())
		case "port":
			out.Port = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeWaveResourcesConfig(out *jwriter.Writer, in ServerConfiguration) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"host\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Host))
	}
	{
		const prefix string = ",\"port\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Port))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ServerConfiguration) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeWaveResourcesConfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ServerConfiguration) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeWaveResourcesConfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ServerConfiguration) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeWaveResourcesConfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ServerConfiguration) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeWaveResourcesConfig(l, v)
}
func easyjson6615c02eDecodeWaveResourcesConfig1(in *jlexer.Lexer, out *DatabaseConfiguration) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			out.User = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "dbname":
			out.DBName = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeWaveResourcesConfig1(out *jwriter.Writer, in DatabaseConfiguration) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"dbname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.DBName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v DatabaseConfiguration) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeWaveResourcesConfig1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v DatabaseConfiguration) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeWaveResourcesConfig1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *DatabaseConfiguration) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeWaveResourcesConfig1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *DatabaseConfiguration) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeWaveResourcesConfig1(l, v)
}
func easyjson6615c02eDecodeWaveResourcesConfig2(in *jlexer.Lexer, out *Configuration) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "server":
			(out.SC).UnmarshalEasyJSON(in)
		case "database":
			(out.DC).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6615c02eEncodeWaveResourcesConfig2(out *jwriter.Writer, in Configuration) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"server\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.SC).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"database\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.DC).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Configuration) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6615c02eEncodeWaveResourcesConfig2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Configuration) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6615c02eEncodeWaveResourcesConfig2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Configuration) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6615c02eDecodeWaveResourcesConfig2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Configuration) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6615c02eDecodeWaveResourcesConfig2(l, v)
}
