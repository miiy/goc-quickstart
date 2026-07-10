import { PostStatus } from "@/api/generated/models";

function editPostIDFromPath() {
  if (typeof window === "undefined") {
    return "";
  }
  const match = window.location.pathname.match(/\/posts\/([^/]+)\/edit$/);
  return match?.[1] ? decodeURIComponent(match[1]) : "";
}

function redirectAfterSave(id: string, status: PostStatus) {
  if (!id) {
    window.location.href = "/posts";
    return;
  }
  window.location.href = status === PostStatus.POST_STATUS_PUBLISHED ? `/posts/${id}` : `/posts/${id}/edit`;
}

export { editPostIDFromPath, redirectAfterSave };
