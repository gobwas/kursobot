package  yahoo

import (
  jlexer "github.com/mailru/easyjson/jlexer"
  json "encoding/json"
  jwriter "github.com/mailru/easyjson/jwriter"
)

var _ = json.RawMessage{} // suppress unused package warning

func easyjson_decode_finance_provider_yahoo_YahooRate(in *jlexer.Lexer, out *YahooRate) {
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
    case "id":
      out.Id = in.String()
    case "Name":
      out.Name = in.String()
    case "Rate":
      out.Rate = in.String()
    case "Date":
      out.Date = in.String()
    case "Time":
      out.Time = in.String()
    case "Ask":
      out.Ask = in.String()
    case "Bid":
      out.Bid = in.String()
    default:
      in.SkipRecursive()
    }
    in.WantComma()
  }
  in.Delim('}')
}
func easyjson_encode_finance_provider_yahoo_YahooRate(out *jwriter.Writer, in *YahooRate) {
  out.RawByte('{')
  first := true
  _ = first
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"id\":")
  out.String(in.Id)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"Name\":")
  out.String(in.Name)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"Rate\":")
  out.String(in.Rate)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"Date\":")
  out.String(in.Date)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"Time\":")
  out.String(in.Time)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"Ask\":")
  out.String(in.Ask)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"Bid\":")
  out.String(in.Bid)
  out.RawByte('}')
}
func (v *YahooRate) MarshalJSON() ([]byte, error) {
  w := jwriter.Writer{}
  easyjson_encode_finance_provider_yahoo_YahooRate(&w, v)
  return w.Buffer.BuildBytes(), w.Error
}
func (v *YahooRate) MarshalEasyJSON(w *jwriter.Writer) {
  easyjson_encode_finance_provider_yahoo_YahooRate(w, v)
}
func (v *YahooRate) UnmarshalJSON(data []byte) error {
  r := jlexer.Lexer{Data: data}
  easyjson_decode_finance_provider_yahoo_YahooRate(&r, v)
  return r.Error()
}
func (v *YahooRate) UnmarshalEasyJSON(l *jlexer.Lexer) {
  easyjson_decode_finance_provider_yahoo_YahooRate(l, v)
}
func easyjson_decode_finance_provider_yahoo_ResultRate(in *jlexer.Lexer, out *ResultRate) {
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
    case "rate":
      (out.Rate).UnmarshalEasyJSON(in)
    default:
      in.SkipRecursive()
    }
    in.WantComma()
  }
  in.Delim('}')
}
func easyjson_encode_finance_provider_yahoo_ResultRate(out *jwriter.Writer, in *ResultRate) {
  out.RawByte('{')
  first := true
  _ = first
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"rate\":")
  (in.Rate).MarshalEasyJSON(out)
  out.RawByte('}')
}
func (v *ResultRate) MarshalJSON() ([]byte, error) {
  w := jwriter.Writer{}
  easyjson_encode_finance_provider_yahoo_ResultRate(&w, v)
  return w.Buffer.BuildBytes(), w.Error
}
func (v *ResultRate) MarshalEasyJSON(w *jwriter.Writer) {
  easyjson_encode_finance_provider_yahoo_ResultRate(w, v)
}
func (v *ResultRate) UnmarshalJSON(data []byte) error {
  r := jlexer.Lexer{Data: data}
  easyjson_decode_finance_provider_yahoo_ResultRate(&r, v)
  return r.Error()
}
func (v *ResultRate) UnmarshalEasyJSON(l *jlexer.Lexer) {
  easyjson_decode_finance_provider_yahoo_ResultRate(l, v)
}
func easyjson_decode_finance_provider_yahoo_Result(in *jlexer.Lexer, out *Result) {
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
    case "count":
      out.Count = in.Int()
    case "created":
      out.Created = in.String()
    case "lang":
      out.Lang = in.String()
    case "results":
      (out.Results).UnmarshalEasyJSON(in)
    default:
      in.SkipRecursive()
    }
    in.WantComma()
  }
  in.Delim('}')
}
func easyjson_encode_finance_provider_yahoo_Result(out *jwriter.Writer, in *Result) {
  out.RawByte('{')
  first := true
  _ = first
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"count\":")
  out.Int(in.Count)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"created\":")
  out.String(in.Created)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"lang\":")
  out.String(in.Lang)
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"results\":")
  (in.Results).MarshalEasyJSON(out)
  out.RawByte('}')
}
func (v *Result) MarshalJSON() ([]byte, error) {
  w := jwriter.Writer{}
  easyjson_encode_finance_provider_yahoo_Result(&w, v)
  return w.Buffer.BuildBytes(), w.Error
}
func (v *Result) MarshalEasyJSON(w *jwriter.Writer) {
  easyjson_encode_finance_provider_yahoo_Result(w, v)
}
func (v *Result) UnmarshalJSON(data []byte) error {
  r := jlexer.Lexer{Data: data}
  easyjson_decode_finance_provider_yahoo_Result(&r, v)
  return r.Error()
}
func (v *Result) UnmarshalEasyJSON(l *jlexer.Lexer) {
  easyjson_decode_finance_provider_yahoo_Result(l, v)
}
func easyjson_decode_finance_provider_yahoo_Response(in *jlexer.Lexer, out *Response) {
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
    case "query":
      (out.Query).UnmarshalEasyJSON(in)
    default:
      in.SkipRecursive()
    }
    in.WantComma()
  }
  in.Delim('}')
}
func easyjson_encode_finance_provider_yahoo_Response(out *jwriter.Writer, in *Response) {
  out.RawByte('{')
  first := true
  _ = first
  if !first { out.RawByte(',') }
  first = false
  out.RawString("\"query\":")
  (in.Query).MarshalEasyJSON(out)
  out.RawByte('}')
}
func (v *Response) MarshalJSON() ([]byte, error) {
  w := jwriter.Writer{}
  easyjson_encode_finance_provider_yahoo_Response(&w, v)
  return w.Buffer.BuildBytes(), w.Error
}
func (v *Response) MarshalEasyJSON(w *jwriter.Writer) {
  easyjson_encode_finance_provider_yahoo_Response(w, v)
}
func (v *Response) UnmarshalJSON(data []byte) error {
  r := jlexer.Lexer{Data: data}
  easyjson_decode_finance_provider_yahoo_Response(&r, v)
  return r.Error()
}
func (v *Response) UnmarshalEasyJSON(l *jlexer.Lexer) {
  easyjson_decode_finance_provider_yahoo_Response(l, v)
}
