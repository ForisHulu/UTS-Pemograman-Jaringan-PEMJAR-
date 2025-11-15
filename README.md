# UTS-Pemograman-Jaringan-PEMJAR-

Proyek Game Tic-Tac-Toe

Proyek ini adalah implementasi dari game Tic-Tac-Toe online yang menggunakan model Client-Server dengan WebSocket dan sistem matchmaking berbasis peringkat (Rank/RR).

1. Anggota Kelompok
Daud Aldo S (NIM: 223400019)

Foris Juniawan Hulu (NIM: 223400015)

2. Deskripsi Aplikasi
Aplikasi ini adalah game Tic-Tac-Toe multiplayer real-time yang dibangun dengan backend Go (Golang) dan frontend HTML/CSS/JavaScript. Komunikasi antara klien (browser) dan server dijembatani oleh WebSocket.

Fitur Utama:

Server Go: Server utama ditulis dalam Go, bertugas mengelola koneksi pemain, logika game, dan proses matchmaking.

Matchmaking Berbasis Peringkat (RR):

Pemain memiliki Poin Peringkat (Rank Rating / RR) yang disimpan di localStorage browser.

Server (Hub) mengelola satu antrian matchmaking dan akan menjodohkan dua pemain dengan perbedaan RR yang berada dalam batas toleransi (maxRRDifference).

Sistem Peringkat (Rank): RR pemain akan bertambah saat menang dan berkurang saat kalah (dengan floor di 100 RR). Peringkat memiliki nama (Bronze, Silver, Gold, dst.).

Komunikasi Real-time: Seluruh status game (pergerakan, giliran, status menang/kalah, rematch) dikirim dari server ke klien secara instan menggunakan WebSocket.

4. Petunjuk Cara Menjalankan Aplikasi
Prasyarat (Prerequisites)

Go (Golang) (versi 1.18 atau lebih baru)

Browser web modern (Chrome, Firefox, Safari, dll.)

Setup (Instalasi Dependensi)

Aplikasi ini memerlukan paket gorilla/websocket dari Go. Untuk menginstalnya, buka terminal dan jalankan:

go get github.com/gorilla/websocket

Menjalankan Server (Run)

Pastikan ketiga file (main.go, index.html, style.css) berada dalam satu direktori yang sama.

Buka terminal di dalam direktori tersebut.

Jalankan server Go dengan perintah:

go run main.go

Jika berhasil, server akan berjalan dan Anda akan melihat pesan di terminal:

Server Matchmaking (v4.3 - FINAL FIX) dimulai di http://localhost:8080

Memulai Game (Akses Aplikasi)

Buka browser web Anda dan kunjungi alamat: http://localhost:8080

Untuk menguji mode multiplayer, buka dua tab browser (atau dua browser yang berbeda, misal Chrome dan Safari) dan arahkan keduanya ke alamat http://localhost:8080.

Klik tombol "Cari Game" di kedua tab tersebut untuk memulai proses matchmaking.

4. Cuplikan Tampilan / Interaksi
Berikut adalah contoh alur interaksi aplikasi:

1. Lobi / Menu Utama Pemain melihat status dan tombol "Cari Game". Jika pemain sudah pernah bermain, peringkat (RR) mereka dari sesi sebelumnya akan ditampilkan.
<img width="611" height="336" alt="image" src="https://github.com/user-attachments/assets/28c3fa79-ba0a-478b-b50e-55eb250aa6d9" />
<img width="608" height="334" alt="image" src="https://github.com/user-attachments/assets/07dbaf61-3184-4296-9542-2e4bb7062bca" />


2. Proses Matchmaking Setelah mengklik "Cari Game", server akan mencari lawan dengan RR yang sesuai.
<img width="608" height="334" alt="image" src="https://github.com/user-attachments/assets/b6366713-5f1c-4b56-9360-39c4ca173878" />

3. Game Berlangsung Kedua pemain dijodohkan. Papan permainan aktif dan pemain bisa bergantian menempatkan 'X' atau 'O'. Status giliran ditampilkan.
<img width="608" height="337" alt="image" src="https://github.com/user-attachments/assets/4304c4ca-3fc7-4e38-bc9c-1a278a14dbcc" />
<img width="610" height="331" alt="image" src="https://github.com/user-attachments/assets/ddd4b2fd-b196-499e-a5d5-eab841cbb18b" />


4. Hasil Permainan & Update Peringkat Setelah game selesai (Menang, Kalah, atau Seri), hasilnya akan ditampilkan. Poin RR pemain akan diperbarui dan ditampilkan di layar (misal: (+20) atau (-15)). Pemain kemudian dapat memilih untuk "Rematch" atau "Search for More".
<img width="607" height="335" alt="image" src="https://github.com/user-attachments/assets/9722279a-3fff-464a-90c2-93446372ee6f" />
<img width="604" height="332" alt="image" src="https://github.com/user-attachments/assets/4e4f96bb-4d58-4dca-918d-a5f94a7b7c9c" />


