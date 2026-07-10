import "./styles/globals.css";
import "./styles/app.css";
import "./styles/template-ui.css";

import { HeaderActions } from "./components/HeaderActions";
import { mountReactIslands, onDomReady } from "./lib/reactIslands";

function bindNavToggle(root: ParentNode = document) {
  root.querySelectorAll<HTMLButtonElement>("[data-nav-toggle]").forEach((button) => {
    const targetId = button.getAttribute("aria-controls");
    const menu = targetId ? document.getElementById(targetId) : null;
    if (!menu) {
      return;
    }

    button.addEventListener("click", () => {
      const open = menu.getAttribute("data-open") !== "true";
      menu.setAttribute("data-open", open ? "true" : "false");
      button.setAttribute("aria-expanded", open ? "true" : "false");
    });
  });
}

onDomReady(() => {
  bindNavToggle();
  mountReactIslands("[data-header-actions]", HeaderActions);
});
