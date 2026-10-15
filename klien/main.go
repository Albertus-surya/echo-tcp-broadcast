package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	alamat := "localhost:9001"
	if len(os.Args) > 1 {
		alamat = os.Args[1]
	}

	conn, err := net.Dial("tcp", alamat)
	if err != nil {
		fmt.Println("gagal tersambung ke", alamat, ":", err)
		return
	}
	defer conn.Close()

	fmt.Println("tersambung ke", alamat, "- ketik KELUAR untuk berhenti.")

	// Menampilkan semua pesan dari server. io.Copy dipakai agar prompt
	// tanpa baris baru tetap tampil.
	go func() {
		io.Copy(os.Stdout, conn)
		fmt.Println("\nserver menutup sambungan")
		os.Exit(0)
	}()

	// Mengirim ketikan pengguna ke server.
	keyboard := bufio.NewScanner(os.Stdin)
	for keyboard.Scan() {
		baris := strings.TrimSpace(keyboard.Text())
		if strings.ToUpper(baris) == "KELUAR" {
			return
		}
		if _, err := fmt.Fprintf(conn, "%s\n", baris); err != nil {
			fmt.Println("gagal mengirim:", err)
			return
		}
	}
}
