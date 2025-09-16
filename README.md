# itu-pickle

**itu-pickle** tarayıcı üzerinden arayüzüne ulaşıp kontrol edebileceğiniz, HTTP request tabanlı bir İTÜ ders seçicisidir.

## Build

```bash
git clone https://github.com/captainadorable/itu-pickle.git
cd itu-pickle/src
go get github.com/a-h/templ
templ generate
go build cmd/main.go
```
http://localhost:5454 adresinden arayüze ulaşabilirsin.

--port flagı ile port belirleyebilirsin.
--url flagı ile request atılacak URL'yi belirleyebilirsin.
