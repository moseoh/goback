/**
 * GitHub PR merge 방식을 Squash merge로 설정하는 스크립트
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
 * - 저장소 설정에서 PR merge 방식을 Squash merge로 설정합니다.
 * - 다른 merge 방식(merge commit, rebase merge)은 비활성화합니다.
 * - Squash merge 시 PR 제목과 설명을 유지하도록 설정합니다.
 */

import * as dotenv from "dotenv";
import { getGitRemoteInfo, createGitHubClient } from "../utils/github";
import { isConfigured, recordConfigured } from "../utils/gh-settings";

dotenv.config();

const featureName = "set-squash-merge";

/**
 * 저장소의 PR merge 방식을 Squash merge로 설정하는 함수
 */
export async function setSquashMergePreference() {
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

  console.log(`🔍 저장소 merge 설정 업데이트 중: ${owner}/${repo}`);
  
  try {
    // 저장소 설정 업데이트 - Squash merge 활성화, 다른 방식 비활성화
    await octokit.repos.update({
      owner,
      repo,
      // Squash merge만 활성화
      allow_squash_merge: true,
      allow_merge_commit: false,
      allow_rebase_merge: false,
      // Squash merge 시 PR 제목과 설명 사용
      use_squash_pr_title_as_default: true,
      squash_merge_commit_title: "PR_TITLE",
      squash_merge_commit_message: "PR_BODY"
    });
    
    console.log(`✅ PR merge 방식이 Squash merge로 설정되었습니다: ${owner}/${repo}`);
    console.log(`ℹ️ 설정 내용:`);
    console.log(`  - Squash merge: 활성화`);
    console.log(`  - Merge commit: 비활성화`);
    console.log(`  - Rebase merge: 비활성화`);
    console.log(`  - Squash merge 시 PR 제목과 설명을 유지`);
    
    recordConfigured(featureName);
  } catch (error) {
    console.error("❌ 저장소 설정 업데이트 중 오류 발생:", error);
    process.exit(1);
  }
}

// 직접 실행될 때 사용
if (require.main === module) {
  setSquashMergePreference();
} 