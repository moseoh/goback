# GitHub 저장소 관리 도구

이 디렉토리는 GitHub 저장소 관리와 관련된 자동화 스크립트와 설정 파일을 포함하고 있습니다.

## 주요 기능

- **브랜치 보호 설정**: 지정된 브랜치에 대한 보호 규칙을 자동으로 설정
- **자동 브랜치 삭제**: PR 병합 후 소스 브랜치를 자동으로 삭제하도록 설정
- **기타 GitHub 관련 자동화 도구**: 저장소 관리에 필요한 다양한 기능 제공

## Getting Started

### 1. 패키지 설치

```bash
cd .repo/github
npm install
```

### 2. 환경 변수 설정

`.env` 파일을 생성하고 필요한 환경 변수를 설정합니다:

```bash
# .env 파일 생성
cat > .env << EOL
# GitHub 개인 액세스 토큰 (PAT)
GITHUB_TOKEN=ghp_your_personal_access_token

# 보호할 브랜치 목록 (쉼표로 구분)
PROTECTED_BRANCHES=main,develop,staging,release
EOL
```

GitHub Personal Access Token은 다음 스코프 권한이 필요합니다:
- `repo`: 저장소 관리 권한
- `admin:org` (조직 레포지토리의 경우)

### 3. 브랜치 보호 규칙 설정

```bash
# 환경 변수에 명시된 브랜치들에 보호 규칙 적용
npm run protect
```

이 명령은 다음 작업을 수행합니다:
- 존재하지 않는 브랜치는 자동 생성
- PR에 최소 1개의 승인이 필요하도록 설정
- 관리자(owner)는 제한 없이 변경할 수 있도록 설정

### 4. PR 병합 후 브랜치 자동 삭제 설정

```bash
# PR 병합 시 소스 브랜치 자동 삭제 옵션 활성화
npm run auto-delete
```

### 5. 모든 설정 한 번에 적용

```bash
# 브랜치 보호 규칙 및 자동 삭제 옵션 모두 설정
npm run start all
```

## 고급 기능

### 브랜치 자동 삭제 상태 확인

```bash
npm run check-auto-delete
```

### 도움말 보기

```bash
npm run start help
```

## API 문서

모든 스크립트는 TypeScript로 작성되었으며, 주요 기능들을 모듈화하여 제공합니다:

- `protectBranches()`: 브랜치 보호 규칙 설정
- `enableAutoDeleteMergedBranches()`: PR 병합 시 브랜치 자동 삭제 옵션 활성화
- `checkAutoDeleteMergedBranchesStatus()`: 자동 삭제 옵션 상태 확인

## 문제 해결

### 권한 오류

권한 오류가 발생하는 경우, GitHub 토큰에 적절한 권한이 부여되었는지 확인하세요.

### 브랜치 보호 규칙이 적용되지 않는 경우

저장소 소유자 또는 관리자 권한이 있는지 확인하세요. 브랜치 보호 규칙을 설정하려면 관리자 권한이 필요합니다.
