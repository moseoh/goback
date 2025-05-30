/**
 * GitHub 브랜치 보호 규칙 설정 스크립트
 * 
 * 필요한 권한:
 * - GitHub 개인 액세스 토큰(PAT)이 필요합니다.
 * - 토큰에는 최소한 'repo' 스코프 권한이 필요합니다.
 * - 조직 레포지토리의 경우 'repo'와 'admin:org' 스코프가 필요할 수 있습니다.
 * - GitHub Enterprise의 경우 추가 권한이 필요할 수 있습니다.
 * 
 * 환경 변수 설정 방법:
 * .env 파일에 다음과 같이 설정합니다:
 * 
 * ```
 * # GitHub 개인 액세스 토큰
 * GITHUB_TOKEN=ghp_your_personal_access_token
 * 
 * # 보호할 브랜치 목록 (쉼표로 구분)
 * PROTECTED_BRANCHES=main,develop,staging,release
 * ```
 */

import * as dotenv from "dotenv";
import { 
  getGitRemoteInfo, 
  createGitHubClient, 
  getDefaultBranch, 
  getAllBranches, 
  createBranch 
} from "../utils/github";
import { isConfigured, recordConfigured } from "../utils/gh-settings";

dotenv.config();

const featureName = "protect-branch";

export async function protectBranches() {
  if (isConfigured(featureName)) {
    return;
  }

  const token = process.env.GITHUB_TOKEN!;
  const branchesEnv = process.env.PROTECTED_BRANCHES;

  if (!token) {
    console.error("❌ GITHUB_TOKEN이 필요합니다.");
    process.exit(1);
  }

  if (!branchesEnv) {
    console.error("❌ PROTECTED_BRANCHES 환경 변수가 필요합니다.");
    process.exit(1);
  }

  // 쉼표로 구분된 브랜치 목록을 배열로 변환
  const branches = branchesEnv.split(",").map(branch => branch.trim());

  const { owner, repo } = getGitRemoteInfo();
  const octokit = createGitHubClient(token);

  try {
    const defaultBranch = await getDefaultBranch(octokit, owner, repo);
    console.log(`ℹ️ 기본 브랜치: ${defaultBranch}`);
    
    const existingBranches = await getAllBranches(octokit, owner, repo);
    console.log(`ℹ️ 현재 브랜치 수: ${existingBranches.length}`);
    
    // 각 브랜치에 대해 작업 수행
    for (const branch of branches) {
      console.log(`\n--------------------------------`);
      console.log(`ℹ️ 브랜치 처리 중: ${branch}`);

      // 브랜치가 이미 존재하는지 확인
      if (existingBranches.includes(branch)) {
        console.log(`ℹ️ 브랜치가 이미 존재합니다: ${branch}`);
      } else {
        // 브랜치가 없으면 생성
        console.log(`🔍 브랜치가 존재하지 않습니다. 생성 중: ${branch}`);
        try {
          await createBranch(octokit, owner, repo, branch, defaultBranch);
        } catch (error) {
          console.error(`⚠️ ${branch} 브랜치 생성 오류. 생성을 건너뜁니다.`);
          continue; // 현재 브랜치의 처리를 건너뛰고 다음 브랜치로 이동
        }
      }
      
      // 브랜치 보호 규칙 설정
      console.log(`🔒 브랜치 보호 규칙 설정 중: ${owner}/${repo}#${branch}`);
      
      await octokit.repos.updateBranchProtection({
        owner,
        repo,
        branch,
        required_status_checks: null,
        enforce_admins: false,  // 관리자(owner)가 제한 없이 변경할 수 있도록 false로 설정
        required_pull_request_reviews: {
          required_approving_review_count: 1,
        },
        restrictions: null,
      });
      
      console.log(`✅ 브랜치 보호 규칙 설정 완료: ${owner}/${repo}#${branch}`);
    }
    
    console.log(`\n🎉 모든 브랜치(${branches.join(", ")})에 대한 보호 규칙이 설정되었습니다.`);
    
    recordConfigured(featureName);
  } catch (error) {
    console.error(`❌ 브랜치 보호 규칙 설정 중 오류 발생:`, error);
    process.exit(1);
  }
}

// 직접 실행될 때만 실행
if (require.main === module) {
  protectBranches();
} 