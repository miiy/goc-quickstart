import { useEffect, useMemo, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { PostsApi } from "@/api/generated/apis/PostsApi";
import { UsersApi } from "@/api/generated/apis/UsersApi";
import type { Post, PublicUser, User } from "@/api/generated/models";
import { ResponseError } from "@/api/generated/runtime";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious
} from "@/components/ui/pagination";

import "./UserShowPage.css";

type UserShowPageProps = DOMStringMap;
type LoadStatus = "loading" | "ready" | "error";
type UserSection = "posts" | "activity" | "followers" | "following";

const pageSize = 10;

function UserShowPage(props: UserShowPageProps) {
  const username = props.username || "";
  const [section, setSection] = useState<UserSection>(() => normalizeSection(props.section || sectionFromLocation()));
  const [user, setUser] = useState<PublicUser | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [total, setTotal] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [status, setStatus] = useState<LoadStatus>("loading");
  const [message, setMessage] = useState("");
  const configuration = useMemo(() => apiConfiguration(), []);
  const usersAPI = useMemo(() => new UsersApi(configuration), [configuration]);
  const postsAPI = useMemo(() => new PostsApi(configuration), [configuration]);

  useEffect(() => {
    function syncSection() {
      setSection(sectionFromLocation());
    }
    window.addEventListener("popstate", syncSection);
    return () => window.removeEventListener("popstate", syncSection);
  }, []);

  useEffect(() => {
    setCurrentPage(1);
  }, [section, username]);

  useEffect(() => {
    let canceled = false;
    setStatus("loading");
    setMessage("");

    Promise.all([usersAPI.getUser({ username }), getOptionalProfile(usersAPI)])
      .then(async ([userPayload, nextCurrentUser]) => {
        if (canceled) {
          return;
        }
        const nextUser = userPayload.user;
        const nextIsOwnProfile = isProfileOwner(nextCurrentUser, username);
        setUser(nextUser);

        if ((section === "posts" || section === "activity") && nextUser?.id) {
          const postsPayload = nextIsOwnProfile
            ? await postsAPI.listUserPosts({ username, page: currentPage, pageSize })
            : await postsAPI.listPosts({ userId: nextUser.id, page: currentPage, pageSize });
          if (canceled) {
            return;
          }
          setPosts(postsPayload.posts || []);
          setTotal(postsPayload.total || 0);
          setTotalPages(postsPayload.totalPages || 0);
          setCurrentPage(postsPayload.currentPage || currentPage);
        } else {
          setPosts([]);
          setTotal(0);
          setTotalPages(0);
        }
        setStatus("ready");
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
  }, [currentPage, postsAPI, section, username, usersAPI]);

  function selectSection(nextSection: UserSection) {
    setSection(nextSection);
    setCurrentPage(1);
    updateSectionURL(username, nextSection);
  }

  if (status === "loading") {
    return (
      <Card>
        <CardContent className="text-slate-500">加载中...</CardContent>
      </Card>
    );
  }

  if (status === "error") {
    return (
      <Alert variant="destructive">
        <AlertDescription>{message}</AlertDescription>
      </Alert>
    );
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[minmax(0,320px)_1fr]">
      <aside>
        <Card>
          <CardContent className="text-center">
            {user?.avatar ? (
              <img src={user.avatar} alt={user.username} className="mx-auto mb-4 h-28 w-28 rounded-full object-cover" />
            ) : (
              <div className="mx-auto mb-4 flex h-28 w-28 items-center justify-center rounded-full bg-slate-100 text-4xl font-semibold text-slate-500">
                U
              </div>
            )}
            <h1 className="text-xl font-semibold text-slate-950">{userDisplayName(user, username)}</h1>
            <div className="mt-1 text-slate-500">@{user?.username || username}</div>
            {user?.createdAt ? <div className="mt-2 text-sm text-slate-500">加入于 {formatDate(user.createdAt)}</div> : null}
          </CardContent>
          <div className="divide-y divide-slate-200">
            {userTabs(user?.username || username, section).map((tab) => (
              <a
                key={tab.section}
                href={tab.href}
                className={`user-show-tab-link ${tab.active ? "user-show-tab-link-active" : ""}`}
                onClick={(event) => {
                  event.preventDefault();
                  selectSection(tab.section);
                }}
              >
                {tab.label}
              </a>
            ))}
          </div>
        </Card>
      </aside>

      <main>
        {section === "followers" || section === "following" ? (
          <Card>
            <CardContent className="text-slate-500">{section === "followers" ? "暂无粉丝。" : "暂无关注。"}</CardContent>
          </Card>
        ) : (
          <Card>
            <CardHeader className="flex flex-col gap-1 bg-slate-50 sm:flex-row sm:items-center sm:justify-between">
              <h2 className="text-xl font-semibold text-slate-950">{section === "activity" ? "动态" : "文章"}</h2>
              <span className="text-sm text-slate-500">{total} 篇</span>
            </CardHeader>
            <ul className="divide-y divide-slate-200">
              {posts.length > 0 ? (
                posts.map((post) => (
                  <li key={post.id} className="px-5 py-4">
                    <div className="flex min-w-0 items-center gap-4">
                      {post.coverUrl ? (
                        <img src={post.coverUrl} alt="" className="h-16 w-24 shrink-0 rounded-md object-cover sm:h-14 sm:w-[88px]" />
                      ) : null}
                      <div className="min-w-0">
                        <a href={postHref(post)} className="font-medium text-slate-950 no-underline hover:text-brand-700 hover:underline">
                          {post.title}
                        </a>
                        <div className="mt-1 flex flex-wrap items-center gap-2 text-sm text-slate-500">
                          <span>{formatDate(post.createdAt)}</span>
                          {post.canManage && post.status !== "published" ? <PostStatusBadge status={post.status} /> : null}
                        </div>
                      </div>
                    </div>
                  </li>
                ))
              ) : (
                <li className="px-5 py-4 text-slate-500">暂无内容。</li>
              )}
            </ul>
            {totalPages > 1 ? (
              <div className="border-t border-border px-5 py-4">
                <UserPostsPagination currentPage={currentPage} totalPages={totalPages} onPageChange={setCurrentPage} />
              </div>
            ) : null}
          </Card>
        )}
      </main>
    </div>
  );
}

function UserPostsPagination({
  currentPage,
  totalPages,
  onPageChange
}: {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}) {
  const pages = paginationItems(currentPage, totalPages);

  function goToPage(page: number) {
    if (page < 1 || page > totalPages || page === currentPage) {
      return;
    }
    onPageChange(page);
  }

  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            href="#"
            disabled={currentPage <= 1}
            onClick={(event) => {
              event.preventDefault();
              goToPage(currentPage - 1);
            }}
          />
        </PaginationItem>
        {pages.map((page, index) => (
          <PaginationItem key={`${page}-${index}`}>
            {page === "ellipsis" ? (
              <PaginationEllipsis />
            ) : (
              <PaginationLink
                href="#"
                isActive={page === currentPage}
                onClick={(event) => {
                  event.preventDefault();
                  goToPage(page);
                }}
              >
                {page}
              </PaginationLink>
            )}
          </PaginationItem>
        ))}
        <PaginationItem>
          <PaginationNext
            href="#"
            disabled={currentPage >= totalPages}
            onClick={(event) => {
              event.preventDefault();
              goToPage(currentPage + 1);
            }}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}

function paginationItems(currentPage: number, totalPages: number): Array<number | "ellipsis"> {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, index) => index + 1);
  }

  if (currentPage <= 4) {
    return [1, 2, 3, 4, 5, "ellipsis", totalPages];
  }

  if (currentPage >= totalPages - 3) {
    return [1, "ellipsis", totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages];
  }

  return [1, "ellipsis", currentPage - 1, currentPage, currentPage + 1, "ellipsis", totalPages];
}

