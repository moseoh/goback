/**
 * GitHub PR 머지 후 브랜치 자동 삭제 옵션 활성화 스크립트
 * 
 * 필요한 권한:
 * - GitHub 개인 액세스 토큰(PAT)이 필요합니다.
 * - 토큰에는 최소한 'repo' 스코프 권한이 필요합니다.
 * 
 * 환경 변수 설정 방법:
 * .env 파일에 다음과 같이 설정합니다:
 * 
 * ```
 * # GitHub 개인 액세스 토큰
 * GITHUB_TOKEN=ghp_your_personal_access_token
 * ```
 * 
 * 기능:
 * - 저장소 설정에서 '자동으로 PR 소스 브랜치 삭제' 옵션을 활성화합니다.
 * - 이 설정을 활성화하면 PR이 머지될 때 소스 브랜치가 자동으로 삭제됩니다.
 */

import * as dotenv from "dotenv";
import { getGitRemoteInfo, createGitHubClient } from "../utils/github";
import { isConfigured, recordConfigured } from "../utils/gh-settings";

dotenv.config();

const featureName = "auto-delete-branch";

/**
 * 저장소에 대한 '자동으로 병합된 브랜치 삭제' 옵션을 활성화합니다.
 */
export async function enableAutoDeleteMergedBranches() {
  if (isConfigured(featureName)) {
    return;
  }

  const token = process.env.GITHUB_TOKEN!;

  if (!token) {
    console.error("❌ GITHUB_TOKEN이 필요합니다.");
    process.exit(1);
  }

  const { owner, repo } = getGitRemoteInfo();
  const octokit = createGitHubClient(token);

  console.log(`🔍 저장소 설정 확인 중: ${owner}/${repo}`);
  
  try {
    // 저장소 설정 업데이트 - 병합된 브랜치 자동 삭제 옵션 활성화
    await octokit.repos.update({
      owner,
      repo,
      delete_branch_on_merge: true
    });
    
    console.log(`✅ PR 병합 시 브랜치 자동 삭제 옵션이 활성화되었습니다: ${owner}/${repo}`);
    console.log(`ℹ️ 이제부터 PR이 머지되면 소스 브랜치가 자동으로 삭제됩니다.`);
    
    recordConfigured(featureName);
  } catch (error) {
    console.error("❌ 저장소 설정 업데이트 중 오류 발생:", error);
    process.exit(1);
  }
}

// 직접 실행될 때 사용
if (require.main === module) {
  enableAutoDeleteMergedBranches();
} 