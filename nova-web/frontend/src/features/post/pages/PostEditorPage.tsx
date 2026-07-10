import { useEffect, useMemo, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { UsersApi } from "@/api/generated/apis/UsersApi";
import { ResponseError } from "@/api/generated/runtime";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { PostEditor } from "../components/PostEditor";
import { responseErrorMessage } from "../lib/postErrors";
import { editPostIDFromPath } from "../lib/postUrls";

type PostEditorPageProps = DOMStringMap;

function PostEditorPage(props: PostEditorPageProps) {
  const mode = props.mode === "edit" ? "edit" : "create";
  const postID = props.postId || (mode === "edit" ? editPostIDFromPath() : "");
  const [username, setUsername] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "error">(mode === "edit" ? "loading" : "idle");
  const [message, setMessage] = useState("");
  const usersAPI = useMemo(() => new UsersApi(apiConfiguration()), []);

  useEffect(() => {
    if (mode !== "edit") {
      return;
    }
    let canceled = false;
    setStatus("loading");
    setMessage("");

    usersAPI
      .getProfile()
      .then((payload) => {
        if (canceled) {
          return;
        }
        setUsername(payload.user.username || "");
        setStatus("idle");
      })
      .catch(async (error: unknown) => {
        if (canceled) {
          return;
        }
        if (error instanceof ResponseError && error.response.status === 401) {
          window.location.href = "/login";
          return;
        }
        setStatus("error");
        setMessage(await responseErrorMessage(error, "加载用户信息失败"));
      });

    return () => {
      canceled = true;
    };
  }, [mode, usersAPI]);

  return (
    <Card>
      <CardHeader>
        <CardTitle>{mode === "edit" ? "编辑文章" : "创建文章"}</CardTitle>
      </CardHeader>
      <CardContent>
        {status === "loading" ? <div className="text-slate-500">加载中...</div> : null}
        {status === "error" ? (
          <Alert variant="destructive">
            <AlertDescription>{message}</AlertDescription>
          </Alert>
        ) : null}
        {status === "idle" ? <PostEditor mode={mode} postId={postID} username={username} /> : null}
      </CardContent>
    </Card>
  );
}

export { PostEditorPage };
