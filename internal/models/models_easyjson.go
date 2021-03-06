// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjsonD2b7633eDecodeWaveInternalModels(in *jlexer.Lexer, out *UserScore) {
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
		case "username":
			out.Username = string(in.String())
		case "score":
			out.Score = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels(out *jwriter.Writer, in UserScore) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"score\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Score))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserScore) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserScore) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserScore) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserScore) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels1(in *jlexer.Lexer, out *UserExtended) {
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
		case "username":
			out.Username = string(in.String())
		case "score":
			out.Score = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "locale":
			out.Locale = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels1(out *jwriter.Writer, in UserExtended) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"score\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Score))
	}
	{
		const prefix string = ",\"avatar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"locale\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Locale))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserExtended) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserExtended) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserExtended) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserExtended) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels1(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels2(in *jlexer.Lexer, out *UserEdit) {
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
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "locale":
			out.Locale = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels2(out *jwriter.Writer, in UserEdit) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
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
		const prefix string = ",\"avatar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"locale\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Locale))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserEdit) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserEdit) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserEdit) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserEdit) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels2(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels3(in *jlexer.Lexer, out *UserCredentials) {
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
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels3(out *jwriter.Writer, in UserCredentials) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserCredentials) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserCredentials) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserCredentials) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserCredentials) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels3(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels4(in *jlexer.Lexer, out *UserApplicationsInstalled) {
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
		case "user_apps_installed":
			if in.IsNull() {
				in.Skip()
				out.UserApplications = nil
			} else {
				in.Delim('[')
				if out.UserApplications == nil {
					if !in.IsDelim(']') {
						out.UserApplications = make([]UserApplication, 0, 1)
					} else {
						out.UserApplications = []UserApplication{}
					}
				} else {
					out.UserApplications = (out.UserApplications)[:0]
				}
				for !in.IsDelim(']') {
					var v1 UserApplication
					(v1).UnmarshalEasyJSON(in)
					out.UserApplications = append(out.UserApplications, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonD2b7633eEncodeWaveInternalModels4(out *jwriter.Writer, in UserApplicationsInstalled) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_apps_installed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.UserApplications == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.UserApplications {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserApplicationsInstalled) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserApplicationsInstalled) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserApplicationsInstalled) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserApplicationsInstalled) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels4(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels5(in *jlexer.Lexer, out *UserApplications) {
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
		case "user_apps":
			if in.IsNull() {
				in.Skip()
				out.UserApplications = nil
			} else {
				in.Delim('[')
				if out.UserApplications == nil {
					if !in.IsDelim(']') {
						out.UserApplications = make([]UserApplication, 0, 1)
					} else {
						out.UserApplications = []UserApplication{}
					}
				} else {
					out.UserApplications = (out.UserApplications)[:0]
				}
				for !in.IsDelim(']') {
					var v4 UserApplication
					(v4).UnmarshalEasyJSON(in)
					out.UserApplications = append(out.UserApplications, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonD2b7633eEncodeWaveInternalModels5(out *jwriter.Writer, in UserApplications) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_apps\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.UserApplications == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.UserApplications {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserApplications) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserApplications) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserApplications) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserApplications) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels5(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels6(in *jlexer.Lexer, out *UserApplicationInstalled) {
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
		case "link":
			out.Link = string(in.String())
		case "url":
			out.Url = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "name_de":
			out.NameDE = string(in.String())
		case "name_ru":
			out.NameRU = string(in.String())
		case "image":
			out.Image = string(in.String())
		case "about":
			out.About = string(in.String())
		case "about_de":
			out.AboutDE = string(in.String())
		case "about_ru":
			out.AboutRU = string(in.String())
		case "installs":
			out.Installations = int(in.Int())
		case "price":
			out.Price = int(in.Int())
		case "category":
			out.Category = string(in.String())
		case "installed":
			out.Installed = bool(in.Bool())
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
func easyjsonD2b7633eEncodeWaveInternalModels6(out *jwriter.Writer, in UserApplicationInstalled) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"link\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Link))
	}
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Url))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"name_de\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NameDE))
	}
	{
		const prefix string = ",\"name_ru\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NameRU))
	}
	{
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Image))
	}
	{
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"about_de\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.AboutDE))
	}
	{
		const prefix string = ",\"about_ru\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.AboutRU))
	}
	{
		const prefix string = ",\"installs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Installations))
	}
	{
		const prefix string = ",\"price\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Price))
	}
	{
		const prefix string = ",\"category\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"installed\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Installed))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserApplicationInstalled) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserApplicationInstalled) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserApplicationInstalled) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserApplicationInstalled) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels6(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels7(in *jlexer.Lexer, out *UserApplication) {
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
		case "link":
			out.Link = string(in.String())
		case "url":
			out.Url = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "image":
			out.Image = string(in.String())
		case "about":
			out.About = string(in.String())
		case "installs":
			out.Installations = int(in.Int())
		case "price":
			out.Price = int(in.Int())
		case "category":
			out.Category = string(in.String())
		case "time_total":
			out.TimeTotal = int(in.Int())
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
func easyjsonD2b7633eEncodeWaveInternalModels7(out *jwriter.Writer, in UserApplication) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"link\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Link))
	}
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Url))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Image))
	}
	{
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"installs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Installations))
	}
	{
		const prefix string = ",\"price\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Price))
	}
	{
		const prefix string = ",\"category\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"time_total\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.TimeTotal))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserApplication) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserApplication) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserApplication) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserApplication) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels7(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels8(in *jlexer.Lexer, out *Pagination) {
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
		case "page":
			out.Page = string(in.String())
		case "count":
			out.Count = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels8(out *jwriter.Writer, in Pagination) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"page\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Page))
	}
	{
		const prefix string = ",\"count\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Count))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Pagination) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Pagination) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Pagination) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Pagination) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels8(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels9(in *jlexer.Lexer, out *Leaders) {
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
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]UserScore, 0, 2)
					} else {
						out.Users = []UserScore{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v7 UserScore
					(v7).UnmarshalEasyJSON(in)
					out.Users = append(out.Users, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "total":
			out.Total = int(in.Int())
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
func easyjsonD2b7633eEncodeWaveInternalModels9(out *jwriter.Writer, in Leaders) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"users\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Users {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"total\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Total))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Leaders) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Leaders) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Leaders) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Leaders) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels9(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels10(in *jlexer.Lexer, out *ForbiddenRequest) {
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
		case "reason":
			out.Reason = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels10(out *jwriter.Writer, in ForbiddenRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"reason\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Reason))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ForbiddenRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ForbiddenRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ForbiddenRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ForbiddenRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels10(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels11(in *jlexer.Lexer, out *Categories) {
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
		case "categories":
			if in.IsNull() {
				in.Skip()
				out.Categories = nil
			} else {
				in.Delim('[')
				if out.Categories == nil {
					if !in.IsDelim(']') {
						out.Categories = make([]string, 0, 4)
					} else {
						out.Categories = []string{}
					}
				} else {
					out.Categories = (out.Categories)[:0]
				}
				for !in.IsDelim(']') {
					var v10 string
					v10 = string(in.String())
					out.Categories = append(out.Categories, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonD2b7633eEncodeWaveInternalModels11(out *jwriter.Writer, in Categories) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"categories\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Categories == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Categories {
				if v11 > 0 {
					out.RawByte(',')
				}
				out.String(string(v12))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Categories) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Categories) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Categories) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Categories) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels11(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels12(in *jlexer.Lexer, out *Applications) {
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
		case "apps":
			if in.IsNull() {
				in.Skip()
				out.Applications = nil
			} else {
				in.Delim('[')
				if out.Applications == nil {
					if !in.IsDelim(']') {
						out.Applications = make([]Application, 0, 1)
					} else {
						out.Applications = []Application{}
					}
				} else {
					out.Applications = (out.Applications)[:0]
				}
				for !in.IsDelim(']') {
					var v13 Application
					(v13).UnmarshalEasyJSON(in)
					out.Applications = append(out.Applications, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjsonD2b7633eEncodeWaveInternalModels12(out *jwriter.Writer, in Applications) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"apps\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Applications == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Applications {
				if v14 > 0 {
					out.RawByte(',')
				}
				(v15).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Applications) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Applications) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Applications) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Applications) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels12(l, v)
}
func easyjsonD2b7633eDecodeWaveInternalModels13(in *jlexer.Lexer, out *Application) {
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
		case "link":
			out.Link = string(in.String())
		case "url":
			out.Url = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "name_de":
			out.NameDE = string(in.String())
		case "name_ru":
			out.NameRU = string(in.String())
		case "image":
			out.Image = string(in.String())
		case "about":
			out.About = string(in.String())
		case "about_de":
			out.AboutDE = string(in.String())
		case "about_ru":
			out.AboutRU = string(in.String())
		case "installs":
			out.Installations = int(in.Int())
		case "price":
			out.Price = int(in.Int())
		case "category":
			out.Category = string(in.String())
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
func easyjsonD2b7633eEncodeWaveInternalModels13(out *jwriter.Writer, in Application) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"link\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Link))
	}
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Url))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"name_de\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NameDE))
	}
	{
		const prefix string = ",\"name_ru\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NameRU))
	}
	{
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Image))
	}
	{
		const prefix string = ",\"about\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"about_de\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.AboutDE))
	}
	{
		const prefix string = ",\"about_ru\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.AboutRU))
	}
	{
		const prefix string = ",\"installs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Installations))
	}
	{
		const prefix string = ",\"price\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Price))
	}
	{
		const prefix string = ",\"category\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Category))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Application) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeWaveInternalModels13(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Application) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeWaveInternalModels13(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Application) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeWaveInternalModels13(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Application) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeWaveInternalModels13(l, v)
}
