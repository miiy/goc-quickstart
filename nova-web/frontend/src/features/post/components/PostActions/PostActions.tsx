import { useMemo, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { PostsApi } from "@/api/generated/apis/PostsApi";
import { ResponseError } from "@/api/generated/runtime";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { csrfHeaders } from "@/lib/csrf";

import { responseErrorMessage } from "../../lib/postErrors";

type PostActionsProps = DOMStringMap;

function PostActions(props: PostActionsProps) {
  const postID = props.postId || "";
  const [status, setStatus] = useState<"idle" | "deleting" | "error">("idle");
  const [message, setMessage] = useState("");
  const api = useMemo(() => new PostsApi(apiConfiguration()), []);

  async function deletePost() {
    if (!postID || !window.confirm("确定删除这篇文章？")) {
      return;
    }

    setStatus("deleting");
    setMessage("");

    try {
      await api.deletePost(
        { id: postID },
        {
          headers: csrfHeaders()
        }
      );
      window.location.href = "/posts";
    } catch (error) {
      if (error instanceof ResponseError && error.response.status === 401) {
        window.location.href = "/login";
        return;
      }
      setStatus("error");
      setMessage(await responseErrorMessage(error, "删除文章失败"));
    }
  }

  return (
    <div className="mt-4 space-y-3">
      {message ? (
        <Alert variant="destructive">
          <AlertDescription>{message}</AlertDescription>
        </Alert>
      ) : null}
      <div className="flex flex-col gap-2 sm:flex-row">
        <Button asChild>
          <a href={`/posts/${postID}/edit`}>编辑</a>
        </Button>
        <Button type="button" variant="destructive" onClick={() => void deletePost()} disabled={status === "deleting"}>
          {status === "deleting" ? "删除中..." : "删除"}
        </Button>
      </div>
    </div>
  );
}

export { PostActions };
