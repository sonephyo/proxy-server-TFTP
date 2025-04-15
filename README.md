# Building proxy server using TFTP protocol

# 📄 Proxy Server Assignment Requirements (TFTP-over-TCP)

## ✅ Components

- [X] Client connects to proxy and sends a URL.
- [X] Proxy fetches file/page via HTTP and caches the most recent request.
- [X] Proxy relays file to client.
- [X] Only image files (e.g., .jpg) are required to be supported.

---

## ⚙️ Protocol Specifications

- [X] Use TCP instead of UDP (not traditional TFTP).
- [X] Extend TFTP (RFC 1350) where possible.
- [X] Use TFTP Option Extension (RFC 2347) if applicable.
- [X] Custom packet headers must be designed.
- [X] Support binary (octet) mode only.

---

## 🔐 Security and Session

- [X] Begin each session with:
  - [X] Sender ID
  - [X] Random number
- [X] Use both values to derive a shared session key.
- [X] Encrypt data using XOR or a better scheme.

---

## 📤 Transmission Protocol

- [X] Use TCP-style **sliding windows**
- [X] Implement TCP-style **Retransmission Timeout (RTO)** scheme.
- [ ] Add command-line option to simulate **1% packet drop**.

---

## 🧪 File Handling and Validation

- [ ] Received files should be stored in a **temporary directory** (e.g., `/tmp`).
- [ ] Validate file content with `cmp` or byte comparison.

---

## 📊 Throughput Report

Create a web page showing throughput under:
- [ ] At least **2 different host pairs**
- [ ] Window sizes: **1, 8, 64**
- [ ] With and without **1% simulated drop**