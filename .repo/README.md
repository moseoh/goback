# .repo

이 디렉토리는 GitHub 템플릿 레포지토리 초기화 및 관리를 위한 도구들을 포함하고 있습니다.

## 폴더 구조

### .repo/githook

Git 커밋 메시지 규칙을 강제하고 코드 품질을 유지하기 위한 Git hooks 구성 파일들을 포함합니다.

- **Husky**: Git hooks를 쉽게 설정하고 관리할 수 있도록 도와줍니다.
- **Commitlint**: 커밋 메시지가 [Conventional Commits](https://www.conventionalcommits.org/) 규칙을 따르도록 강제합니다.

설정 및 사용 방법은 해당 폴더의 README를 참조하세요.

### .repo/git

GitHub 저장소 관리 자동화 도구를 제공합니다.

- **브랜치 보호 설정**: 지정된 브랜치에 대한 보호 규칙을 자동으로 설정합니다.
- **자동 브랜치 삭제**: PR 병합 후 소스 브랜치를 자동으로 삭제하도록 설정합니다.
- **기타 GitHub 관련 자동화 도구**: 저장소 관리에 필요한 다양한 기능을 제공합니다.

각 도구의 상세 사용법은 해당 폴더의 문서를 참조하세요.

## 사용 방법

각 폴더 내의 README 파일과 스크립트 상단 주석을 참조하여 필요한 기능을 사용하세요.

## 기타

### 로컬 브랜치 정리 (주의: 병합 안된 브랜치도 삭제)

```bash
# ~/.zshrc
alias cleanb='git branch | grep -Ev "(^\*|master|develop|staging|release)" | xargs --no-run-if-empty git branch -D'
```
