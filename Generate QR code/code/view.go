package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

/* fungsi yang digunakan untuk mengambil endpoint sesuai akses , dengan function yang telah dipilih */
func main() {
	fmt.Println("halo dunia")

	// membuat fungsi mux sebagai router
	router := mux.NewRouter()

	// membuat handlefuction untuk mengakses endpoint dengan method get
	router.HandleFunc("/view",view).Methods("GET")
	router.HandleFunc("/download",download).Methods("GET")

	// membuat server untuk mengakses endpoint
	log.Fatal(http.ListenAndServe(":8000",router))
}	

func sayaHello(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "Hallo dunia")
}

// function untuk view qr saja
func view(w http.ResponseWriter, r *http.Request)  {

	// mengambil form sesuai inputan user 
	teks := r.FormValue("teks")
	label := r.FormValue("label")

	// mengambil qr yang telah generate secara final
	qr := buatQR(teks,label)

	// mengirimkan gambar kepada user
	w.Header().Set("Content-type","image/jpeg")
	jpeg.Encode(w,qr,&jpeg.Options{Quality: 100})

}

// fungsi untuk download atau menyimpan gambar di komputer
func download(w http.ResponseWriter, r *http.Request)  {

	// mengambil inputan dari user dengan form value
	teks := r.FormValue("teks")
	label := r.FormValue("label")
	unduh := r.FormValue("unduh")

	unduh_final := unduh + ".jpeg"
	// mengambil qr nya
	qr := buatQR(teks,label)

	// membuat file untuk disimpan di file

	
	folder := "C:/Smartlink QR/"
	os.MkdirAll(folder, os.ModePerm)

	file_final := filepath.Join(folder,unduh_final)	
	file,_ := os.Create(file_final)
	name_file := io.Writer(file)

	// mengambil 3 karakter char terakhir
	// tipe_data := unduh[(len(unduh)-3):]

	// penkondisian untuk menetukan tipe file yang digenerate
	// if tipe_data == "png" {
	// 	png.Encode(name_file,qr)
	// }else if tipe_data == "jpeg" {
	jpeg.Encode(name_file,qr,&jpeg.Options{Quality: 100})
	// }else{
	// 	log.Fatal("tipe tidak didukung")
	// }

	http.ServeFile(w,r,file_final)
}

//buat QR nya
func buatQR(teks string , label string) image.Image{	


	//random teks nya
	text := teks

	//random teksnya
	random_teks := random_kode_view(text)
	fmt.Println(random_teks)

	//ini qr nya
	kodeqr , _:= qr.Encode(random_teks,qr.M,qr.Auto)

	// mengatur lebar dan tinggi dari qr
	width := 315
	height := 290 
	kodeqr,_ = barcode.Scale(kodeqr,width,height)

	// mengambil logo untuk disisipkan di qr
	gambar,err_gambar := os.Open("logo_sm_3.png")
	if err_gambar != nil {
		log.Fatal(err_gambar)
	}

	// ubah tipe gambar dari os ke image
	gambar_decode,_ := png.Decode(gambar)

	// unutk mengubah ukuran pada gambar atau logo yang disisipkan
	buat_qr := resize.Resize(uint(width)/5, uint(height)/5, gambar_decode, resize.Bilinear) 

	// mendapatkan gambar qr final yang telah di setting
	gambar_final := settingFinal(kodeqr,buat_qr,380,389,label)

	// mengembalikan variabel gambar_final
	return gambar_final

	// w.Header().Set("Content-type","image/jpeg") 
	// final := jpeg.Encode(w,kodeqr,&jpeg.Options{100})
	// if final == nil {
	// 	http.Error(w,"gagal menampilkan qr code",http.StatusInternalServerError)
	// 	return
	// }
	// //buat nama filenya
	// save , err_save:= os.Create("qr_pake_booler_20.jpg")
	// if err_save != nil {
	// 	log.Fatal(err_save)
	// }

	// //ini berfungsi untuk mengatur kualitas qr --> qr dimasukan kedalam file --> file baru ada qr nya
	// final_qr := jpeg.Encode(save,gambar_final, &jpeg.Options{Quality: 100})
	// if final_qr == nil {
	// 	fmt.Println("qr selesai dibuat")
	// }

}

// fungsi untuk menset atau menggabungkan atara qr dengan logo dan label
func settingFinal(kodeqr barcode.Barcode,gambar_decode image.Image,width int, height int,label string) image.Image {
	
	// mendapatkan gambar label
	label_final := buat_label(label)

	// untuk lebar(x) dan tinggi(y)
	x := width
	y := height

	// membuat kanvas kosong
	kanvas:= image.NewNRGBA(image.Rect(0,0,x,y))
	kotak := image.Rect(0,0,400,300)

	// membuat color
	random_warna_1 := rand.Int31n(256)
	random_warna_2 := rand.Int31n(256)
	random_warna_3 := rand.Int31n(256)
	random_warna_4 := rand.Int31n(256)

	warna := color.RGBA{uint8(random_warna_1), uint8(random_warna_2), uint8(random_warna_3), uint8(random_warna_4)}

	// mendraw kanvas kosong dengan luasnya sesuai dengan bounds dari kanvas
	draw.Draw(kanvas,kanvas.Bounds(),&image.Uniform{warna},image.ZP,draw.Over)
	
	// menambahkan qr di kanvas
	draw.Draw(kanvas,kotak.Bounds().Add(image.Pt(32,30)),kodeqr,image.ZP,draw.Over)

	// menambahkan logo di kanvas
	draw.Draw(kanvas,kodeqr.Bounds().Add(image.Pt((width/2)-25,(height/3)+20)),gambar_decode,image.ZP,draw.Src)

	// menambahkan label di kanvas
	draw.Draw(kanvas,kanvas.Bounds().Add(image.Pt((x/3)-95,y-50)),label_final,image.ZP,draw.Src)

	// mengembalikan kanvas yang telah di draw
	return kanvas
}

