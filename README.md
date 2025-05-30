## Pulumi dan Golang untuk Automasi Deploy Nomad di DigitalOcean

![image](https://github.com/user-attachments/assets/73bb763f-ed5a-450b-b417-6810845eb609)

Tulisan ini menjelaskan langkah-langkah untuk membuat dan mengelola cluster Nomad menggunakan Pulumi (Go) dan DigitalOcean.

---

### 📦 Pre-Requisites

Pastikan kamu sudah menginstal:

* [Pulumi](https://www.pulumi.com/docs/get-started/install/)
* [Go](https://go.dev/dl/)
* [Nomad CLI](https://developer.hashicorp.com/nomad/downloads)
* Akses akun [DigitalOcean](https://www.digitalocean.com/) + token API aktif
* `ssh-keygen` tersedia (biasanya sudah ada di Linux/macOS)

---

### 🛠️ Langkah-Langkah Deploy

1. **Inisialisasi Project Pulumi (Go)**

   ```bash
   pulumi new go
   ```

   Ikuti instruksi untuk memilih nama project, deskripsi, dan stack.
2. **Apply Token DigitalOcean ke Pulumi**

   ```bash
   pulumi config set digitalocean:token {DIGITAL_OCEAN_TOKEN} --secret
   ```
3. **Buat SSH Key untuk akses Droplet**

   ```bash
   ssh-keygen -f dostart-keypair
   ```
4. **Generate TLS Certificate untuk Nomad**

   ```bash
   mkdir cert
   cd cert
   nomad tls cert create -server -region global
   cd ..
   ```
5. **Tambahkan Path SSH Key ke Konfigurasi Pulumi**

   ```bash
   pulumi config set publicKeyPath dostart-keypair.pub
   pulumi config set privateKeyPath dostart-keypair.key
   ```
6. **Deploy Droplet**
   Jalankan:

   ```bash
   pulumi up
   ```

   Setelah sukses, ambil IP Droplet:

   ```bash
   pulumi stack output dropletIP
   ```

---

### 🔒 Akses Nomad via TLS

1. **Cek status node Nomad dari lokal menggunakan TLS cert:**

   ```bash
   nomad node status \
     -ca-cert=nomad-agent-ca.pem \
     -client-cert=global-cli-nomad.pem \
     -client-key=global-cli-nomad-key.pem \
     -address=https://{IP_ADDRESS}:4646
   ```

   Ganti `{IP_ADDRESS}` dengan IP hasil `pulumi stack output dropletIP`.
2. **Set Nomad environment variable (jika sudah terhubung melalui SSH atau port forwarding):**

   ```bash
   export NOMAD_ADDR=https://127.0.0.1:4646
   ```
3. **Cek status layanan di Nomad**

   ```bash
   nomad status
   ```

---

### 🛡️ Bootstrap ACL (Access Control List)

1. **Inisialisasi token ACL**

   ```bash
   nomad acl bootstrap
   ```
2. **Set token ke environment variable:**

   ```bash
   export NOMAD_TOKEN="{TOKEN}"
   ```

   ⚠️ Ganti token sesuai hasil `nomad acl bootstrap`.
3. Akses nomad server via browser

   ```
   https://{IP_ADDRESS}}:4646
   ```

---

### 📚 Catatan Tambahan

* File `.pem` yang dihasilkan dari `nomad tls cert` harus disimpan dengan aman.
* Gunakan `ssh -i dostart-keypair.key root@{IP}` untuk mengakses Droplet secara manual.
* Jika ingin menghancurkan resource:
  ```bash
  pulumi destroy
  ```
