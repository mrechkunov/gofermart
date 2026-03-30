#!/bin/bash
curl -X POST http://localhost:8080/api/user/orders \
     -H "Content-Type: text/plain" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzQ2NDczMzEsImlhdCI6MTc3NDY0MDEzMSwidXNlcm5hbWUiOiJtaWNoIn0.xH7uyIB54RJr1Hyan40Pxm2IdzicLupYFgY2FV4r7VQ" \
     -d "577424872"