# Receipt Processor

Receipt Processor is a web service designed to process receipt data and calculate reward points based on specific business rules. This lightweight application is implemented in Go and uses an in-memory store to manage data.

---

## **Table of Contents**
- [Overview](#overview)
- [API Endpoints](#api-endpoints)
- [Validation Rules](#validation-rules)
- [Points Calculation Rules](#points-calculation-rules)
- [Setup and Run Instructions](#setup-and-run-instructions)
- [Examples](#examples)

---

## **Overview**

This service processes receipts to generate a unique ID for each receipt, calculates reward points based on predefined rules, and allows retrieval of the calculated points. It is built using the Gin framework and supports only in-memory storage.

---

## **API Endpoints**

### **1. Process Receipts**
- **Endpoint**: `/receipts/process`
- **Method**: `POST`
- **Description**: Submits a receipt for processing and returns a unique receipt ID.
- **Request Body**:
  ```json
  {
    "retailer": "M&M Corner Market",
    "purchaseDate": "2022-03-20",
    "purchaseTime": "14:33",
    "items": [
      { "shortDescription": "Gatorade", "price": "2.25" }
    ],
    "total": "9.00"
  }
  ```
- **Response**:
  ```json
  { "id": "adb6b560-0eef-42bc-9d16-df48f30e89b2" }
  ```
- **Error**:
  ```json
  { "description": "The receipt is invalid." }
  ```

---

### **2. Get Points**
- **Endpoint**: `/receipts/{id}/points`
- **Method**: `GET`
- **Description**: Retrieves the reward points for a given receipt ID.
- **Path Parameter**:
  - **`id`**: The unique ID assigned to the receipt.
- **Response**:
  ```json
  { "points": 109 }
  ```
- **Error**:
  ```json
  { "description": "No receipt found for that ID." }
  ```

---

## **Validation Rules**

### **Required Fields**
Each receipt must include the following fields:
- **`retailer`** (string): The name of the retailer or store.
- **`purchaseDate`** (string): The date of the purchase (format: `YYYY-MM-DD`).
- **`purchaseTime`** (string): The time of the purchase (format: `HH:MM` in 24-hour format).
- **`items`** (array): A non-empty list of items purchased.
  - Each item must include:
    - **`shortDescription`** (string): Description of the item.
    - **`price`** (string): The price of the item (format: `X.XX`).
- **`total`** (string): The total amount paid on the receipt (format: `X.XX`).

### **Validation Rules**
#### **General Rules**
- All fields must be present and adhere to their specified formats.

#### **Field-Specific Rules**
1. **`retailer`**:
   - Must not be empty.
   - Can contain alphanumeric characters, spaces, hyphens (`-`), and ampersands (`&`).

2. **`purchaseDate`**:
   - Must be in `YYYY-MM-DD` format.
   - Must be a valid calendar date (e.g., February cannot have 30 days).
   - Leap years are validated (e.g., February 29 is valid only in leap years).
   - Cannot be a future date.

3. **`purchaseTime`**:
   - Must be in `HH:MM` format (24-hour time).
   - If `purchaseDate` is today, `purchaseTime` cannot be in the future.

4. **`items`**:
   - Must contain at least one item.
   - Each item's `shortDescription` must not be empty and can only include alphanumeric characters, spaces, and hyphens.
   - Each item's `price` must be in the format `X.XX` and greater than `0`.

5. **`total`**:
   - Must be in the format `X.XX`.
   - Must equal the sum of the prices of all items.

6. **Path Parameter (`id`)**:
   - Must consist of non-whitespace characters (`^\S+$`).

---

## **Points Calculation Rules**

Points are awarded based on the following rules:

1. **Retailer Name**:
   - 1 point for every alphanumeric character in the retailer name.

2. **Total**:
   - 50 points if the total is a round dollar amount (e.g., `5.00`).
   - 25 points if the total is a multiple of `0.25`.

3. **Items**:
   - 5 points for every two items.
   - If the trimmed length of the item's description is a multiple of `3`, multiply the item price by `0.2` and round up.

4. **Purchase Date**:
   - 6 points if the purchase day is odd.

5. **Purchase Time**:
   - 10 points if the time is after 2:00 PM and before 4:00 PM.

---

## **Setup and Run Instructions**

### **Prerequisites**
- Go 1.19+ installed on your system.

### **Steps to Run**
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd receipt-processor
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

4. Access the API at `http://localhost:8080`.

---

## **Examples**

### **Processing a Receipt**
**Request**:
```http
POST /receipts/process HTTP/1.1
Content-Type: application/json

{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    { "shortDescription": "Mountain Dew 12PK", "price": "6.49" }
  ],
  "total": "6.49"
}
```

**Response**:
```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

---

### **Retrieving Points**
**Request**:
```http
GET /receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points HTTP/1.1
```

**Response**:
```json
{ "points": 32 }
```

---

## **Notes**
- Data is stored in memory and will reset when the application restarts.
- Refer to `api.yml` for a formal API definition.

---
