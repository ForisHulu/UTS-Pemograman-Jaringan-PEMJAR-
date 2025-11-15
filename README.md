#UTS-Pemograman-Jaringan-PEMJAR-#

Proyek Game Tic-Tac-Toe

Proyek ini adalah implementasi dari game Tic-Tac-Toe online yang menggunakan model Client-Server dengan WebSocket dan sistem matchmaking berbasis peringkat (Rank/RR).

1. Anggota Kelompok :
  - Foris Juniawan Hulu (NIM: 223400015)
  - Daud Aldo S (NIM: 223400019)

2. Deskripsi Aplikasi

Aplikasi ini adalah game Tic-Tac-Toe multiplayer real-time yang dibangun dengan backend Go (Golang) dan frontend HTML/CSS/JavaScript. Komunikasi antara klien (browser) dan server dijembatani oleh WebSocket.

Fitur Utama:

Server Go: Server utama ditulis dalam Go, bertugas mengelola koneksi pemain, logika game, dan proses matchmaking.

Matchmaking Berbasis Peringkat (RR):

Pemain memiliki Poin Peringkat (Rank Rating / RR) yang disimpan di localStorage browser.

Server (Hub) mengelola satu antrian matchmaking dan akan menjodohkan dua pemain dengan perbedaan RR yang berada dalam batas toleransi (maxRRDifference).

Sistem Peringkat (Rank): RR pemain akan bertambah saat menang dan berkurang saat kalah (dengan floor di 100 RR). Peringkat memiliki nama (Bronze, Silver, Gold, dst.).

Komunikasi Real-time: Seluruh status game (pergerakan, giliran, status menang/kalah, rematch) dikirim dari server ke klien secara instan menggunakan WebSocket.


3. Petunjuk Cara Menjalankan Aplikasi

Prasyarat (Prerequisites)

Go (Golang) (versi 1.18 atau lebih baru)

Browser web modern (Chrome, Firefox, Safari, dll.)

Setup (Instalasi Dependensi)

Aplikasi ini memerlukan paket gorilla/websocket dari Go. Untuk menginstalnya, buka terminal dan jalankan:

go get [github.com/gorilla/websocket](https://github.com/gorilla/websocket)


Menjalankan Server (Run)

Pastikan ketiga file (main.go, index.html, style.css) berada dalam satu direktori yang sama.

Buka terminal di dalam direktori tersebut.

Jalankan server Go dengan perintah:

go run main.go


Jika berhasil, server akan berjalan dan Anda akan melihat pesan di terminal:

Server Matchmaking (v4.3 - FINAL FIX) dimulai di http://localhost:8080


Memulai Game (Akses Aplikasi)

Buka browser web Anda dan kunjungi alamat:
http://localhost:8080

Untuk menguji mode multiplayer, buka dua tab browser (atau dua browser yang berbeda, misal Chrome dan Safari) dan arahkan keduanya ke alamat http://localhost:8080.

Klik tombol "Cari Game" di kedua tab tersebut untuk memulai proses matchmaking.


4. Cuplikan Tampilan / Interaksi

Berikut adalah contoh alur interaksi aplikasi:

1. Lobi / Menu Utama
Pemain melihat status dan tombol "Cari Game". Jika pemain sudah pernah bermain, peringkat (RR) mereka dari sesi sebelumnya akan ditampilkan.

<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 48 03" src="https://github.com/user-attachments/assets/c7d58ebe-3a86-4da5-8544-0c3390536188" />
<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 50 20" src="https://github.com/user-attachments/assets/928db709-08c8-4962-a3c9-93d7aeaa958e" />

2. Proses Matchmaking
Setelah mengklik "Cari Game", server akan mencari lawan dengan RR yang sesuai.

<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 51 44" src="https://github.com/user-attachments/assets/0980d79b-06dc-473e-9bd1-6e72e2975b38" />

3. Game Berlangsung
Kedua pemain dijodohkan. Papan permainan aktif dan pemain bisa bergantian menempatkan 'X' atau 'O'. Status giliran ditampilkan.

<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 52 30" src="https://github.com/user-attachments/assets/069e4d85-b708-450b-bfb8-1bf4aa19085d" />
<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 53 01" src="https://github.com/user-attachments/assets/55daeb49-8d7f-4ad8-ae43-b729152b6e9d" />


4. Hasil Permainan & Update Peringkat
Setelah game selesai (Menang, Kalah, atau Seri), hasilnya akan ditampilkan. Poin RR pemain akan diperbarui dan ditampilkan di layar (misal: (+20) atau (-15)). Pemain kemudian dapat memilih untuk "Rematch" atau "Search for More".

<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 54 39" src="https://github.com/user-attachments/assets/c9b82b57-615d-4f4c-8a28-a7100a619de1" />
<img width="1680" height="929" alt="Screenshot 2025-11-16 at 01 55 03" src="https://github.com/user-attachments/assets/e47701b0-7a3c-4886-bb20-902abdac79e0" />

