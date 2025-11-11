# QuestFlow API Documentation 🌊

<p align="center">
  <em>Your one-stop guide to interacting with the QuestFlow backend.</em><br>
  <strong>Version: 1.0</strong>
</p>

---

Welcome to the official API documentation for QuestFlow! This document provides a comprehensive and developer-friendly guide to all available endpoints. Whether you're building a new client or integrating with our services, this is the place to start. 🚀

## 📜 Table of Contents

- [General Information](#-general-information)
  - [Base URL](#base-url)
  - [Authentication](#authentication)
  - [Standard Response Format](#standard-response-format)
  - [Common Error Codes](#common-error-codes)
- [👤 Auth Module](#-auth-module)
  - [Register a New User](#register-a-new-user)
  - [User Login](#user-login)
  - [Get Current User Info](#get-current-user-info)
- [📝 Form Module](#-form-module)
  - [Create a New Form](#create-a-new-form)
  - [Get Public Form Definition](#get-public-form-definition)
  - [Get Form Statistics](#get-form-statistics)
- [📩 Submission Module](#-submission-module)
  - [Submit Form Data](#submit-form-data)

---

## 🌐 General Information

### Base URL

All API endpoints are prefixed with `/api/v1`.

```
https://your-domain.com/api/v1
```

### Authentication

Endpoints that require authentication are marked with a 🛡️ icon.

- **Scheme:** `Bearer Token`
- **Header:** `Authorization: Bearer <YOUR_JWT_TOKEN>`

You can obtain a JWT token by successfully calling the `POST /users/login` endpoint. The client is expected to store this token and include it in the header for all subsequent requests to protected routes.

### Standard Response Format

We believe in consistency. All API responses, whether successful or not, follow this standard JSON structure:

```json
{
  "code": 0,
  "message": "success",
  "data": { ... } | [ ... ] | null
}
```

| Field     | Type           | Description                                                                                             |
| :-------- | :------------- | :------------------------------------------------------------------------------------------------------ |
| `code`    | `Integer`      | A business-specific status code. `0` indicates a successful operation. Non-zero values indicate an error. |
| `message` | `String`       | A human-readable message describing the result of the operation.                                        |
| `data`    | `Object|Array|null` | The payload of the response. This will be `null` for operations that don't return data or in case of an error. |

### Common Error Codes

Here are some common business error codes (`code` field) you might encounter:

| Code   | HTTP Status | Meaning                                      |
| :----- | :---------- | :------------------------------------------- |
| `4000` | `400`       | Bad Request - Invalid request parameters.    |
| `4001` | `401`       | Unauthorized - Invalid or missing token.     |
| `4003` | `403`       | Forbidden - Insufficient permissions.        |
| `4004` | `404`       | Not Found - The requested resource does not exist. |
| `5000` | `500`       | Internal Server Error - A generic server error. |
| `5001` | `500`       | Business Logic Error (e.g., username exists). |

---

## 👤 Auth Module

Endpoints for user registration, login, and session management.

### Register a New User

`POST /users/register`

Creates a new user account. On successful registration, the user is considered logged in, and a token is returned.

- **Authentication:** None
- **Request Body:** `application/json`

  | Field      | Type     | Constraints              | Description                               |
  | :--------- | :------- | :----------------------- | :---------------------------------------- |
  | `username` | `String` | **Required**, 3-20 chars | The unique username for the new account.  |
  | `password` | `String` | **Required**, 6-30 chars | The user's password.                      |
  | `email`    | `String` | Optional, valid email  | The user's email address (must be unique if provided). |

- **✅ Success Response (`200 OK`)**

  ```json
  {
    "code": 0,
    "message": "注册成功",
    "data": {
      "user_info": {
        "id": 123,
        "username": "newuser",
        "nickname": "",
        "email": "user@example.com",
        "role": 1
      },
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```

### User Login

`POST /users/login`

Authenticates a user with their credentials and returns a JWT token.

- **Authentication:** None
- **Request Body:** `application/json`

  | Field      | Type     | Constraints | Description                   |
  | :--------- | :------- | :---------- | :---------------------------- |
  | `username` | `String` | **Required**  | The user's registered username. |
  | `password` | `String` | **Required**  | The user's password.          |

- **✅ Success Response (`200 OK`)**

  ```json
  {
    "code": 0,
    "message": "登录成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user_info": {
        "id": 123,
        "username": "newuser",
        "nickname": "Tester",
        "role": 1
      }
    }
  }
  ```

### Get Current User Info

`GET /users/me` 🛡️

Retrieves the profile information of the currently authenticated user based on the provided JWT.

- **Authentication:** **Required**
- **✅ Success Response (`200 OK`)**

  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 123,
      "username": "newuser",
      "role": 1
    }
  }
  ```

---

## 📝 Form Module

Endpoints for creating, retrieving, and managing forms.

### Create a New Form

`POST /forms` 🛡️

Creates a new form. The entire structure, including questions and settings, is defined within the `definition` JSON object.

- **Authentication:** **Required**
- **Request Body:** `application/json`

  | Field         | Type          | Constraints  | Description                                        |
  | :------------ | :------------ | :----------- | :------------------------------------------------- |
  | `title`       | `String`      | **Required**   | The main title of the form.                        |
  | `description` | `String`      | Optional     | A longer description displayed below the title.    |
  | `definition`  | `Object`      | **Required**   | A JSON object defining the form's structure. See below. |

- **The `definition` Object:**

  ```json
  {
    "settings": {
      "type": "questionnaire" // "questionnaire", "exam", "registration", etc.
    },
    "questions": [
      {
        "id": "q1", // A unique ID for the question within the form
        "type": "single_choice", // e.g., "single_choice", "multi_choice", "text_input"
        "title": "What is your favorite programming language?",
        "required": true,
        "options": [
          { "id": "opt1", "text": "Go" },
          { "id": "opt2", "text": "Python" },
          { "id": "opt3", "text": "JavaScript" }
        ]
      }
      // ... more question objects
    ]
  }
  ```

- **✅ Success Response (`200 OK`)**

  Returns the complete form object, including the system-generated `ID` and public `FormKey`.

  ```json
  {
    "code": 0,
    "message": "创建成功",
    "data": {
      "ID": 5,
      "FormKey": "2yuOxSMnM",
      "CreatorID": 123,
      "Title": "Developer Survey",
      "Description": "...",
      "Definition": { /* ... definition object ... */ },
      "Status": 1, // 1: Draft, 2: Published, 3: Closed
      "CreatedAt": "2023-11-01T12:00:00Z",
      "UpdatedAt": "2023-11-01T12:00:00Z"
    }
  }
  ```

### Get Public Form Definition

`GET /public/forms/{form_key}`

Retrieves the public-facing definition of a form. This endpoint is used by the frontend to render the form for a user to fill out.

- **Authentication:** None
- **URL Parameter:**
  - `form_key` (string, **Required**): The unique public identifier for the form (e.g., `2yuOxSMnM`).
- **✅ Success Response (`200 OK`)**

  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "form_key": "2yuOxSMnM",
      "title": "Developer Survey",
      "description": "...",
      "definition": { /* ... definition object ... */ }
    }
  }
  ```

### Get Form Statistics

`GET /forms/{form_id}/stats` 🛡️

Retrieves aggregated statistics for a form's submissions. 🔐 **Only the creator of the form can access this endpoint.**

- **Authentication:** **Required**
- **URL Parameter:**
  - `form_id` (integer, **Required**): The internal numeric ID of the form (e.g., `5`).
- **✅ Success Response (`200 OK`)**

  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "total_submissions": 250,
      "question_stats": [
        {
          "question_id": "q1",
          "question_type": "single_choice",
          "title": "What is your favorite programming language?",
          "option_stats": [
            { "text": "Go", "count": 120 },
            { "text": "Python", "count": 80 },
            { "text": "JavaScript", "count": 50 }
          ]
        }
      ]
    }
  }
  ```

---

## 📩 Submission Module

Endpoints related to submitting form data.

### Submit Form Data

`POST /public/forms/{form_key}/submissions`

Submits a user's answers to a specific form. This is a high-performance, asynchronous endpoint designed to handle high traffic.

- **Authentication:** None
- **URL Parameter:**
  - `form_key` (string, **Required**): The unique public identifier for the form.
- **Request Body:** `application/json`

  The `data` object should contain `question_id` as keys and the user's answer(s) as values.
  - For `single_choice`: The value is the `option_id` string.
  - For `text_input`: The value is the text string.
  - For `multi_choice`: The value is an array of `option_id` strings.

  ```json
  {
    "data": {
      "q1": "opt1" // User selected "Go"
    }
  }
  ```
- **✅ Success Response (`200 OK`)**

  The response is returned **immediately** after the submission is accepted by the message queue for processing.
  ```json
  {
    "code": 0,
    "message": "提交成功，正在处理中...",
    "data": {
      "message_id": "1762266441332-0"
    }
  }
  ```

---
<p align="center">
  Happy coding! 🎉
</p>