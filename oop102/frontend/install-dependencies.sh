#!/bin/bash
# install-dependencies.sh
# Quick setup script for the frontend

echo "☕ Coffee Shop Frontend - Setup Script"
echo "======================================="
echo ""

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install Node.js and npm first."
    exit 1
fi

echo "📦 Installing dependencies..."
npm install

if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies"
    exit 1
fi

echo ""
echo "✅ Dependencies installed!"
echo ""
echo "📝 Next steps:"
echo "1. Copy environment file:"
echo "   cp .env.example .env.local"
echo ""
echo "2. Edit .env.local and set API URL:"
echo "   REACT_APP_API_URL=http://localhost:8080"
echo ""
echo "3. Start the frontend:"
echo "   npm start"
echo ""
echo "4. Make sure backend is running:"
echo "   cd ../backend/cmd/api && go run main.go"
echo ""
echo "5. Open browser:"
echo "   http://localhost:3000"
echo ""
