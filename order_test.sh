#!/bin/bash
curl -X POST http://localhost:8080/api/user/orders \
     -H "Content-Type: text/plain" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzUyMzYyNjUsImlhdCI6MTc3NTIyOTA2NSwidXNlcm5hbWUiOiJtaWNoIn0.YQpm-aK-BOpRi3KauMXJTKFxkW09dmMoPFcJ3B5khro" \
     -d "577424872"