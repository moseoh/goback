// .repo/githook/postinstall.mjs
import { execSync } from "node:child_process";
import { mkdirSync, writeFileSync } from "node:fs";
import { join } from "node:path";

// Git 사용자 이름 가져오기
let username = "unknown";
try {
  username = execSync("git config user.name", { encoding: "utf-8" }).trim();
} catch {
  throw new Error("git user.name을 찾을 수 없습니다");
}

const folder = ".hook-installed";

// 사용자 이름 파일 생성
const path = join(folder, username);
writeFileSync(path, `installed\n`);

console.log(`✅ githook 설치 기록: ${path}`);
