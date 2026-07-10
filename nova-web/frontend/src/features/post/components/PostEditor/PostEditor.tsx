import { useEffect, useMemo, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { FilesApi } from "@/api/generated/apis/FilesApi";
import { PostsApi } from "@/api/generated/apis/PostsApi";
import {
  FileScene,
  PostStatus,
  UpdatePostRequestUpdateFieldsEnum
} from "@/api/generated/models";
import { ResponseError } from "@/api/generated/runtime";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { csrfHeaders } from "@/lib/csrf";
import { cn } from "@/lib/utils";

import { responseErrorMessage } from "../../lib/postErrors";
import { redirectAfterSave } from "../../lib/postUrls";

type PostEditorProps = DOMStringMap;
type SaveStatus = "idle" | "loading" | "saving" | "success" | "error";

const fieldControlClassName =
  "block w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground shadow-xs transition-colors outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:bg-muted disabled:text-muted-foreground disabled:opacity-70";

function PostEditor(props: PostEditorProps) {
  const mode = props.mode === "edit" ? "edit" : "create";
  const postID = props.postId || "";
  const username = props.username || "";
  const [title, setTitle] = useState("");
  const [summary, setSummary] = useState("");
  const [content, setContent] = useState("");
  const [status, setStatus] = useState<PostStatus>(PostStatus.POST_STATUS_DRAFT);
  const [coverURL, setCoverURL] = useState("");
  const [coverFile, setCoverFile] = useState<File | null>(null);
  const [coverPreviewURL, setCoverPreviewURL] = useState("");
  const [saveStatus, setSaveStatus] = useState<SaveStatus>(mode === "edit" ? "loading" : "idle");
  const [message, setMessage] = useState("");
  const configuration = useMemo(() => apiConfiguration(), []);
  const postsAPI = useMemo(() => new PostsApi(configuration), [configuration]);
  const filesAPI = useMemo(() => new FilesApi(configuration), [configuration]);
  const formDisabled = saveStatus === "loading";
  const submitDisabled = formDisabled || saveStatus === "saving";

  useEffect(() => {
    if (mode !== "edit") {
      return;
    }
    if (!postID || !username) {
      setSaveStatus("error");
      setMessage(postID ? "缺少用户名" : "缺少文章 ID");
      return;
    }

    let canceled = false;
    setSaveStatus("loading");
    setMessage("");
    postsAPI
      .getUserPost({ username, id: postID })
      .then((payload) => {
        if (canceled) {
          return;
        }
        const post = payload.post;
        setTitle(post.title || "");
        setSummary(post.summary || "");
        setContent(post.content || "");
        setStatus(post.status || PostStatus.POST_STATUS_DRAFT);
        setCoverURL(post.coverUrl || "");
        setSaveStatus("idle");
      })
      .catch(async (error: unknown) => {
        if (canceled) {
          return;
        }
        if (error instanceof ResponseError && error.response.status === 401) {
          window.location.href = "/login";
          return;
        }
        setSaveStatus("error");
        setMessage(await responseErrorMessage(error, "加载文章失败"));
      });

    return () => {
      canceled = true;
    };
  }, [mode, postID, postsAPI, username]);

  useEffect(() => {
    return () => {
      if (coverPreviewURL) {
        URL.revokeObjectURL(coverPreviewURL);
      }
    };
  }, [coverPreviewURL]);

  async function savePost() {
    if (mode === "edit" && !postID) {
      setSaveStatus("error");
      setMessage("缺少文章 ID");
      return;
    }

    setSaveStatus("saving");
    setMessage("");

    try {
      const nextCoverURL = await uploadCoverIfNeeded(filesAPI, coverFile, coverURL);
      if (mode === "edit") {
        const payload = await postsAPI.updatePost(
          {
            id: postID,
            updatePostRequest: {
              post: {
                title,
                summary,
                content,
                status,
                coverUrl: nextCoverURL
              },
              updateFields: [
                UpdatePostRequestUpdateFieldsEnum.Title,
                UpdatePostRequestUpdateFieldsEnum.Summary,
                UpdatePostRequestUpdateFieldsEnum.Content,
                UpdatePostRequestUpdateFieldsEnum.Status,
                UpdatePostRequestUpdateFieldsEnum.CoverUrl
              ]
            }
          },
          { headers: csrfHeaders() }
        );
        redirectAfterSave(payload.post?.id || postID, payload.post?.status || status);
        return;
      }

      const payload = await postsAPI.createPost(
        {
          createPostRequest: {
            post: {
              title,
              summary,
              content,
              status,
              coverUrl: nextCoverURL
            }
          }
        },
        { headers: csrfHeaders() }
      );
      redirectAfterSave(payload.post?.id || "", payload.post?.status || status);
    } catch (error) {
      if (error instanceof ResponseError && error.response.status === 401) {
        window.location.href = "/login";
        return;
      }
      setSaveStatus("error");
      setMessage(await responseErrorMessage(error, mode === "edit" ? "保存文章失败" : "创建文章失败"));
    }
  }

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        void savePost();
      }}
      className="space-y-5"
    >
      {message ? (
        <Alert variant={saveStatus === "error" ? "destructive" : "success"}>
          <AlertDescription>{message}</AlertDescription>
        </Alert>
      ) : null}

      <div>
        <Label className="mb-1.5" htmlFor="title">
          标题
        </Label>
        <Input
          id="title"
          name="title"
          value={title}
          onChange={(event) => setTitle(event.currentTarget.value)}
          disabled={formDisabled}
          required
        />
      </div>

      <div>
        <Label className="mb-1.5" htmlFor="summary">
          摘要
        </Label>
        <textarea
          id="summary"
          name="summary"
          className={cn(fieldControlClassName, "min-h-24")}
          value={summary}
          onChange={(event) => setSummary(event.currentTarget.value)}
          disabled={formDisabled}
        />
      </div>

      <div>
        <Label className="mb-1.5" htmlFor="cover">
          封面图
        </Label>
        <Input
          id="cover"
          name="cover"
          type="file"
          accept="image/png,image/jpeg,image/webp"
          disabled={formDisabled}
          onChange={(event) => {
            const file = event.currentTarget.files?.[0] || null;
            setCoverFile(file);
            setCoverPreviewURL((current) => {
              if (current) {
                URL.revokeObjectURL(current);
              }
              return file ? URL.createObjectURL(file) : "";
            });
          }}
        />
        {coverPreviewURL || coverURL ? (
          <img className="mt-3 max-h-60 w-full rounded-md object-cover" src={coverPreviewURL || coverURL} alt="" />
        ) : null}
      </div>

      <div>
        <Label className="mb-1.5" htmlFor="status">
          状态
        </Label>
        <select
          id="status"
          name="status"
          className={fieldControlClassName}
          value={status}
          onChange={(event) => setStatus(event.currentTarget.value as PostStatus)}
          disabled={formDisabled}
        >
          <option value={PostStatus.POST_STATUS_DRAFT}>草稿</option>
          <option value={PostStatus.POST_STATUS_PUBLISHED}>发布</option>
          {status === PostStatus.POST_STATUS_PENDING_REVIEW ? (
            <option value={PostStatus.POST_STATUS_PENDING_REVIEW}>审核中</option>
          ) : null}
        </select>
      </div>

      <div>
        <Label className="mb-1.5" htmlFor="content">
          正文
        </Label>
        <textarea
          id="content"
          name="content"
          className={cn(fieldControlClassName, "min-h-64")}
          rows={10}
          value={content}
          onChange={(event) => setContent(event.currentTarget.value)}
          disabled={formDisabled}
          required
        />
      </div>

      <div className="flex flex-col gap-2 sm:flex-row">
        <Button type="submit" disabled={submitDisabled}>
          {saveStatus === "saving" ? "保存中..." : mode === "edit" ? "保存" : "创建"}
        </Button>
        <Button type="button" variant="outline" onClick={() => (window.location.href = "/posts")}>
          取消
        </Button>
      </div>
    </form>
  );
}

async function uploadCoverIfNeeded(filesAPI: FilesApi, file: File | null, currentCoverURL: string) {
  if (!file) {
    return currentCoverURL;
  }
  const payload = await filesAPI.uploadFile(
    {
      scene: FileScene.FILE_SCENE_POST_COVER,
      file
    },
    { headers: csrfHeaders() }
  );
  return String(payload.file?.objectKey || payload.file?.url || "");
}

export { PostEditor };