function postHref(post: Post) {
  if (post.canManage && post.status !== "published") {
    return `/posts/${post.id}/edit`;
  }
  return `/posts/${post.id}`;
}

async function getOptionalProfile(api: UsersApi): Promise<User | null> {
  try {
    const payload = await api.getProfile();
    return payload.user || null;
  } catch {
    return null;
  }
}

function isProfileOwner(currentUser: User | null, username: string) {
  return username !== "" && currentUser?.username === username;
}

function PostStatusBadge({ status }: { status?: string }) {
  return (
    <span className="rounded border border-border bg-slate-50 px-1.5 py-0.5 text-xs text-slate-600">
      {postStatusLabel(status)}
    </span>
  );
}

function postStatusLabel(status?: string) {
  switch (status) {
    case "draft":
      return "草稿";
    case "pending_review":
      return "审核中";
    default:
      return "未发布";
  }
}

function userDisplayName(user: PublicUser | null, fallback: string) {
  return user?.nickname || user?.username || fallback || "用户";
}

function userTabs(username: string, section: UserSection) {
  const tabs: Array<{ section: UserSection; label: string; href: string }> = [
    { section: "posts", label: "文章", href: `/users/${username}` },
    { section: "activity", label: "动态", href: `/users/${username}?section=activity` },
    { section: "followers", label: "粉丝", href: `/users/${username}?section=followers` },
    { section: "following", label: "关注", href: `/users/${username}?section=following` }
  ];
  return tabs.map((tab) => ({ ...tab, active: tab.section === section }));
}

function normalizeSection(value?: string | null): UserSection {
  switch (value) {
    case "activity":
    case "followers":
    case "following":
      return value;
    default:
      return "posts";
  }
}

function sectionFromLocation(): UserSection {
  if (typeof window === "undefined") {
    return "posts";
  }
  return normalizeSection(new URLSearchParams(window.location.search).get("section"));
}

function updateSectionURL(username: string, section: UserSection) {
  if (typeof window === "undefined" || username === "") {
    return;
  }
  const url = new URL(window.location.href);
  url.pathname = `/users/${username}`;
  if (section === "posts") {
    url.searchParams.delete("section");
  } else {
    url.searchParams.set("section", section);
  }
  window.history.pushState({}, "", `${url.pathname}${url.search}${url.hash}`);
}

function formatDate(value: string) {
  if (!value) {
    return "";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit"
  });
}

async function responseErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof ResponseError)) {
    return "网络错误，请稍后重试";
  }
  const payload = await error.response.json().catch(() => ({}));
  return payload?.error?.message || fallback;
}

export { UserShowPage };
