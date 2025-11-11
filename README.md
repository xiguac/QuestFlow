<div align="center">

  <img src="https://raw.githubusercontent.com/xiguac/QuestFlow/main/.github/assets/logo.png" alt="QuestFlow Logo" width="180"/>

  <h1>QuestFlow 🌊</h1>

  <p>
    <strong>A blazingly fast, highly extensible, and beautiful information collection platform, built with Go & Vue.</strong>
  </p>

  <p>
    From dynamic surveys 📝 and challenging exams 🧠 to simple team sign-ups 🤝, QuestFlow provides a powerful engine for all your data collection needs. It's engineered for performance, scalability, and an amazing developer experience.
  </p>

  <!-- Badges -->
  <p>
    <a href="https://github.com/xiguac/QuestFlow/actions/workflows/go.yml"><img src="https://img.shields.io/github/actions/workflow/status/xiguac/QuestFlow/go.yml?branch=main&style=for-the-badge&logo=githubactions&logoColor=white" alt="Backend CI"></a>
    <!-- Add frontend CI badge later -->
    <!-- <a href="#"><img src="https://img.shields.io/github/actions/workflow/status/xiguac/QuestFlow/node.js.yml?branch=main&label=Frontend%20CI&style=for-the-badge&logo=githubactions&logoColor=white" alt="Frontend CI"></a> -->
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"></a>
    <a href="https://vuejs.org/"><img src="https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" alt="Vue Version"></a>
    <a href="https://github.com/xiguac/QuestFlow/blob/main/LICENSE"><img src="https://img.shields.io/github/license/xiguac/QuestFlow?style=for-the-badge&color=blue" alt="License"></a>
    <img src="https://img.shields.io/badge/PRs-Welcome-brightgreen.svg?style=for-the-badge" alt="PRs Welcome">
  </p>
</div>

---

> ✨ **Project Status:** The backend is feature-complete, robust, and ready for production. The Vue 3 frontend is currently under active development. Join us and contribute!

## 🚀 Core Philosophy: Metadata-Driven Design

At the heart of QuestFlow lies a powerful concept: **everything is metadata**. A survey, an exam, or a sign-up form are not treated as separate entities. Instead, they are all just different representations of a flexible JSON `definition`.

This approach provides **unparalleled extensibility**. Want to add a new "file upload" question type? There's no need for complex database migrations. Simply teach the frontend how to render it and the backend how to validate it. The core system remains unchanged.

```json
{
  "title": "Your Awesome Form",
  "definition": {
    "type": "exam", // Could be "questionnaire", "registration", etc.
    "settings": { "timeLimit": 3600 },
    "questions": [
      // ... your array of question objects
    ]
  }
}
```

## 🌟 Key Features

### 🚄 High-Performance Architecture
QuestFlow is built to handle massive, spiky traffic without breaking a sweat.

*   **Asynchronous Submission Processing:** We use **Redis Streams** as a high-throughput message queue. When a user submits a form, the request is captured in milliseconds and the user gets an instant response.
*   **Decoupled Consumer Service:** A separate Go microservice (the consumer) reads from the message queue at a steady pace and writes to the database. This "shock absorption" system protects your database from being overwhelmed during peak traffic, ensuring system stability.
*   **Blazing Fast Go Backend:** Leveraging the power of Go's concurrency and the high-performance Gin framework, the API server is designed for low latency and a small memory footprint.

### 🎨 Powerful & Intuitive Form Builder (Frontend - In Progress)
Creating beautiful forms should be a joy, not a chore.

*   **Drag & Drop Interface:** Visually construct your forms by dragging components from a palette onto a canvas.
*   **Real-time Preview:** See exactly what your users will see as you build.
*   **Rich Question Types:** From basic multiple choice and text inputs to more advanced options, all fully customizable.
*   **Advanced Settings:** Configure submission deadlines, response limits, password protection, and more for each form.

### 🔐 Enterprise-Grade Security
Security is not an afterthought; it's a core design principle.

