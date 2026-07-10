import { StrictMode } from "react";
import type { ComponentType } from "react";
import { createRoot } from "react-dom/client";

type IslandProps = DOMStringMap;
type ReactIslandElement = HTMLElement & {
  __novaReactMounted?: boolean;
};

function mountReactIslands(
  selector: string,
  Component: ComponentType<IslandProps>,
  root: ParentNode = document
) {
  root.querySelectorAll<HTMLElement>(selector).forEach((element) => {
    const mountElement = element as ReactIslandElement;
    if (mountElement.__novaReactMounted) {
      return;
    }

    mountElement.__novaReactMounted = true;
    createRoot(element).render(
      <StrictMode>
        <Component {...element.dataset} />
      </StrictMode>
    );
  });
}

function onDomReady(callback: () => void) {
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", callback);
    return;
  }

  callback();
}

export { mountReactIslands, onDomReady };
