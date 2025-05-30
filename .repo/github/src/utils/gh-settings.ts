/**
 * 기능 설정 및 설치 상태를 기록하는 간단한 유틸리티
 */

import { mkdirSync, writeFileSync, existsSync } from 'fs';
import { join } from 'path';

// 설치 기록을 저장할 디렉토리
const INSTALL_DIR = '.gh-settings';

// 설치 디렉토리 확인 및 생성
function ensureInstallDir(): void {
  if (!existsSync(INSTALL_DIR)) {
    mkdirSync(INSTALL_DIR, { recursive: true });
  }
}

/**
 * 기능 설치/실행 여부를 기록
 * @param featureName 기능 이름
 */
export function recordConfigured(featureName: string): void {
  ensureInstallDir();
  
  const timestamp = new Date().toISOString();
  const filePath = join(INSTALL_DIR, featureName);
  
  try {
    writeFileSync(filePath, `configured=${timestamp}\n`);
    console.log(`✅ github 설정 기록: ${featureName}`);
  } catch (error) {
    console.error(`❌ 설치 기록 중 오류 발생 (${featureName}):`, error);
  }
}

/**
 * 기능이 이미 설치/설정되었는지 확인
 * @param featureName 기능 이름
 * @returns 설치/설정 여부 (true/false)
 */
export function isConfigured(featureName: string): boolean {
  ensureInstallDir();
  const filePath = join(INSTALL_DIR, featureName);

  const exists = existsSync(filePath);
  if (exists) {
    console.log(`✅ github 설정 기록: ${featureName}`);
  }
  return exists;
}

// 초기화: 설치 디렉토리 확인
ensureInstallDir(); 