// fungsi untuk generate isi qr
func random_kode_view(teks string) string  {	
	biner := teks

	final_random := biner

	return final_random
}

// fungsi untuk membuat label atau tulisan dengan tipe gambar
func buat_label(label string) image.Image  {	

	// mengambil file font
	gaya_tulisan,_ := os.ReadFile("D:/PROGRAM/go program/src/Project/ostrich-sans/ostrich-regular.ttf")
	
	// membuat agar font tadi bisa di terhubung dengan library freetype
	gaya,_ := truetype.Parse(gaya_tulisan)

	// membuat fungsi freetype
	setting := freetype.NewContext()

	// setting font
	setting.SetFont(gaya)
	// setting kualitas perpixel
	setting.SetDPI(100)
	// setting ukuran tulisan
	setting.SetFontSize(24)	
	var a int
	a = 0
	for i := 0; i < len(label); i++ {
		a += 11
	}
	// setting clip untuk menggambar tulisan
	setting.SetClip(image.Rect(0,0,315,30))


	// membuat kanvas kosong
	ambil := image.NewNRGBA(image.Rect(0,0,315,30))
	// mendraw kanvas dengan warna putih
	draw.Draw(ambil,ambil.Bounds(),image.White,image.ZP,draw.Src)
	
	// setting untuk gambar nya ditulis atau targetnya di destination
	setting.SetDst(ambil)
	// setting untuk tulisan dengan warna hitam
	setting.SetSrc(image.Black)

	// teks akan diikuti dengan kata QR + isi label
	teks_final := label
	
	var x int
	
	if len(label) > 11 {
		x = ambil.Bounds().Dx() / 4
	} else if len(label) >= 5 && len(label) < 12{
		x = (ambil.Bounds().Dx() / 3) + 10
	} else if len(label) == 4 {
		x = (ambil.Bounds().Dx() / 2) -10
	}else if len(label) == 3 {
		x = (ambil.Bounds().Dx() / 2)-13 
	}else if len(label) == 2 {
		x = (ambil.Bounds().Dx() / 2)
	}else if len(label) == 1 {
		x = (ambil.Bounds().Dx() / 2)
	}
	
	fmt.Println(len(label))
	
	y := (ambil.Bounds().Dy()) - 5
	fmt.Println(x)

	// mendraw label dari string ke image dengan teks dari label dan posisi yang disesuaikan
	setting.DrawString(teks_final,freetype.Pt(x,y))
	fmt.Println(label)
	// mengembalikan nilai ambil
	return ambil
}

// func qrSettingsLogo(gambar image.Image , kodeqr barcode.Barcode) image.Image {


// 	ukuran_kanvas := image.Rect(0,0,256,256)
// 	kanvas_kosong := image.NewNRGBA(ukuran_kanvas)

// 	draw.Draw(kanvas_kosong,ukuran_kanvas,kodeqr,image.ZP,draw.Over)

// 	draw.Draw(kanvas_kosong,gambar.Bounds().Add(image.Point{X: 10,Y: 10}),gambar,image.Point{},draw.Over)

// 	return kanvas_kosong
// }

// func matematika(a,b int) string{
// 	c := a + b
// 	cover := strconv.Itoa(c)
// 	return cover
// }

	//buat qr atau barcode nya
	// Barcode := qrcode.WriteFile("https://www.youtube.com/channel/UCc0QQoWjmqmxPCk0UQHmW5g",qrcode.Medium,512,"qr_ku_512.png")
	// // Barcode,err := qrcode.NewWithForcedVersion("konten saya", 25, qrcode.Medium)
	
	// if Barcode == nil {
	// 	fmt.Println("sukses untuk barcode",Barcode)
	// }else{
	// 	fmt.Println("tidak sukses")
	// }

// 	Barcode,err := qrcode.New("https://www.youtube.com/channel/UCc0QQoWjmqmxPCk0UQHmW5g",qrcode.Medium)

	

// 	Barcode.ForegroundColor = color.RGBA{
// 		R:255,
// 		G:0,
// 		B:0,
// 		A: 255,
// 	}
// 	Barcode.BackgroundColor = color.RGBA{
// 		R: 255,
// 		G: 255,
// 		B: 156,
// 		A: 255,
// 	}
// 	Barcode.DrawLabel("qr code ku ",20,10)

// 	err = Barcode.WriteFile(512,"qr_ku_ya3.png")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

/**
1. barcode
2. idestifikasi rest api golang
3. explore library http/
4. tambah image dan label
5. fungsi generate
*/