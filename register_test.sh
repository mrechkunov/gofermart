#!/bin/bash
curl -X POST http://localhost:8080/api/user/register -H 'Content-Type: application/json' -d '{"login":"mich","password":"retdsre"}'

