import { useEffect, useMemo, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { UsersApi } from "@/api/generated/apis/UsersApi";
import type { User } from "@/api/generated/models";
import { ResponseError } from "@/api/generated/runtime";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { PasswordSettings } from "../components/PasswordSettings";
import { ProfileSettings } from "../components/ProfileSettings";

type ProfilePageProps = DOMStringMap;

function ProfilePage(props: ProfilePageProps) {
  const [user, setUser] = useState<User | null>(null);
  const [status, setStatus] = useState<"loading" | "ready" | "error">("loading");
  const [message, setMessage] = useState("");
  const api = useMemo(() => new UsersApi(apiConfiguration()), []);

  useEffect(() => {
    let canceled = false;

    api
      .getProfile()
      .then((payload) => {
        if (canceled) {
          return;
        }
        setUser(payload.user);
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
        setMessage(await responseErrorMessage(error, "加载资料失败"));
      });

    return () => {
      canceled = true;
    };
  }, [api]);

  if (status === "loading") {
    return (
      <div className="mx-auto max-w-5xl">
        <Card>
          <CardContent className="text-slate-500">加载中...</CardContent>
        </Card>
      </div>
    );
  }

  if (status === "error" || user === null) {
    return (
      <div className="mx-auto max-w-5xl">
        <Alert variant="destructive">
          <AlertDescription>{message || "加载资料失败"}</AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="mx-auto grid max-w-5xl gap-6 md:grid-cols-[minmax(0,320px)_1fr]">
      <ProfileSummary user={user} />
      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>资料设置</CardTitle>
          </CardHeader>
          <CardContent>
            <ProfileSettings user={user} onUserChange={setUser} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>修改密码</CardTitle>
          </CardHeader>
          <CardContent>
            <PasswordSettings />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function ProfileSummary({ user }: { user: User }) {
  const nickname = user.nickname?.trim() || "";

  return (
    <Card>
      <CardContent className="text-center">
        {user.avatar ? (
          <img src={user.avatar} alt={user.username} className="mx-auto mb-4 h-24 w-24 rounded-full object-cover" />
        ) : (
          <div className="mx-auto mb-4 flex h-24 w-24 items-center justify-center rounded-full bg-slate-100 text-3xl font-semibold text-slate-500">
            U
          </div>
        )}
        <h1 className="text-lg font-semibold text-slate-950">{user.username || "用户中心"}</h1>
        {nickname ? <div className="mt-1 text-slate-500">{nickname}</div> : null}
      </CardContent>
      <ul className="divide-y divide-slate-200">
        <ProfileSummaryItem label="username" value={user.username} />
        <ProfileSummaryItem label="nickname" value={nickname || "未设置"} />
        <ProfileSummaryItem label="email" value={user.email || "未设置"} />
      </ul>
    </Card>
  );
}

function ProfileSummaryItem({ label, value }: { label: string; value: string }) {
  return (
    <li className="px-5 py-4">
      <div className="text-sm text-slate-500">{label}</div>
      <div className="mt-1 text-slate-900">{value}</div>
    </li>
  );
}

async function responseErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof ResponseError)) {
    return "网络错误，请稍后重试";
  }
  const payload = await error.response.json().catch(() => ({}));
  return payload?.error?.message || fallback;
}

export { ProfilePage };
