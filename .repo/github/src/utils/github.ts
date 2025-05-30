import { Octokit } from "@octokit/rest";
import { execSync } from "child_process";

// Git remote 정보에서 owner와 repo 가져오기
export function getGitRemoteInfo() {
  try {
    const remoteUrl = execSync("git config --get remote.origin.url")
      .toString()
      .trim();

    // HTTPS 형식: https://github.com/owner/repo.git
    // SSH 형식: git@github.com:owner/repo.git
    let owner = "";
    let repo = "";
    
    if (remoteUrl.includes("github.com")) {
      if (remoteUrl.startsWith("git@")) {
        // SSH 형식
        const match = remoteUrl.match(
          /git@github\.com:([^\/]+)\/([^\.]+)\.git/
        );
        if (match) {
          owner = match[1];
          repo = match[2];
        }
      } else {
        // HTTPS 형식
        const match = remoteUrl.match(
          /https:\/\/github\.com\/([^\/]+)\/([^\.]+)(?:\.git)?/
        );
        if (match) {
          owner = match[1];
          repo = match[2];
        }
      }
    }
    
    if (!owner || !repo) {
      throw new Error(
        "GitHub 원격 URL에서 소유자와 저장소 정보를 추출할 수 없습니다"
      );
    }
    
    return { owner, repo };
  } catch (error) {
    console.error("Git 원격 정보를 가져오는 중 오류가 발생했습니다:", error);
    process.exit(1);
  }
}

// GitHub API 클라이언트 생성
export function createGitHubClient(token: string): Octokit {
  if (!token) {
    throw new Error("GitHub 토큰이 필요합니다.");
  }
  return new Octokit({ auth: token });
}

// 저장소의 기본 브랜치 가져오기
export async function getDefaultBranch(octokit: Octokit, owner: string, repo: string): Promise<string> {
  try {
    const { data: repository } = await octokit.repos.get({
      owner,
      repo,
    });
    return repository.default_branch;
  } catch (error) {
    console.error("기본 브랜치 정보를 가져오는 중 오류가 발생했습니다:", error);
    throw error;
  }
}

// 저장소의 모든 브랜치 목록 가져오기
export async function getAllBranches(octokit: Octokit, owner: string, repo: string): Promise<string[]> {
  try {
    const branches: string[] = [];
    let page = 1;
    let hasNextPage = true;

    while (hasNextPage) {
      const response = await octokit.repos.listBranches({
        owner,
        repo,
        per_page: 100,
        page,
      });

      const branchNames = response.data.map(branch => branch.name);
      branches.push(...branchNames);
      
      hasNextPage = response.data.length === 100;
      page++;
    }

    return branches;
  } catch (error) {
    console.error("브랜치 목록을 가져오는 중 오류가 발생했습니다:", error);
    throw error;
  }
}

// 브랜치 생성하기
export async function createBranch(
  octokit: Octokit, 
  owner: string, 
  repo: string, 
  branchName: string, 
  sourceBranch: string
): Promise<void> {
  try {
    // 소스 브랜치의 최신 커밋 SHA 가져오기
    const { data: refData } = await octokit.git.getRef({
      owner,
      repo,
      ref: `heads/${sourceBranch}`,
    });

    const sha = refData.object.sha;

    // 새 브랜치 생성
    await octokit.git.createRef({
      owner,
      repo,
      ref: `refs/heads/${branchName}`,
      sha,
    });

    console.log(`✅ 새 브랜치 생성 완료: ${branchName} (소스: ${sourceBranch})`);
  } catch (error) {
    console.error(`❌ 브랜치 생성 중 오류 발생 (${branchName}):`, error);
    throw error;
  }
}

// 브랜치 삭제하기
export async function deleteBranch(
  octokit: Octokit,
  owner: string,
  repo: string,
  branchName: string
): Promise<void> {
  try {
    await octokit.git.deleteRef({
      owner,
      repo,
      ref: `heads/${branchName}`,
    });
    console.log(`✅ 브랜치 삭제 완료: ${branchName}`);
  } catch (error) {
    console.error(`❌ 브랜치 삭제 중 오류 발생 (${branchName}):`, error);
    throw error;
  }
} 