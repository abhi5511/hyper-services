## 📡 Hyper-SMS: The Notification Nerve Center

**Hyper-SMS** ek high-throughput, asynchronous notification engine hai jo Apache Kafka ka use karke massive scale par messages handle karta hai. Iska main maqsad systems ke beech "Decoupling" create karna hai.

### 🛠️ Technical Stack
*   **Language**: Go (Golang) for high concurrency.
*   **Message Broker**: Apache Kafka & Zookeeper (For queuing and reliability).
*   **Protocol**: REST for Producers, Protobuf for internal messaging.
*   **Infrastructure**: Kubernetes (Scaled via HPA).

---

### 🏗️ Architecture Flow
System ka logic **Producer-Consumer** model par chalta hai:



1.  **The Producer (Trigger)**: Koi bhi service (like Hyper-ID ya Hyper-Transit) ek event generate karti hai (e.g., `USER_LOGIN_OTP` ya `BUS_DELAY_ALERT`).
2.  **The Queue (Kafka Topic)**: Ye event Kafka ke ek specific topic mein "Append" ho jata hai. Isse faida ye hai ki agar SMS service down bhi ho, toh messages loss nahi honge; wo queue mein wait karenge.
3.  **The Consumer (SMS Service)**: Hamari Go-based service Kafka se messages ko "consume" karti hai, usey format karti hai, aur gateway ke through bhej deti hai.

---

### 📂 Service Logic (Internal Working)

SMS service teen main parts mein divided hai:
*   **Dispatcher**: Ye decide karta hai ki message ki priority kya hai. OTPs "High Priority" queue mein jaate hain, jabki marketing updates "Low Priority" mein.
*   **Provider Manager**: Agar ek SMS gateway (jaise Twilio) down ho, toh ye automatically doosre gateway par switch kar jata hai.
*   **Rate Limiter**: Ye ensure karta hai ki system kisi ek user ko spam na kare ya gateway ki limits cross na ho jayein.

---

### ☸️ Kubernetes Deployment
K8s mein humne isey aise configure kiya hai:

*   **StatefulSets**: Kafka aur Zookeeper ke liye, taaki unka message data persist rahe.
*   **Horizontal Pod Autoscaler (HPA)**: Agar queue mein messages badhte hain, toh K8s apne aap SMS service ke extra "Pods" chalu kar deta hai.
*   **Health Checks**: `Liveness` aur `Readiness` probes ensure karti hain ki sirf wahi pods traffic lein jo Kafka se successfully connected hain.

---

### 📝 Key API Example (Internal)
Producer is tarah ka payload Kafka topic mein push karta hai:

```json
{
  "event_id": "uuid-v4",
  "service_origin": "hyper-id",
  "recipient": "+91XXXXXXXXXX",
  "message_type": "otp",
  "payload": {
    "code": "5511",
    "expiry": "5m"
  }
}
```


---

