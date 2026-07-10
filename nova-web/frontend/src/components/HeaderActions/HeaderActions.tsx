import { useEffect, useId, useRef, useState } from "react";

import { apiConfiguration } from "@/api/config";
import { UsersApi } from "@/api/generated/apis/UsersApi";
import type { User } from "@/api/generated/models";
import { ResponseError } from "@/api/generated/runtime";
import { csrfToken } from "@/lib/csrf";

import "./HeaderActions.css";

type HeaderActionsProps = DOMStringMap;

type HeaderUser = {
  username: string;
};

let headerUserRequest: Promise<HeaderUser | null> | null = null;

function HeaderActions(props: HeaderActionsProps) {
  const variant = props.variant === "mobile" ? "mobile" : "desktop";
  const registerEnabled = props.registerEnabled === "true";
  const [user, setUser] = useState<HeaderUser | null>(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    let canceled = false;

    async function loadProfile() {
      try {
        const nextUser = await getHeaderUser();
        if (!canceled) {
          setUser(nextUser);
        }
      } catch (error) {
        if (!canceled && error instanceof ResponseError && error.response.status === 401) {
          setUser(null);
        }
      } finally {
        if (!canceled) {
          setLoaded(true);
        }
      }
    }

    void loadProfile();

    return () => {
      canceled = true;
    };
  }, []);

  if (!loaded) {
    return null;
  }

  if (variant === "mobile") {
    return <MobileHeaderActions registerEnabled={registerEnabled} user={user} />;
  }

  return <DesktopHeaderActions registerEnabled={registerEnabled} user={user} />;
}

function getHeaderUser() {
  headerUserRequest ??= new UsersApi(apiConfiguration())
    .getProfile()
    .then((payload) => headerUser(payload.user))
    .catch((error: unknown) => {
      if (error instanceof ResponseError && error.response.status === 401) {
        return null;
      }
      throw error;
    });
  return headerUserRequest;
}

function DesktopHeaderActions(props: {
  registerEnabled: boolean;
  user: HeaderUser | null;
}) {
  return (
    <ul className="flex items-center gap-1">
      {props.user ? (
        <li>
          <UserMenu user={props.user} />
        </li>
      ) : (
        <GuestLinks registerEnabled={props.registerEnabled} listItems />
      )}
    </ul>
  );
}

function MobileHeaderActions(props: {
  registerEnabled: boolean;
  user: HeaderUser | null;
}) {
  if (!props.user) {
    return <GuestLinks registerEnabled={props.registerEnabled} />;
  }

  const label = userLabel(props.user);
  return (
    <>
      <a className="site-nav-link" href={userHref(props.user)}>
        {label}
      </a>
      <a className="site-nav-link" href="/profile">
        设置
      </a>
      <LogoutForm buttonClassName="site-nav-action w-full" />
    </>
  );
}

function GuestLinks(props: { registerEnabled: boolean; listItems?: boolean }) {
  if (!props.listItems) {
    return (
      <>
        {props.registerEnabled ? (
          <a className="site-nav-link" href="/register">
            注册
          </a>
        ) : null}
        <a className="site-nav-link" href="/login">
          登录
        </a>
      </>
    );
  }

  return (
    <>
      {props.registerEnabled ? (
        <li>
          <a className="site-nav-link" href="/register">
            注册
          </a>
        </li>
      ) : null}
      <li>
        <a className="site-nav-link" href="/login">
          登录
        </a>
      </li>
    </>
  );
}

function UserMenu(props: { user: HeaderUser }) {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const menuID = useId();
  const label = userLabel(props.user);

  useEffect(() => {
    if (!open) {
      return;
    }

    function closeOnPointerDown(event: PointerEvent) {
      if (event.target instanceof Node && !containerRef.current?.contains(event.target)) {
        setOpen(false);
      }
    }

    function closeOnEscape(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setOpen(false);
      }
    }

    document.addEventListener("pointerdown", closeOnPointerDown);
    document.addEventListener("keydown", closeOnEscape);
    return () => {
      document.removeEventListener("pointerdown", closeOnPointerDown);
      document.removeEventListener("keydown", closeOnEscape);
    };
  }, [open]);

  return (
    <div className="site-user-menu" ref={containerRef}>
      <button
        aria-controls={menuID}
        aria-expanded={open ? "true" : "false"}
        aria-haspopup="menu"
        className="site-user-menu-trigger"
        onClick={() => setOpen((current) => !current)}
        type="button"
      >
        {label}
        <span aria-hidden="true" className="site-user-menu-caret" />
      </button>
      <div className="site-user-menu-panel" data-open={open ? "true" : "false"} id={menuID} role="menu">
        <a className="site-user-menu-item" href={userHref(props.user)} role="menuitem">
          {label}
        </a>
        <a className="site-user-menu-item" href="/profile" role="menuitem">
          设置
        </a>
        <LogoutForm buttonClassName="site-user-menu-item site-user-menu-item-danger" menuItem />
      </div>
    </div>
  );
}

function LogoutForm(props: { buttonClassName: string; menuItem?: boolean }) {
  return (
    <form action="/auth/logout" method="POST" role={props.menuItem ? "none" : undefined}>
      <input name="_csrf" type="hidden" value={csrfToken()} />
      <button className={props.buttonClassName} role={props.menuItem ? "menuitem" : undefined} type="submit">
        退出
      </button>
    </form>
  );
}

function headerUser(user: User | undefined): HeaderUser | null {
  const username = (user?.username || "").trim();
  if (!username) {
    return null;
  }
  return { username };
}

function userLabel(user: HeaderUser) {
  return user.username || "用户";
}

function userHref(user: HeaderUser) {
  if (!user.username) {
    return "/profile";
  }
  return `/users/${encodeURIComponent(user.username)}`;
}

export { HeaderActions };
