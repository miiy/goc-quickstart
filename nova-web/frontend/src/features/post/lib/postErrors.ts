import { ResponseError } from "@/api/generated/runtime";

async function responseErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof ResponseError)) {
    return "网络错误，请稍后重试";
  }
  const payload = await error.response.json().catch(() => ({}));
  return payload?.error?.message || fallback;
}

export { responseErrorMessage };
