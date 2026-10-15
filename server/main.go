package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// klien adalah satu anggota grup yang sedang tersambung.
type klien struct {
	nama  string
	kirim chan string // antrean pesan keluar milik klien ini
}

// pesan adalah teks beserta pengirimnya.
type pesan struct {
	dari *klien
	teks string
}

// Channel komunikasi antara goroutine klien dan hub.
var (
	masuk  = make(chan *klien)
	keluar = make(chan *klien)
	siaran = make(chan pesan)
)

func main() {
	ln, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("gagal membuka port 9001: ", err)
	}
	defer ln.Close()

	log.Println("server mendengarkan di port 9001")

	go hub()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("gagal menerima sambungan:", err)
			continue
		}
		go layani(conn)
	}
}

// hub adalah satu-satunya pemilik daftar anggota. Karena hanya satu
// goroutine yang mengaksesnya, map tidak memerlukan mutex.
func hub() {
	anggota := make(map[*klien]bool)

	for {
		select {
		case k := <-masuk:
			anggota[k] = true
			sebar(anggota, nil, fmt.Sprintf("* %s bergabung (%d online)", k.nama, len(anggota)))

		case k := <-keluar:
			if anggota[k] {
				delete(anggota, k)
				close(k.kirim) // menghentikan goroutine penulis klien
				sebar(anggota, nil, fmt.Sprintf("* %s keluar (%d online)", k.nama, len(anggota)))
			}

		case p := <-siaran:
			sebar(anggota, p.dari, fmt.Sprintf("[%s] %s", p.dari.nama, p.teks))
		}
	}
}

// sebar mengirim teks ke seluruh anggota kecuali "kecuali" (boleh nil).
func sebar(anggota map[*klien]bool, kecuali *klien, teks string) {
	log.Println(teks)
	for k := range anggota {
		if k == kecuali {
			continue
		}
		select {
		case k.kirim <- teks:
		default:
			// Antrean penuh: pesan dilewati agar klien lambat tidak menahan grup.
		}
	}
}

// layani menangani satu koneksi sampai klien terputus.
func layani(conn net.Conn) {
	defer conn.Close()

	alamat := conn.RemoteAddr().String()
	log.Println("sambungan baru:", alamat)

	scanner := bufio.NewScanner(conn)

	fmt.Fprint(conn, "Nama kamu: ")
	if !scanner.Scan() {
		return
	}
	nama := strings.TrimSpace(scanner.Text())
	if nama == "" {
		nama = "tamu-" + alamat
	}

	k := &klien{nama: nama, kirim: make(chan string, 16)}

	masuk <- k
	defer func() { keluar <- k }()

	// Penulis: satu-satunya yang menulis ke koneksi setelah klien terdaftar.
	// Dijalankan setelah mendaftar agar tidak bocor jika klien putus lebih awal.
	go func() {
		for teks := range k.kirim {
			fmt.Fprintln(conn, teks)
		}
	}()

	for scanner.Scan() {
		teks := strings.TrimSpace(scanner.Text())
		if teks == "" {
			continue
		}
		siaran <- pesan{dari: k, teks: teks}
	}
}
