import { mountReactIslands, onDomReady } from "@/lib/reactIslands";

import { PostActions } from "..";

onDomReady(() => {
  mountReactIslands("[data-post-actions]", PostActions);
});