*   **Stateless JWT Authentication:** Secure and scalable authentication using JSON Web Tokens.
*   **Robust Authorization:** Clear separation of permissions. A user can *only* access and manage the forms they create.
*   **Built-in Protection:** Out-of-the-box defense against common web vulnerabilities, including SQL Injection (via GORM's prepared statements), XSS (proper output escaping), and CSRF (inherently mitigated by JWT Bearer token usage).

### 📊 Insightful Analytics
Turn raw data into actionable insights.

*   **Real-time Data Aggregation:** A dedicated statistics endpoint processes all submissions for a form and provides an aggregated summary.
*   **Beautiful Visualizations (Frontend - In Progress):** The frontend will use libraries like ECharts to render stunning, interactive charts (pie, bar, etc.) from the statistics data.
*   **Data Export (Planned):** Future support for exporting raw submission data to CSV/Excel for deeper analysis.

## 🛠️ Technology Stack

| Area          | Technology / Library                                       |
|:--------------|:-----------------------------------------------------------|
| **Backend**     | 🐹 **Go** (Language)                                       |
|               | 🌐 [Gin](https://github.com/gin-gonic/gin) (Web Framework)   |
|               | 🗃️ [GORM](https://gorm.io/) (ORM)                            |
|               | 🐬 [MySQL 8.0+](https://www.mysql.com/) (Database)            |
|               | ⚡️ [Redis Streams](https://redis.io/) (Message Queue)        |
|               | 🔑 [Go-JWT](https://github.com/golang-jwt/jwt) (Auth)          |
|               | ⚙️ [Viper](https://github.com/spf13/viper) (Configuration)     |
| **Frontend**    | 💚 **Vue 3** (Framework)                                   |
|               |  Vite (Build Tool)                                        |
|               | 🍍 [Pinia](https://pinia.vuejs.org/) (State Management)      |
|               | 🎨 [Element Plus](https://element-plus.org/) (UI Components) |
|               | 📊 [ECharts](https://echarts.apache.org/) (Data Visualization) |
|               | 📡 [Axios](https://axios-http.com/) (HTTP Client)            |

## 🚀 Getting Started

Get the QuestFlow backend up and running on your local machine in minutes.

### Prerequisites

*   [Go](https://go.dev/doc/install) (v1.18+)
*   [MySQL](https://dev.mysql.com/downloads/mysql/) (v8.0+)
*   [Redis](https://redis.io/topics/quickstart)

### 1. Clone the Repository
```bash
git clone https://github.com/xiguac/QuestFlow.git
cd QuestFlow
```

### 2. Configure Your Environment

1.  Create a fresh MySQL database:
    ```sql
    CREATE DATABASE questflow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    ```
2.  Copy the example configuration file:
    ```bash
    cp configs/config.example.yaml configs/config.yaml
    ```
3.  Edit `configs/config.yaml` with your database and Redis credentials. **Crucially, set a strong, unique `jwt.secret`!**

    ```yaml
    database:
      # Replace with your MySQL DSN (user:password@tcp(host:port)/dbname)
      dsn: "root:your_password@tcp(127.0.0.1:3306)/questflow?charset=utf8mb4&parseTime=True&loc=Local"
    
    redis:
      addr: "127.0.0.1:6379"
      password: "" # Your Redis password, if any

    jwt:
      # ❗️ IMPORTANT: Change this to a long, random, and secret string!
      secret: "replace-this-with-a-very-secure-random-string-!@#$%^"
    ```

### 3. Run the Backend Services

The backend consists of two services. You'll need two separate terminal windows to run them.

1.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```
2.  **Terminal 1: Start the API Server**
    This service handles all HTTP requests. Database tables are automatically migrated on the first run.
    ```bash
    go run ./cmd/server/
    ```

3.  **Terminal 2: Start the Submissions Consumer**
    This service processes the submission queue.
    ```bash
    go run ./cmd/consumer/
    ```

✅ **You're all set!** The QuestFlow backend is now live on `http://localhost:8080`.

## 📚 API Documentation

Our API is fully documented with detailed information on every endpoint, including request/response schemas and examples. This is the ultimate guide for frontend development or third-party integrations.

### **[➡️ Read the Full API Documentation](./API.md)**

## 🤝 How to Contribute

We welcome contributions of all kinds! Whether you're a Go guru, a Vue virtuoso, or just passionate about great software, we'd love your help.

1.  **Fork** the repository.
2.  Create your **Feature Branch** (`git checkout -b feature/AmazingFeature`).
3.  **Commit** your Changes (`git commit -m 'Add some AmazingFeature'`).
4.  **Push** to the Branch (`git push origin feature/AmazingFeature`).
5.  Open a **Pull Request**.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.