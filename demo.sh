#!/bin/bash

echo "=========================================="
echo "  DEVOPS CHALLENGE 2025 - DEMO"
echo "=========================================="

echo ""
echo ">>> release-service (cache: 10s)"
echo "--- 1ª chamada (esperado: MISS) ---"
curl -si http://localhost/release/ | grep -E "X-Cache-Status|x-cache-status"
echo "--- 2ª chamada imediata (esperado: HIT) ---"
curl -si http://localhost/release/ | grep -E "X-Cache-Status|x-cache-status"

echo ""
echo ">>> infra-service (cache: 60s)"
echo "--- 1ª chamada (esperado: MISS) ---"
curl -si http://localhost/infra/ | grep -E "X-Cache-Status|x-cache-status"
echo "--- 2ª chamada imediata (esperado: HIT) ---"
curl -si http://localhost/infra/ | grep -E "X-Cache-Status|x-cache-status"

echo ""
echo ">>> Aguardando 11s para expirar cache do release-service..."
sleep 11
echo "--- 3ª chamada após expirar (esperado: MISS ou EXPIRED) ---"
curl -si http://localhost/release/ | grep -E "X-Cache-Status|x-cache-status"

echo ""
echo ">>> infra-service ainda em cache (esperado: HIT)"
curl -si http://localhost/infra/ | grep -E "X-Cache-Status|x-cache-status"

echo ""
echo "=========================================="
echo "  FIM DA DEMO"
echo "=========================================="