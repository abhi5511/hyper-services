## 📡 1. Hyper-SMS: High-Throughput Notification Engine
Ye service poore ecosystem ka "Messenger" hai. Iska kaam real-time alerts aur notifications ko handle karna hai.

*   **Core Architecture**: Apache Kafka par based asynchronous messaging engine.
*   **Workflow**:
    *   **Producers**: Koi bhi service (like Transit) Kafka topic par notification request bhejti hai.
    *   **Consumer Groups**: SMS service in messages ko consume karke priority ke hisaab se delivery queue mein daalti hai.
*   **Why Kafka?**: Taaki agar ek saath hazaron bus tracking alerts ya login OTPs aayein, toh system crash na ho (High Availability).

---

## 📦 2. Hyper-Storage: Hybrid Data & Media Node
Ye service files, images, aur database caching ko handle karti hai.

*   **Dual-Layer Storage**:
    1.  **Fast Cache (Redis)**: Frequently accessed data (like active user sessions or live bus coordinates) ko store karne ke liye.
    2.  **Persistent Storage (S3/Ceph)**: User documents, profile pictures, aur logs ko permanently save karne ke liye.
*   **Unified API**: Ek single endpoint jo decide karta hai ki data Redis mein jayega ya S3 mein, developer ko tension lene ki zaroorat nahi.

---

## 🏗️ 3. K8s-Configs: The Orchestration Brain
Ye folder poore **Hyper-realm** ko auto-pilot par chalata hai.

*   **Manifests**: Har service ke liye `Deployment` (scaling ke liye) aur `Service` (connectivity ke liye) YAML files.
*   **Security Management**:
    *   **Secrets**: RS256 Private/Public keys aur Database passwords ko environment variables ke bajaye K8s Secrets mein rakha gaya hai taaki GitHub par leak na hon.
    *   **ConfigMaps**: Services ki global settings (like Kafka brokers list) yahan manage hoti hain.
*   **Scalability**: `HorizontalPodAutoscaler` use karke hum traffic badhne par apne aap server badha sakte hain.

---

## 📂 Updated Folder Structure
Ab tere repo ka layout aisa dikhega:
```text
.
├── hyper-id/          # Identity & Auth Logic (Go + React)
├── sms-service/       # Kafka Consumers & SMS Gateways
├── storage-api/       # Redis Caching & S3 Adapters
├── k8s-configs/       # The "Infrastructure as Code" layer
│   ├── secrets/       # RSA & DB Secrets (Templated)
│   ├── deployments/   # Service Pod Definitions
│   └── services/      # Networking & Load Balancers
└── migrations/        # SQL & Redis Schema definitions
```

---
