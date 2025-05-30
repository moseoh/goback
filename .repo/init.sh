#! /bin/bash

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 현재 디렉토리가 저장소 루트인지 확인
if [ ! -d ".git" ]; then
  echo "❌ 현재 디렉토리가 저장소 루트가 아닙니다."
  exit 1
fi

cd $DIR/githook || exit 1
npm install --silent
