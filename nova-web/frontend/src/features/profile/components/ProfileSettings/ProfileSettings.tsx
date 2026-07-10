import { useEffect, useMemo, useRef, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { UsersApi } from "@/api/generated/apis/UsersApi";
import type { User } from "@/api/generated/models";
import { ResponseError } from "@/api/generated/runtime";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { csrfHeaders } from "@/lib/csrf";

import { buildProfileUpdate, type ProfileFields } from "../../lib/profileUpdate";

type ProfileSettingsProps = {
  user: User;
  onUserChange: (user: User) => void;
};

function ProfileSettings(props: ProfileSettingsProps) {
  const initialProfile = useMemo(() => userProfileFields(props.user), [props.user]);
  const [profile, setProfile] = useState<ProfileFields>(initialProfile);
  const [savedProfile, setSavedProfile] = useState<ProfileFields>(initialProfile);
  const [status, setStatus] = useState<"idle" | "saving" | "success" | "error">("idle");
  const [message, setMessage] = useState("");
  const savingRef = useRef(false);
  const api = useMemo(() => new UsersApi(apiConfiguration()), []);
  const profileUpdate = buildProfileUpdate(savedProfile, profile);
  const isSaving = status === "saving";

  useEffect(() => {
    setProfile(initialProfile);
    setSavedProfile(initialProfile);
    setStatus("idle");
    setMessage("");
  }, [initialProfile]);

  async function saveProfile() {
    const update = buildProfileUpdate(savedProfile, profile);
    if (!update || savingRef.current) {
      return;
    }

    savingRef.current = true;
    setStatus("saving");
    setMessage("");

    try {
      const payload = await api.updateProfile({
        updateProfileRequest: {
          user: update.user,
          updateFields: update.updateFields
        }
      }, {
        headers: csrfHeaders()
      });
      const nextProfile: ProfileFields = {
        nickname: payload.user.nickname ?? "",
        email: payload.user.email ?? ""
      };
      setProfile(nextProfile);
      setSavedProfile(nextProfile);
      props.onUserChange(payload.user);
      setStatus("success");
      setMessage("资料已更新");
    } catch (error) {
      if (error instanceof ResponseError && error.response.status === 401) {
        window.location.href = "/login";
        return;
      }
      setStatus("error");
      setMessage(await responseErrorMessage(error, "保存资料失败"));
    } finally {
      savingRef.current = false;
    }
  }

  function updateField(field: keyof ProfileFields, value: string) {
    setProfile((current) => ({ ...current, [field]: value }));
    setStatus("idle");
    setMessage("");
  }

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        void saveProfile();
      }}
    >
      {message ? (
        <Alert className="mb-4" variant={status === "error" ? "destructive" : "success"}>
          <AlertDescription>{message}</AlertDescription>
        </Alert>
      ) : null}
      <fieldset className="m-0 min-w-0 border-0 p-0" disabled={isSaving}>
        <div className="mb-4">
          <Label className="mb-1.5" htmlFor="username">
            username
          </Label>
          <Input type="text" id="username" value={props.user.username ?? ""} disabled />
        </div>
        <div className="mb-4">
          <Label className="mb-1.5" htmlFor="nickname">
            nickname
          </Label>
          <Input
            type="text"
            id="nickname"
            name="nickname"
            value={profile.nickname}
            onChange={(event) => updateField("nickname", event.currentTarget.value)}
          />
        </div>
        <div className="mb-4">
          <Label className="mb-1.5" htmlFor="email">
            email
          </Label>
          <Input
            type="email"
            id="email"
            name="email"
            value={profile.email}
            onChange={(event) => updateField("email", event.currentTarget.value)}
          />
        </div>
        <Button type="submit" disabled={isSaving || !profileUpdate}>
          {isSaving ? "保存中..." : "保存资料"}
        </Button>
      </fieldset>
    </form>
  );
}

function userProfileFields(user: User): ProfileFields {
  return {
    nickname: user.nickname ?? "",
    email: user.email ?? ""
  };
}

async function responseErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof ResponseError)) {
    return "网络错误，请稍后重试";
  }
  const payload = await error.response.json().catch(() => ({}));
  return payload?.error?.message || fallback;
}

export { ProfileSettings };
