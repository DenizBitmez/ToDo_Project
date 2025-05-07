# TO-DO App REST API

## Genel Bakış

Bu proje, **Golang** ve **Gin** framework'ü kullanılarak geliştirilmiş bir TO-DO uygulamasının REST API servisidir. Kullanıcılar yapılacaklar listesi (TO-DO) oluşturabilir, bu listelere adımlar (steps) ekleyebilir ve bunları yönetebilir. Sistemde iki kullanıcı rolü vardır: **normal kullanıcı** ve **admin**.

## Özellikler

*  JWT ile kimlik doğrulama
*  Rol bazlı yetkilendirme (admin/user)
*  TO-DO listesi oluşturma, güncelleme ve silme (soft delete)
*  Her listeye ait adım (step) oluşturma, güncelleme ve silme
*  Soft delete özelliği ile veri kaybı olmadan silme işlemleri
*  Admin tüm kullanıcıların verilerine erişebilir

## Kullanılan Teknolojiler

* **Programlama Dili**: Golang
* **Framework**: Gin
* **Kimlik Doğrulama**: JWT
* **Veritabanı**: In-memory (Mock repository)
* **API Testi**: Postman
* **Frontend**: HTML, CSS, JavaScript
* **Hosting**: Render (frontend için)

## Kullanıcı Bilgileri

**Normal Kullanıcı**

    * Kullanıcı Adı: `user1`
    * Şifre: `user123`
   
**Admin**

    * Kullanıcı Adı: `admin`
    * Şifre: `admin123`

## Kurulum

```bash
git clone https://github.com/kullaniciadi/todo-api.git
cd todo-api
go mod download
go run main.go
```

## API Endpointleri (Postman ile test edilebilir)

### Authentication

| Yöntem | Endpoint | Açıklama        |
| ------ | -------- | --------------- |
| POST   | `/login` | JWT token alımı |

### TO-DO Listeleri

| Yöntem | Endpoint     | Açıklama                               |
| ------ | ------------ | -------------------------------------- |
| GET    | `/todos`     | Kullanıcının TO-DO listelerini getirir |
| POST   | `/todos`     | Yeni TO-DO listesi oluşturur           |
| PUT    | `/todos/:id` | TO-DO listesini günceller              |
| DELETE | `/todos/:id` | TO-DO listesini soft-delete yapar      |

### TO-DO Adımları (Steps)

| Yöntem | Endpoint     | Açıklama                                              |
| ------ | ------------ | ----------------------------------------------------- |
| GET    | `/steps`     | Kullanıcının adımlarını getirir (admin hepsini görür) |
| POST   | `/steps`     | Yeni adım oluşturur                                   |
| PUT    | `/steps/:id` | Adımı günceller                                       |
| DELETE | `/steps/:id` | Adımı siler (soft delete)                             |

> Admin kullanıcı `/admin/steps` endpointi ile tüm adımları görebilir.

## Mimari

Proje **Clean Architecture** prensiplerine uygun geliştirilmiştir:

* **Handler (Controller)**: HTTP isteklerini işler.
* **Service**: İş mantığını barındırır.
* **Repository**: Veri erişim işlemleri.
* **Model**: Veri yapıları (struct'lar).
* **MiddleWare**: Yetkilendirme işlemlerini yapar.
* **pkg/jwt**: Token işlemlerini yürütür.
## Frontend

Bu projeye ait bir frontend arayüzü de geliştirilmiştir. HTML, CSS ve JavaScript kullanılarak hazırlanmıştır. Aşağıdaki bağlantıdan uygulamayı canlı olarak görüntüleyebilirsiniz:

🔗 [Canlı Demo (Render)](https://todo-project-69kz.onrender.com)


##  Katkı

Katkıda bulunmak isterseniz lütfen fork edip pull request gönderin. Büyük değişiklikler öncesinde issue açmanız önerilir.
##  Lisans

MIT Lisansı (veya kullandığın diğer lisans)
