# TCP Broadcast Server di Go

Server chat grup sederhana: pesan dari satu klien disiarkan ke klien lain.

## Menjalankan

```
go run ./server    # terminal 1 (port 9001)
go run ./klien     # terminal 2, 3, dst. Isi nama lalu ketik pesan
```

## Cara kerja

- **TCP**: tiap klien membuka satu koneksi ke server, pesan dikirim per baris.
- **Goroutine**: server melayani setiap klien di goroutine sendiri.
- **Hub**: satu goroutine yang menyimpan daftar anggota dan menyebarkan pesan.
- **Channel**: klien menitipkan pesan ke hub lewat channel, lalu hub menaruhnya di antrean tiap klien.