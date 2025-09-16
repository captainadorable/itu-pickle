package utils

import (
	"bufio"
	"net/http"
	"strings"
)

type Course struct {
	Crn string
	Code string
	Type string
	Teacher string
	Building string
	Exists bool
}

var ReturnValues = map[string]string{
	"successResult": "CRN %s için işlem başarıyla tamamlandı.",
	"errorResult": "CRN %s için Operasyon tamamlanamadı.",
	"": "CRN %s için Operasyon tamamlanamadı.",
	"error": "CRN %s için bir hata meydana geldi.",
	"VAL01": "CRN %s bir problemden dolayı alınamadı.",
	"VAL02": "CRN %s kayıt zaman engelinden dolayı alınamadı.",
	"VAL03": "CRN %s bu dönem zaten alındığından dolayı tekrar alınamadı.",
	"VAL04": "CRN %s ders planında yer almadığından dolayı alınamadı.",
	"VAL05": "CRN %s dönemlik maksimum kredi sınırını aştığından dolayı alınamadı.",
	"VAL06": "CRN %s kontenjan yetersizliğinden dolayı alınamadı.",
	"VAL07": "CRN %s daha önce AA notuyla verildiğinden dolayı alınamadı.",
	"VAL08": "CRN %s program şartını sağlamadığından dolayı alınamadı.",
	"VAL09": "CRN %s başka bir dersle çakıştığından dolayı alınamadı.",
	"VAL10": "CRN %s dersine kayıtlı olmadığınızdan dolayı hiç bir işlem yapılmadı.",
	"VAL11": "CRN %s önşartlardan dolayı alınamadı.",
	"VAL12": "CRN %s şu anki dönemde açılmadığından dolayı alınamadı.",
	"VAL13": "CRN %s geçici olarak engellenmiş olması sebebiyle alınamadı.",
	"VAL14": "Sistem geçici olarak yanıt vermiyor.",
	"VAL15": "Maksimum 12 CRN alabilirsiniz.",
	"VAL16": "Aktif bir işleminiz devam ettiğinden dolayı işlem yapılamadı.",
	"VAL18": "CRN %s engellendiğinden dolayı alınamadı.",
	"VAL19": "CRN %s önlisans dersi olduğundan dolayı alınamadı.",
	"VAL20": "Dönem başına sadece 1 ders bırakabilirsiniz.",
	"CRNListEmpty": "CRN %s listesi boş göründüğünden alınamadı.",
	"CRNNotFound": "CRN %s bulunamadığından dolayı alınamadı.",
	"ERRLoad": "Sistem geçici olarak yanıt vermiyor.",
	"NULLParam-CheckOgrenciKayitZamaniKontrolu": "CRN %s kayıt zaman engelinden dolayı alınamadı.",
	"Ekleme İşlemi Başarılı": "CRN %s için ekleme işlemi başarıyla tamamlandı.",
	"Kontenjan Dolu": "CRN %s için kontenjan dolu olduğundan dolayı alınamadı.",
}

var crnToCourse map[string]string = nil

func GetCrns() (map[string]string, error) {
	url := "https://raw.githubusercontent.com/itu-helper/data/main/lessons.psv"

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
    return nil, err
	}
	defer resp.Body.Close()

	// Create a map to store lessons
	crnToLesson := make(map[string]string)

	// Read the response line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		elements := strings.Split(line, "|")
		if len(elements) > 0 {
			crnToLesson[elements[0]] = line
		}
	}

	if err := scanner.Err(); err != nil {
    return nil, err
	}

	return crnToLesson, nil
}

func FindCrns(crns []string) ([]Course, error) {
	if crnToCourse == nil {
		c, err := GetCrns()
		if err != nil {
			return nil, err
		}
		crnToCourse = c
	}
	
  var courseList []Course
  for _, crn := range crns {
    c, exists := crnToCourse[crn]
    if exists {
      var data = strings.Split(c, "|")
			course := Course{
				Crn: data[0],	
				Code: data[1],	
				Type: data[2],	
				Teacher: data[3],	
				Building: data[4],	
				Exists: true,
			}
			courseList = append(courseList, course)
    } else {
			course := Course{
				Crn: crn,	
				Exists: false,
			}
			courseList = append(courseList, course)
    }
  }
  return courseList, nil
}
