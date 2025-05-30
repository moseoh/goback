/**
 * GitHub 저장소 관리 도구 - 메인 진입점
 * 
 * 이 파일은 모든 기능을 하나로 모아서 제공합니다.
 * 특정 기능만 사용하려면 각 액션 파일을 직접 실행할 수도 있습니다.
 */

import * as dotenv from "dotenv";
import { protectBranches } from "./actions/protect-branch";
import { 
  enableAutoDeleteMergedBranches
} from "./actions/auto-delete-branch";
import {
  setSquashMergePreference
} from "./actions/set-squash-merge";

dotenv.config();

// 사용 가능한 명령어 목록
const COMMANDS = {
  "protect": "브랜치 보호 규칙 설정",
  "auto-delete": "머지된 PR의 브랜치 자동 삭제 옵션 활성화",
  "check-auto-delete": "머지된 PR의 브랜치 자동 삭제 옵션 상태 확인",
  "squash-merge": "PR 병합 방식을 Squash merge로 설정",
  "check-merge": "현재 PR 병합 방식 설정 확인",
  "all": "모든 기능 실행 (브랜치 보호 규칙 설정 + 자동 삭제 옵션 활성화 + Squash merge 설정)"
};

// 도움말 표시
function showHelp() {
  console.log("\n🛠️ GitHub 저장소 관리 도구 🛠️");
  console.log("\n사용 방법: npm run start [명령어]");
  console.log("\n사용 가능한 명령어:");
  
  Object.entries(COMMANDS).forEach(([cmd, desc]) => {
    console.log(`  - ${cmd}: ${desc}`);
  });
  
  console.log("\n예시:");
  console.log("  npm run start protect           # 브랜치 보호 규칙 설정");
  console.log("  npm run start auto-delete       # 머지된 PR의 브랜치 자동 삭제 옵션 활성화");
  console.log("  npm run start squash-merge      # PR 병합 방식을 Squash merge로 설정");
  console.log("  npm run start all               # 모든 기능 실행");
  console.log("");
}

// 메인 함수
async function main() {
  const command = process.argv[2] || "help";

  switch (command) {
    case "protect":
      await protectBranches();
      break;
    case "auto-delete":
      await enableAutoDeleteMergedBranches();
      break;
    case "squash-merge":
      await setSquashMergePreference();
      break;
    case "all":
      await protectBranches();
      await enableAutoDeleteMergedBranches();
      await setSquashMergePreference();
      break;
    case "help":
    default:
      showHelp();
      break;
  }
}

// 직접 실행될 때만 실행
if (require.main === module) {
  main().catch(error => {
    console.error("❌ 오류 발생:", error);
    process.exit(1);
  });
} 