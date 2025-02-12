# itu-pickle

**itu-pickle** tarayıcı üzerinden arayüzüne ulaşıp kontrol edebileceğiniz, HTTP request tabanlı bir İTÜ ders seçicisidir.

![demo](https://safe.captador.space/8ucdHikft01b.png)

## Build

```bash
git clone https://github.com/captainadorable/itu-pickle.git
cd itu-pickle/src
go get github.com/a-h/templ
templ generate
go build cmd/main.go
```
http://localhost:3001 adresinden arayüze ulaşabilirsin.

--port flagı ile port belirleyebilirsin.
--url flago ile request atılacak URL'yi belirleyebilirsin.

## TODO

- Girilen bütün dersler alındığında otomatik olarak duran sistem.
- Autoscroll.
