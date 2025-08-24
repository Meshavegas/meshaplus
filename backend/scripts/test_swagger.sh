#!/bin/bash

# Script pour tester la documentation Swagger
echo "🧪 Test de la documentation Swagger..."

# Vérifier si le serveur est en cours d'exécution
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Le serveur n'est pas en cours d'exécution sur localhost:8080"
    echo "💡 Démarrez le serveur avec : ./main"
    exit 1
fi

echo "✅ Serveur détecté sur localhost:8080"

# Tester l'endpoint de santé
echo "🔍 Test de l'endpoint de santé..."
if curl -s http://localhost:8080/health | grep -q "success"; then
    echo "✅ Endpoint de santé fonctionne"
else
    echo "❌ Endpoint de santé ne fonctionne pas"
fi

# Tester l'endpoint Swagger JSON
echo "🔍 Test de l'endpoint Swagger JSON..."
if curl -s http://localhost:8080/swagger/doc.json | grep -q "swagger"; then
    echo "✅ Endpoint Swagger JSON fonctionne"
else
    echo "❌ Endpoint Swagger JSON ne fonctionne pas"
fi

# Tester l'interface Swagger UI
echo "🔍 Test de l'interface Swagger UI..."
if curl -s http://localhost:8080/swagger/index.html | grep -q "swagger-ui"; then
    echo "✅ Interface Swagger UI fonctionne"
else
    echo "❌ Interface Swagger UI ne fonctionne pas"
fi

echo ""
echo "🌐 Accédez à la documentation :"
echo "   - Interface Swagger UI : http://localhost:8080/swagger/index.html"
echo "   - JSON Swagger : http://localhost:8080/swagger/doc.json"
echo "   - YAML Swagger : http://localhost:8080/swagger/doc.yaml"
echo ""
echo "📊 Si tout fonctionne, vous devriez voir la documentation complète de votre API !" 