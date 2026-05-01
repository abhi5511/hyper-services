## 🏗️ K8s-Configs: The Orchestration & Automation Brain

**K8s-Configs** folder mein wo saari instruction files (YAML manifests) hain jo **Hyper-realm** infrastructure ko automate karti hain. Ye "Infrastructure as Code" (IaC) approach follow karta hai, taaki tu kisi bhi cloud (AWS, Azure) ya local cluster (Minikube) par ek command se poora setup khada kar sake.



### 🛠️ Core Components
*   **Orchestrator**: Kubernetes (Minikube for local dev).
*   **Container Runtime**: Docker.
*   **Scaling**: Horizontal Pod Autoscaler (HPA).
*   **Security Management**: K8s Secrets & ConfigMaps.

---

### 🛡️ 1. Security & Secrets Management
System ki sensitive information ko humne images mein hardcode nahi kiya hai.
*   **Identity Secrets**: RS256 Private/Public keys ko `hyper-id-certs.yaml` mein store kiya gaya hai taaki JWT tokens secure rahein.
*   **Database Credentials**: Postgres, Redis, aur Kafka ke passwords `hyper-id-db-secrets.yaml` ke through inject hote hain.
*   **Environment Injection**: Saari microservices in secrets ko environment variables ki tarah use karti hain, jisse leakage ka khatra khatam ho jata hai.

---

### 📡 2. Networking & Service Discovery
Microservices ek doosre se kaise baat karti hain, ye yahan define hota hai:
*   **ClusterIP**: Internal services (jaise `id-db` ya `redis`) sirf cluster ke andar hi accessible hain, taaki bahar se koi attack na kar sake.
*   **NodePort/LoadBalancer**: `hyper-id-service` ko port `30085` par expose kiya gaya hai taaki bahar ki duniya (browsers) usey access kar sake.
*   **Service Discovery**: Kisi bhi service ko uska naam use karke (e.g., `http://hyper-id-db-service:5432`) cluster ke andar dhunda ja sakta hai.

---

### 🔄 3. Self-Healing & Scalability
*   **Replication**: Har service ke multiple "Replicas" (copies) chal rahi hoti hain. Agar ek pod fail hota hai, toh K8s usey turant restart kar deta hai.
*   **Liveness & Readiness Probes**: K8s lagatar check karta hai ki kya app ready hai traffic lene ke liye. Agar app "hang" ho jaye, toh pod kill karke naya create kiya jata hai.
*   **Persistent Storage**: `PersistentVolumeClaims` (PVC) ka use karke humne database ka data "Immutable" rakha hai, yani pod delete hone par bhi data delete nahi hoga.

---

### 📂 Folder Hierarchy (Inside `k8s-configs/`)
```text
.
├── secrets/          # RSA Keys, DB Passwords, API Tokens
├── deployments/      # Pod definitions for ID, SMS, Storage
├── services/         # Internal & External networking rules
└── volumes/          # Storage disk definitions for Postgres & Redis
```

---

### 🚀 Master Command to Launch Ecosystem
Poore Hyper-realm ko start karne ke liye sirf ye kaafi hai:
```bash
kubectl apply -f k8s-configs/
```

---

