import { useState } from "react";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AuthApi } from "@/api/generated/apis/AuthApi";
import { ResponseError } from "@/api/generated/runtime";
import { apiConfiguration } from "@/api/config";
import { csrfHeaders } from "@/lib/csrf";

function PasswordSettings() {
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newPasswordConfirmation, setNewPasswordConfirmation] = useState("");
  const [status, setStatus] = useState<"idle" | "saving" | "success" | "error">("idle");
  const [message, setMessage] = useState("");
  const api = new AuthApi(apiConfiguration());

  async function savePassword() {
    const validationError = validatePassword(newPassword, newPasswordConfirmation);
    if (validationError) {
      setStatus("error");
      setMessage(validationError);
      return;
    }

    setStatus("saving");
    setMessage("");

    try {
      await api.changePassword({
        changePasswordRequest: {
          oldPassword,
          newPassword,
          newPasswordConfirmation
        }
      }, {
        headers: csrfHeaders()
      });
    } catch (error) {
      if (error instanceof ResponseError && error.response.status === 401) {
        window.location.href = "/login";
        return;
      }
      setStatus("error");
      setMessage(await responseErrorMessage(error, "更新密码失败"));
      return;
    }

    setOldPassword("");
    setNewPassword("");
    setNewPasswordConfirmation("");
    setStatus("success");
    setMessage("密码已更新");
  }

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        void savePassword();
      }}
    >
      {message ? (
        <Alert className="mb-4" variant={status === "error" ? "destructive" : "success"}>
          <AlertDescription>{message}</AlertDescription>
        </Alert>
      ) : null}
      <div className="mb-4">
        <Label className="mb-1.5" htmlFor="old_password">
          当前密码
        </Label>
        <Input
          type="password"
          id="old_password"
          name="old_password"
          value={oldPassword}
          onChange={(event) => setOldPassword(event.currentTarget.value)}
          required
          autoComplete="current-password"
        />
      </div>
      <div className="mb-4">
        <Label className="mb-1.5" htmlFor="new_password">
          新密码
        </Label>
        <Input
          type="password"
          id="new_password"
          name="new_password"
          value={newPassword}
          onChange={(event) => setNewPassword(event.currentTarget.value)}
          required
          minLength={8}
          maxLength={64}
          autoComplete="new-password"
        />
        <div className="mt-1 text-sm text-muted-foreground">新密码需要 8-64 个字符，且必须同时包含字母和数字。</div>
      </div>
      <div className="mb-4">
        <Label className="mb-1.5" htmlFor="new_password_confirmation">
          确认新密码
        </Label>
        <Input
          type="password"
          id="new_password_confirmation"
          name="new_password_confirmation"
          value={newPasswordConfirmation}
          onChange={(event) => setNewPasswordConfirmation(event.currentTarget.value)}
          required
          autoComplete="new-password"
        />
      </div>
      <div className="flex items-center gap-3">
        <Button type="submit" variant="outline" disabled={status === "saving"}>
          {status === "saving" ? "更新中..." : "更新密码"}
        </Button>
      </div>
    </form>
  );
}

function validatePassword(password: string, confirmation: string) {
  if (!/[a-zA-Z]/.test(password) || !/[0-9]/.test(password)) {
    return "密码必须同时包含字母和数字";
  }
  if (password !== confirmation) {
    return "两次输入的密码不一致";
  }
  return "";
}

async function responseErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof ResponseError)) {
    return "网络错误，请稍后重试";
  }
  const payload = await error.response.json().catch(() => ({}));
  return payload?.error?.message || fallback;
}

export { PasswordSettings };
