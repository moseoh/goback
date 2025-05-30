# Git Hooks 설정

이 디렉토리는 Git 커밋 메시지 규칙을 강제하고 코드 품질을 유지하기 위한 Git hooks 구성 파일들을 포함하고 있습니다.

## 주요 기능

- **Husky**: Git hooks를 쉽게 설정하고 관리
- **Commitlint**: 커밋 메시지가 [Conventional Commits](https://www.conventionalcommits.org/) 규칙을 준수하도록 검증

## Getting Started

### 1. 패키지 설치

```bash
cd .repo/githook
npm install
```

### 2. Husky 설정

루트 디렉토리(프로젝트 최상위)에서 다음 명령어를 실행합니다:

```bash
# 루트 디렉토리에서 실행
cd .repo/githook
npx husky install ../../.husky
```

### 3. Commit-msg Hook 설정

```bash
# commit-msg 훅 설정
npx husky add ../../.husky/commit-msg "npx --no -- commitlint --edit $1"
```

## 커밋 메시지 규칙

이 프로젝트는 [Conventional Commits](https://www.conventionalcommits.org/) 규칙을 따릅니다:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 타입 목록

- `feat`: 새로운 기능 추가
- `fix`: 버그 수정
- `docs`: 문서 변경
- `style`: 코드 포맷팅, 세미콜론 누락 등 (코드 변경 없음)
- `refactor`: 코드 리팩토링
- `test`: 테스트 코드 추가 또는 수정
- `chore`: 빌드 프로세스, 패키지 매니저 설정 등 변경 (소스 코드 변경 없음)

### 예시

```
feat(auth): 로그인 기능 추가

사용자 로그인 API 및 폼 구현

Resolves: #123
```

## 문제 해결

### type="module" 오류

commitlint.config.js에서 ES 모듈 구문을 사용하는 경우, package.json에 다음을 추가합니다:

```json
{
  "type": "module"
}
```

### Husky가 작동하지 않을 때

Git hooks 디렉토리 권한을 확인합니다:

```bash
chmod -R +x .husky
```
