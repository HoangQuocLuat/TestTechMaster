Bài 1
- cách chạy
b1. tạo file env và cáu hình
GROQ_API_KEY=
GROQ_URL_CHAT=https://api.groq.com/openai/v1/chat/completions
WEB_PORT=8080
b2 truy cập cmd chạy go run main trên terminal
b3 truy cập folder fe và chạy live server
- Kết quả
![image](https://github.com/user-attachments/assets/1b4382e1-27fa-44b9-825f-c7fceb6663b4)

Bài 2
- cách chạy
truy cập folder fe và chạy live server
- Kết quả
-  ![image](https://github.com/user-attachments/assets/519aa195-4e0a-4888-aba3-94701fffa8a3)

Bài 3
- cách chạy
b1. tạo file env và cáu hình
GROQ_API_KEY=
GROQ_URL_CHAT=https://api.groq.com/openai/v1/chat/completions
WEB_PORT=8080
DATABASE_USERNAME=root
DATABASE_PASSWORD=123
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=golang_test
DATABASE_POOL_IDLE=10
DATABASE_POOL_MAX=100
DATABASE_POOL_LIFETIME=300
LOG_LEVEL=6
b2 truy cập be/cmd và chạy go run main.go trên terminal
- Kết quả
![image](https://github.com/user-attachments/assets/33dae3f9-6a85-464b-80fb-8bc91c1ef92c)
![image](https://github.com/user-attachments/assets/2877cac5-81a0-4d98-9259-c5144a411f1c)

![image](https://github.com/user-attachments/assets/76940176-9e48-4d10-aeff-b07f6fb302d2)
![image](https://github.com/user-attachments/assets/6a8ffbb8-6049-4d8d-8177-5ad58f7157e7)
![image](https://github.com/user-attachments/assets/94c89cd0-e4b9-48d2-b6ff-b5ba8470f480)

các từ đã được lưu trong database và mỗi giai đoạn sẽ được thông báo qua websocket
