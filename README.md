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

> ✨ **Project Status:** This project is actively under development. The core backend architecture is stable and functional. The Vue 3 frontend is evolving rapidly. Join us on this exciting journey!

## 🚀 Core Philosophy: Metadata-Driven Design

At the heart of QuestFlow lies a powerful concept: **everything is metadata**. A survey, an exam, or a sign-up form are not treated as separate entities. Instead, they are all just different representations of a flexible JSON `definition`.

This approach provides **unparalleled extensibility**. Want to add a new "file upload" question type? There's no need for complex database migrations. Simply teach the frontend how to render it and the backend how to validate it. The core system remains unchanged.

```json
{
  "title": "Your Awesome Form",
  "description": "A brief description here.",
  "definition": {
    "settings": { "type": "exam", "timeLimit": 3600 },
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

### 🎨 Powerful & Intuitive Form Builder (Frontend)
Creating beautiful forms should be a joy, not a chore.

*   **Drag & Drop Interface:** Visually construct your forms by dragging components from a palette onto a canvas.
*   **Real-time Preview:** See exactly what your users will see as you build.
*   **Rich Question Types:** From basic multiple choice and text inputs to more advanced options, all fully customizable.
*   **Advanced Settings:** Configure submission deadlines, response limits, and more for each form.

### 🔐 Enterprise-Grade Security
Security is not an afterthought; it's a core design principle.

*   **Stateless JWT Authentication:** Secure and scalable authentication using JSON Web Tokens.
*   **Robust Authorization:** Clear separation of permissions. A user can *only* access and manage the forms they create.
*   **Environment-based Secrets:** No hardcoded passwords or keys. All sensitive information is loaded from environment variables, following best practices for security.

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
|               | 🚀 [Vite](https://vitejs.dev/) (Build Tool)                                        |
|               | 🍍 [Pinia](https://pinia.vuejs.org/) (State Management)      |
|               | 🎨 [Element Plus](https://element-plus.org/) (UI Components) |
|               | 📊 *ECharts (Planned)* (Data Visualization) |
|               | 📡 [Axios](https://axios-http.com/) (HTTP Client)            |

## 🚀 Getting Started

Get the full QuestFlow stack up and running on your local machine.

### Prerequisites

*   [Go](https://go.dev/doc/install) (v1.18+)
*   [Node.js](https://nodejs.org/) (v18+)
*   [MySQL](https://dev.mysql.com/downloads/mysql/) (v8.0+)
*   [Redis](https://redis.io/topics/quickstart)

### 1. Clone the Repository
```bash
git clone https://github.com/xiguac/QuestFlow.git
cd QuestFlow
```

### 2. Configure Your Environment

1.  **Create a fresh MySQL database** and a dedicated user:
    ```sql
    CREATE DATABASE questflow CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
    CREATE USER 'questflow_user'@'localhost' IDENTIFIED BY 'your_strong_password';
    GRANT ALL PRIVILEGES ON questflow.* TO 'questflow_user'@'localhost';
    FLUSH PRIVILEGES;
    ```
2.  **Set up Backend Configuration:**
    *   Copy the example config: `cp configs/config.example.yaml configs/config.yaml`
    *   Create a local environment file (this is ignored by Git): `touch .env`
    *   Add your secrets to the `.env` file. This file will override any empty values in `config.yaml`.

    **File: `.env`**
    ```dotenv
    # Database password
    DATABASE_PASSWORD=your_strong_password

    # ❗️ IMPORTANT: Generate a long, random, and secret string for JWT!
    JWT_SECRET=replace-this-with-a-very-secure-random-string-!@#$%^
    ```
3.  **Review `configs/config.yaml`:**
    Ensure the database `user`, `host`, `dbname`, and Redis settings match your local environment. The `password` and `secret` fields should be left empty as they are loaded from `.env`.

### 3. Run the Services

You'll need two separate terminals for the backend and one for the frontend.

1.  **Terminal 1: Start the Backend API Server**
    ```bash
    # Install Go dependencies
    go mod tidy
    
    # Run the server (it will also auto-migrate the database)
    go run ./cmd/server/main.go
    ```
    The API is now running at `http://localhost:8080`.

2.  **Terminal 2: Start the Frontend Dev Server**
    ```bash
    cd frontend
    npm install
    npm run dev
    ```
    The web interface is now available at `http://localhost:5173`.

✅ **You're all set!** Open your browser to `http://localhost:5173` and start exploring QuestFlow.

## 📚 API Documentation

Our API is fully documented with detailed information on every endpoint, including request/response schemas and examples. This is the ultimate guide for frontend development or third-party integrations.

### **[➡️ Read the Full API Documentation](./API.md)**

## 🤝 How to Contribute

We welcome contributions of all kinds with open arms! Whether you're a Go guru, a Vue virtuoso, or just passionate about great software, we'd love your help.

1.  **Fork** the repository.
2.  Create your **Feature Branch** (`git checkout -b feature/AmazingFeature`).
3.  **Commit** your Changes (`git commit -m 'feat: Add some AmazingFeature'`).
4.  **Push** to the Branch (`git push origin feature/AmazingFeature`).
5.  Open a **Pull Request**.

Check out our [open issues](https://github.com/xiguac/QuestFlow/issues) to find a good place to start!

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

> **Note on AI-Generated Comments**: Some of the code comments within this project may have been partially generated by AI assistants (like Microsoft Copilot or GitHub Copilot) for expediency. While they are reviewed for correctness, they may not always be perfectly accurate or reflect the full intent of the code. Always trust the code itself as the primary source of truth.
